const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const test = require('node:test');
const vm = require('node:vm');

function runAnalyticsScript(overrides = {}) {
  const calls = [];
  let currentNow = overrides.now ?? 1000;
  const listeners = Object.create(null);
  const storage = new Map();

  const sessionStorage = {
    getItem(key) {
      return storage.has(key) ? storage.get(key) : null;
    },
    setItem(key, value) {
      storage.set(key, String(value));
    },
  };

  const context = {
    console,
    fetch: async (url, options) => {
      calls.push({
        url,
        payload: JSON.parse(options.body),
      });
      return { ok: true };
    },
    navigator: {
      sendBeacon: () => false,
    },
    document: {
      referrer: overrides.referrer ?? '',
      visibilityState: 'visible',
      addEventListener(type, handler) {
        listeners[type] = handler;
      },
    },
    location: {
      pathname: overrides.pathname ?? '/',
      search: overrides.search ?? '',
    },
    sessionStorage,
    addEventListener(type, handler) {
      listeners[type] = handler;
    },
    crypto: {
      randomUUID: () => overrides.sessionId ?? 'session-1',
    },
    Date: class extends Date {
      static now() {
        return currentNow;
      }
    },
    Blob: class Blob {
      constructor(parts, options) {
        this.parts = parts;
        this.type = options?.type;
      }
    },
  };

  context.window = context;

  const script = fs.readFileSync(path.join(__dirname, 'analytics.js'), 'utf8');
  vm.runInNewContext(script, context, { filename: 'analytics.js' });

  return {
    calls,
    listeners,
    sessionStorage,
    setNow(value) {
      currentNow = value;
    },
  };
}

test('analytics starts the session before sending the first pageview and preserves duration', async () => {
  const env = runAnalyticsScript({ now: 1000 });

  assert.equal(typeof env.listeners.DOMContentLoaded, 'function');
  await env.listeners.DOMContentLoaded();

  assert.equal(env.sessionStorage.getItem('cricidev.session.started'), '1000');
  assert.equal(env.sessionStorage.getItem('cricidev.session.active'), '1');
  assert.equal(env.calls.length, 2);
  assert.equal(env.calls[0].payload.type, 'session_start');
  assert.equal(env.calls[1].payload.type, 'pageview');

  env.setNow(5000);
  assert.equal(typeof env.listeners.pagehide, 'function');
  env.listeners.pagehide();

  assert.equal(env.calls.length, 3);
  assert.equal(env.calls[2].payload.type, 'session_end');
  assert.equal(env.calls[2].payload.duration_ms, 4000);
});
