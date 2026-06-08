/* theme.js — Light / Dark / System theme switching with persistence.
 *
 * Load this in <head> (after the stylesheet link) so the saved theme is
 * applied to <html> before the body paints. For zero flash you can also
 * inline the two lines marked "EARLY APPLY" directly in <head>.
 */
(function () {
  "use strict";

  var STORAGE_KEY = "dashboard-theme";          // stores "light" | "dark" | "system"
  var root = document.documentElement;
  var media = window.matchMedia("(prefers-color-scheme: dark)");

  function getPref() {
    return localStorage.getItem(STORAGE_KEY) || "system";
  }

  function resolve(pref) {
    if (pref === "system") return media.matches ? "dark" : "light";
    return pref;
  }

  function apply(pref) {
    root.setAttribute("data-theme", resolve(pref));
  }

  function setActiveButtons(pref) {
    var btns = document.querySelectorAll(".theme-option");
    for (var i = 0; i < btns.length; i++) {
      btns[i].classList.toggle("active", btns[i].dataset.pref === pref);
    }
  }

  function setPref(pref) {
    localStorage.setItem(STORAGE_KEY, pref);
    apply(pref);
    setActiveButtons(pref);
  }

  /* EARLY APPLY — runs immediately, before <body> renders. */
  apply(getPref());

  /* Follow the OS when in "system" mode. */
  media.addEventListener("change", function () {
    if (getPref() === "system") apply("system");
  });

  /* Modal open/close (mirrors your existing showModal/closeModal pattern). */
  window.showSettingsModal = function () {
    var m = document.getElementById("settings-modal");
    if (m) m.style.display = "block";
  };
  window.closeSettingsModal = function () {
    var m = document.getElementById("settings-modal");
    if (m) m.style.display = "none";
  };

  document.addEventListener("DOMContentLoaded", function () {
    setActiveButtons(getPref());

    var btns = document.querySelectorAll(".theme-option");
    for (var i = 0; i < btns.length; i++) {
      btns[i].addEventListener("click", function () {
        setPref(this.dataset.pref);
      });
    }

    /* close when clicking the backdrop */
    var modal = document.getElementById("settings-modal");
    if (modal) {
      modal.addEventListener("click", function (e) {
        if (e.target === modal) window.closeSettingsModal();
      });
    }
  });
})();