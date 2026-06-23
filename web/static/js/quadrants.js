/* quadrants.js — injects a collapse + fullscreen toolbar into each
 * dashboard quadrant.
 *   - Collapse: mobile only; shrinks a quadrant to its title bar.
 *   - Fullscreen: all devices; CSS-based so it works on iOS too.
 * Load anywhere; it wires up on DOMContentLoaded.
 */
(function () {
  "use strict";

  var TITLES = {
    "test-quadrant": "Network Tests",
    "control-quadrant": "Control Panel",
    "scheduler-quadrant": "Task Scheduler"
  };

  var CHEVRON =
    '<svg class="chev" width="14" height="14" viewBox="0 0 16 16" fill="none" ' +
    'stroke="currentColor" stroke-width="2" stroke-linecap="round" ' +
    'stroke-linejoin="round"><path d="M4 6l4 4 4-4"/></svg>';

  var EXPAND =
    '<svg class="ic-expand" width="14" height="14" viewBox="0 0 16 16" fill="none" ' +
    'stroke="currentColor" stroke-width="2" stroke-linecap="round" ' +
    'stroke-linejoin="round"><path d="M6 2H2v4M10 2h4v4M6 14H2v-4M10 14h4v-4"/></svg>';

  var COMPRESS =
    '<svg class="ic-compress" width="14" height="14" viewBox="0 0 16 16" fill="none" ' +
    'stroke="currentColor" stroke-width="2" stroke-linecap="round" ' +
    'stroke-linejoin="round"><path d="M2 6h4V2M14 6h-4V2M2 10h4v4M14 10h-4v4"/></svg>';

  function clearFullscreen() {
    var open = document.querySelectorAll(".quadrant.is-fullscreen");
    for (var i = 0; i < open.length; i++) open[i].classList.remove("is-fullscreen");
    document.body.classList.remove("has-fullscreen-quadrant");
  }

  function toggleFullscreen(quadrant) {
    var isOn = quadrant.classList.contains("is-fullscreen");
    clearFullscreen();                       // only one at a time
    if (!isOn) {
      quadrant.classList.remove("is-collapsed");
      quadrant.classList.add("is-fullscreen");
      document.body.classList.add("has-fullscreen-quadrant");
    }
  }

  function build(quadrant) {
    if (quadrant.querySelector(".quadrant-bar")) return;   // already built

    var title = TITLES[quadrant.id] || quadrant.dataset.quadrantTitle || "Panel";

    var bar = document.createElement("div");
    bar.className = "quadrant-bar";
    bar.innerHTML =
      '<span class="quadrant-title">' + title + "</span>" +
      '<div class="quadrant-tools">' +
        '<button type="button" class="quad-btn quad-collapse" ' +
          'title="Collapse / expand" aria-label="Collapse or expand panel">' +
          CHEVRON +
        "</button>" +
        '<button type="button" class="quad-btn quad-fs" ' +
          'title="Toggle fullscreen" aria-label="Toggle fullscreen">' +
          EXPAND + COMPRESS +
        "</button>" +
      "</div>";

    quadrant.insertBefore(bar, quadrant.firstChild);

    bar.querySelector(".quad-collapse").addEventListener("click", function () {
      quadrant.classList.toggle("is-collapsed");
    });
    bar.querySelector(".quad-fs").addEventListener("click", function () {
      toggleFullscreen(quadrant);
    });
  }

  function buildAll() {
    var quads = document.querySelectorAll(".dashboard .quadrant");
    for (var i = 0; i < quads.length; i++) build(quads[i]);
  }

  document.addEventListener("DOMContentLoaded", buildAll);

  // htmx swaps (e.g. outerHTML-replacing #test-quadrant after a delete, or
  // a quick test/chart run) drop in fresh server HTML without the
  // JS-injected toolbar. Listening for specific htmx events is brittle
  // (depends on exactly which element/swap style fired them), so instead
  // watch the DOM directly: any time a .quadrant loses its .quadrant-bar,
  // put it back. build() is a no-op if the bar is already there.
  var dashboard = document.querySelector(".dashboard");
  if (dashboard && window.MutationObserver) {
    new MutationObserver(buildAll).observe(dashboard, { childList: true, subtree: true });
  }

  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape") clearFullscreen();
  });
})();