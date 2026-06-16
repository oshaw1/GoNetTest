/* themedSelect.js — progressively enhances a <select data-themed-select>
 * into a fully themeable custom dropdown (the native option list can't be
 * restyled across browsers). The original <select> stays in the DOM and
 * keeps its value in sync, so forms/htmx (hx-include, form serialization,
 * change events, etc.) keep working unmodified.
 *
 * The open option list is rendered as a direct child of <body>, positioned
 * with `position: fixed` from the trigger's bounding rect. This is a
 * "portal" so the popup is never clipped or covered by an ancestor's
 * `overflow: hidden`/`auto` or stacking context (the dashboard's quadrants
 * use both heavily) — it always renders above everything else, at a fixed
 * top-level z-index.
 *
 * Markup in:
 *   <select id="x" data-themed-select>
 *     <option value="a">A</option>
 *   </select>
 *
 * Markup out:
 *   <div class="themed-select">
 *     <select id="x" data-themed-select class="themed-select-native">...</select>
 *     <button class="themed-select-trigger">A <svg.../></button>
 *   </div>
 *   ...
 *   <body>
 *     ...
 *     <ul class="themed-select-options">  <!-- appended here, position: fixed -->
 *       <li data-value="a" class="is-selected">A</li>
 *     </ul>
 *   </body>
 *
 * If code sets select.value programmatically (or calls form.reset()),
 * no native "change" event fires, so the custom UI won't know to update
 * on its own — call window.ThemedSelect.refresh(select) (or refreshAll)
 * afterwards to resync the label/highlighted option.
 */
(function () {
  "use strict";

  var CARET =
    '<svg class="themed-select-caret" width="10" height="10" viewBox="0 0 16 16" fill="none" ' +
    'stroke="currentColor" stroke-width="2" stroke-linecap="round" ' +
    'stroke-linejoin="round"><path d="M4 6l4 4 4-4"/></svg>';

  var openInstance = null; // the currently-open build() instance, if any

  function closeOpen() {
    if (openInstance) openInstance.close();
  }

  // Resyncs a built themed-select's trigger label + highlighted option from
  // the underlying <select>'s current value. Safe to call on any select;
  // no-ops if it hasn't been (or can't be) built yet.
  function refresh(select) {
    if (typeof select === "string") select = document.getElementById(select);
    if (!select || !select.themedSelect) return;
    select.themedSelect.sync();
  }

  function refreshAll(root) {
    var selects = (root || document).querySelectorAll("[data-themed-select]");
    for (var i = 0; i < selects.length; i++) refresh(selects[i]);
  }

  function build(select) {
    if (select.themedSelect) return; // already built

    var wrap = document.createElement("div");
    wrap.className = "themed-select";

    var trigger = document.createElement("button");
    trigger.type = "button";
    trigger.className = "themed-select-trigger";

    var label = document.createElement("span");
    label.className = "themed-select-label";

    var list = document.createElement("ul");
    list.className = "themed-select-options";
    list.setAttribute("role", "listbox");

    function sync() {
      var opt = select.options[select.selectedIndex];
      label.textContent = opt ? opt.textContent : "";
      var items = list.querySelectorAll("li");
      for (var i = 0; i < items.length; i++) {
        items[i].classList.toggle("is-selected", items[i].dataset.value === select.value);
      }
    }

    function reposition() {
      var rect = trigger.getBoundingClientRect();
      list.style.left = rect.left + "px";
      list.style.minWidth = rect.width + "px";

      // Default below the trigger; flip above it if there isn't room.
      var listHeight = list.offsetHeight;
      var opensBelow = rect.bottom + 4 + listHeight <= window.innerHeight;
      if (opensBelow) {
        list.style.top = rect.bottom + 4 + "px";
        list.style.bottom = "";
      } else {
        list.style.top = "";
        list.style.bottom = window.innerHeight - rect.top + 4 + "px";
      }
    }

    function open() {
      if (openInstance && openInstance !== instance) openInstance.close();
      document.body.appendChild(list); // portal: escape any clipping ancestor
      list.style.position = "fixed";
      wrap.classList.add("is-open");
      reposition();
      list.classList.add("is-visible");
      window.addEventListener("scroll", reposition, true);
      window.addEventListener("resize", reposition, true);
      openInstance = instance;
    }

    function close() {
      wrap.classList.remove("is-open");
      list.classList.remove("is-visible");
      window.removeEventListener("scroll", reposition, true);
      window.removeEventListener("resize", reposition, true);
      if (openInstance === instance) openInstance = null;
    }

    function choose(opt) {
      var changed = select.value !== opt.value;
      select.value = opt.value;
      sync();
      close();
      if (changed) select.dispatchEvent(new Event("change", { bubbles: true }));
    }

    Array.prototype.forEach.call(select.options, function (opt) {
      var li = document.createElement("li");
      li.textContent = opt.textContent;
      li.dataset.value = opt.value;
      li.setAttribute("role", "option");
      li.tabIndex = -1;
      li.addEventListener("click", function () {
        choose(opt);
      });
      list.appendChild(li);
    });

    trigger.appendChild(label);
    trigger.insertAdjacentHTML("beforeend", CARET);

    trigger.addEventListener("click", function (e) {
      e.stopPropagation();
      if (wrap.classList.contains("is-open")) close();
      else open();
    });

    list.addEventListener("keydown", function (e) {
      if (e.key === "Escape") {
        close();
        trigger.focus();
      }
    });
    list.addEventListener("click", function (e) {
      e.stopPropagation();
    });

    select.parentNode.insertBefore(wrap, select);
    wrap.appendChild(select);
    wrap.appendChild(trigger);
    // `list` is intentionally NOT appended into wrap — it's portaled to
    // <body> only while open (see open()) so nothing can clip it.

    select.classList.add("themed-select-native");
    select.tabIndex = -1;
    select.addEventListener("change", sync);

    var instance = { sync: sync, open: open, close: close };
    select.themedSelect = instance;

    sync();
  }

  function scanAndBuild(root) {
    var selects = (root || document).querySelectorAll("[data-themed-select]");
    for (var i = 0; i < selects.length; i++) build(selects[i]);
  }

  document.addEventListener("DOMContentLoaded", function () {
    scanAndBuild(document);
  });

  // Rebuild any selects that arrive via htmx swaps (e.g. the test-type
  // dropdown, which the server re-renders as a plain <select> each time).
  document.addEventListener("htmx:load", function (e) {
    scanAndBuild((e.detail && e.detail.elt) || e.target);
  });

  document.addEventListener("click", closeOpen);
  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape") closeOpen();
  });

  window.ThemedSelect = { refresh: refresh, refreshAll: refreshAll };
})();
