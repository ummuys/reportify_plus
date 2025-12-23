import { BASE_PATH, GET_ACCESS_TOKEN_PATH } from '../config/index.js';
import { showAlert, applyRoleRestrictions } from '../ui/index.js';
import { syncRoleFromToken } from '../core/auth.js';

const AUTH_ALERT_TITLE = "Авторизация";
const AUTH_ALERT_MESSAGE = "Сессия истекла. Войдите снова.";

let refreshPromise = null;
let authModalShown = false;

function isRefreshRequest(url) {
  if (!url) return false;
  if (url === GET_ACCESS_TOKEN_PATH) return true;
  try {
    const parsed = new URL(url, window.location.origin);
    return parsed.pathname.endsWith('/api/v1/secure/refresh');
  } catch {
    return String(url).includes('/api/v1/secure/refresh');
  }
}

async function requestNewAccessToken() {
  const res = await fetch(GET_ACCESS_TOKEN_PATH, {
    method: "GET",
    credentials: "include",
    headers: { "Accept": "application/json" }
  });

  if (res.status === 401) {
    return { ok: false, reason: "unauthorized" };
  }
  if (!res.ok) {
    const txt = await res.text().catch(() => "");
    return { ok: false, reason: `${res.status} ${res.statusText}${txt ? " - " + txt : ""}` };
  }

  let token = "";
  try {
    const data = await res.json();
    token = data.access_token || data.token || data.access || "";
  } catch {
    const txt = await res.text().catch(() => "");
    token = txt.trim();
  }

  if (token) {
    localStorage.setItem("access_token_v1", token);
    const role = syncRoleFromToken(token);
    applyRoleRestrictions(role);
    return { ok: true, token };
  }
  return { ok: false, reason: "empty_token" };
}

async function refreshAccessToken() {
  if (!refreshPromise) {
    refreshPromise = (async () => {
      try {
        return await requestNewAccessToken();
      } catch (err) {
        return { ok: false, reason: err?.message || "refresh_error" };
      }
    })();

    refreshPromise.finally(() => {
      refreshPromise = null;
    });
  }
  return refreshPromise;
}

async function handleAuthFailure() {
  if (authModalShown) return;
  authModalShown = true;

  try {
    await showAlert(AUTH_ALERT_MESSAGE, AUTH_ALERT_TITLE);
  } catch {}

  localStorage.removeItem("access_token_v1");
  document.cookie = "refresh_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
  location.assign(BASE_PATH);
}

export async function fetchWithToken(url, options = {}) {
  const doFetch = async () => {
    const token = localStorage.getItem("access_token_v1") || "";
    const opts = { ...options };
    opts.headers = { ...(opts.headers || {}) };

    const requestedAccept = String(
      opts.headers["Accept"] || opts.headers["accept"] || "application/json"
    ).toLowerCase();
    if (!opts.headers["Accept"] && !opts.headers["accept"]) {
      opts.headers["Accept"] = "application/json";
    }

    if (token) opts.headers["Authorization"] = "Bearer " + token;
    if (!("credentials" in opts)) opts.credentials = "include";

    let res;
    try {
      res = await fetch(url, opts);
    } catch (e) {
      const msg = (e && (e.message || e.toString())) ? (e.message || e.toString()) : String(e);
      throw new Error(`Network error while fetching ${url}: ${msg}`);
    }
    return { res, requestedAccept };
  };

  let { res, requestedAccept } = await doFetch();

  if (res.status === 401) {
    if (isRefreshRequest(url)) {
      await handleAuthFailure();
      throw new Error("401 Unauthorized - refresh token expired");
    }

    const refresh = await refreshAccessToken();
    if (refresh && refresh.ok) {
      ({ res, requestedAccept } = await doFetch());
    } else {
      await handleAuthFailure();
      throw new Error("401 Unauthorized - refresh failed");
    }
  }

  if (res.status === 401) {
    await handleAuthFailure();
    throw new Error("401 Unauthorized - redirect to login");
  }

  if (!res.ok) {
    const txt = await res.text().catch(() => "");
    throw new Error(`${res.status} ${res.statusText}${txt ? " - " + txt : ""}`);
  }

  const ct = (res.headers.get("content-type") || "").toLowerCase();

  const isBinaryRequested =
    requestedAccept.includes("application/pdf") ||
    requestedAccept.includes("text/csv") ||
    requestedAccept.includes("application/vnd.openxmlformats-officedocument.wordprocessingml.document") ||
    requestedAccept.includes("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") ||
    requestedAccept.includes("application/zip") ||
    requestedAccept.includes("application/octet-stream");

  if (isBinaryRequested) return res;

  const isBinaryResponse =
    ct.includes("application/pdf") ||
    ct.includes("text/csv") ||
    ct.includes("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") ||
    ct.includes("application/vnd.openxmlformats-officedocument.wordprocessingml.document") ||
    ct.includes("application/zip") ||
    ct.includes("application/octet-stream") ||
    (!ct && requestedAccept !== "application/json");

  if (isBinaryResponse) return res;

  const txt = await res.text();
  try { return JSON.parse(txt); } catch { return txt; }
}

export async function getJSON(url) {
  return fetchWithToken(url);
}



