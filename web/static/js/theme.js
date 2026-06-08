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

  function setVal(id, val) {
    var el = document.getElementById(id);
    if (el && val !== undefined && val !== null) el.value = val;
  }

  function loadConfigIntoModal(cfg) {
    if (cfg.dashboard) {
      setVal("cfg-recentDays", cfg.dashboard.recentDays);
    }
    if (cfg.tests) {
      var t = cfg.tests;
      if (t.icmp) {
        setVal("cfg-icmp-packetCount", t.icmp.packetCount);
        setVal("cfg-icmp-timeoutSeconds", t.icmp.timeoutSeconds);
      }
      if (t.speedTestURLs) {
        setVal("cfg-downloadUrls", (t.speedTestURLs.downloadUrls || []).join("\n"));
        setVal("cfg-uploadUrls", (t.speedTestURLs.uploadUrls || []).join("\n"));
      }
      if (t.routeTest) {
        setVal("cfg-route-target", t.routeTest.target);
        setVal("cfg-route-maxHops", t.routeTest.maxHops);
        setVal("cfg-route-timeoutSeconds", t.routeTest.timeoutSeconds);
      }
      if (t.jitterTest) {
        setVal("cfg-jitter-target", t.jitterTest.target);
        setVal("cfg-jitter-packetCount", t.jitterTest.packetCount);
        setVal("cfg-jitter-timeoutSeconds", t.jitterTest.timeoutSeconds);
      }
      if (t.bandwidth) {
        setVal("cfg-bw-initialConnections", t.bandwidth.initialConnections);
        setVal("cfg-bw-maxConnections", t.bandwidth.maxConnections);
        setVal("cfg-bw-rampUpStep", t.bandwidth.rampUpStep);
        setVal("cfg-bw-failThreshold", t.bandwidth.failThreshold);
        setVal("cfg-bw-downloadUrl", t.bandwidth.downloadUrl);
      }
    }
  }

  /* Modal open/close (mirrors your existing showModal/closeModal pattern). */
  window.showSettingsModal = function () {
    var m = document.getElementById("settings-modal");
    if (!m) return;
    m.style.display = "block";
    fetch("/config")
      .then(function (r) { return r.json(); })
      .then(loadConfigIntoModal)
      .catch(function (e) { console.error("Failed to load config", e); });
  };

  window.closeSettingsModal = function () {
    var m = document.getElementById("settings-modal");
    if (m) m.style.display = "none";
    var err = document.getElementById("settings-save-error");
    if (err) err.style.display = "none";
  };

  window.saveSettingsConfig = function () {
    function getInt(id) {
      var el = document.getElementById(id);
      return el ? parseInt(el.value, 10) || 0 : 0;
    }
    function getNum(id) {
      var el = document.getElementById(id);
      return el ? parseFloat(el.value) || 0 : 0;
    }
    function getStr(id) {
      var el = document.getElementById(id);
      return el ? el.value.trim() : "";
    }
    function getLines(id) {
      var el = document.getElementById(id);
      if (!el) return [];
      return el.value.split("\n").map(function (s) { return s.trim(); }).filter(function (s) { return s.length > 0; });
    }

    var payload = {
      dashboard: { recentDays: getInt("cfg-recentDays") },
      tests: {
        icmp: {
          packetCount: getInt("cfg-icmp-packetCount"),
          timeoutSeconds: getInt("cfg-icmp-timeoutSeconds")
        },
        speedTestURLs: {
          downloadUrls: getLines("cfg-downloadUrls"),
          uploadUrls: getLines("cfg-uploadUrls")
        },
        routeTest: {
          target: getStr("cfg-route-target"),
          maxHops: getInt("cfg-route-maxHops"),
          timeoutSeconds: getInt("cfg-route-timeoutSeconds")
        },
        jitterTest: {
          target: getStr("cfg-jitter-target"),
          packetCount: getInt("cfg-jitter-packetCount"),
          timeoutSeconds: getInt("cfg-jitter-timeoutSeconds")
        },
        bandwidth: {
          initialConnections: getInt("cfg-bw-initialConnections"),
          maxConnections: getInt("cfg-bw-maxConnections"),
          rampUpStep: getInt("cfg-bw-rampUpStep"),
          failThreshold: getNum("cfg-bw-failThreshold"),
          downloadUrl: getStr("cfg-bw-downloadUrl")
        }
      }
    };

    fetch("/config", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    })
      .then(function (r) {
        if (r.ok) {
          window.closeSettingsModal();
        } else {
          return r.text().then(function (t) {
            var err = document.getElementById("settings-save-error");
            if (err) { err.textContent = t; err.style.display = "block"; }
          });
        }
      })
      .catch(function (e) {
        var err = document.getElementById("settings-save-error");
        if (err) { err.textContent = e.message; err.style.display = "block"; }
      });
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