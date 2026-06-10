(() => {
  const endpoint = "/api/events";
  const sessionKey = "cricidev.session";
  const startKey = "cricidev.session.started";
  const activeKey = "cricidev.session.active";

  const now = () => Date.now();
  const pagePath = () => window.location.pathname + window.location.search;
  const referrer = () => document.referrer || "";
  const sessionId = () => {
    const current = window.sessionStorage.getItem(sessionKey);
    if (current) return current;
    const fresh = crypto.randomUUID();
    window.sessionStorage.setItem(sessionKey, fresh);
    return fresh;
  };
  const startedAt = () => {
    const current = Number(window.sessionStorage.getItem(startKey));
    if (current > 0) return current;
    const fresh = now();
    window.sessionStorage.setItem(startKey, String(fresh));
    return fresh;
  };

  const send = async (payload, useBeacon = false) => {
    const body = JSON.stringify(payload);
    if (useBeacon && navigator.sendBeacon) {
      const ok = navigator.sendBeacon(endpoint, new Blob([body], { type: "application/json" }));
      if (ok) return true;
    }

    try {
      await fetch(endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body,
        keepalive: true,
        credentials: "same-origin",
      });
    } catch (_) {}

    return false;
  };

  const basePayload = () => ({
    session_id: sessionId(),
    path: pagePath(),
    referrer: referrer(),
    created_at: new Date().toISOString(),
  });

  const track = (type, extra = {}, useBeacon = false) => {
    return send({ type, ...basePayload(), ...extra }, useBeacon);
  };

  const trackClick = (event) => {
    const target = event.target.closest("[data-track]");
    if (!target) return;

    track("click", {
      target: target.getAttribute("data-track"),
      label: (target.textContent || "").trim().slice(0, 160),
    });
  };

  const finishSession = () => {
    if (window.sessionStorage.getItem(activeKey) !== "1") return;
    window.sessionStorage.setItem(activeKey, "0");
    track(
      "session_end",
      {
        duration_ms: Math.max(0, now() - startedAt()),
      },
      true,
    );
  };

  const startSession = async () => {
    if (window.sessionStorage.getItem(activeKey) === "1") return;
    const started = now();
    window.sessionStorage.setItem(activeKey, "1");
    window.sessionStorage.setItem(startKey, String(started));

    await track("session_start", {
      duration_ms: 0,
    });
    track("pageview");
  };

  document.addEventListener("click", trackClick, true);
  document.addEventListener("DOMContentLoaded", startSession, { once: true });
  window.addEventListener("pagehide", finishSession);
  document.addEventListener("visibilitychange", () => {
    if (document.visibilityState === "hidden") {
      finishSession();
    }
  });
})();
