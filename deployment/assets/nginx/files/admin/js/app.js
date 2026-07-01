const API_HOST = "http://127.0.0.1:8088";
const API_PREFIX = `${API_HOST}/api/v1`;
const ENDPOINTS = {
  users: "/admin/users",
  deleteUser: (userId) => `/admin/users/${encodeURIComponent(userId)}`
};

const state = {
  token: "",
  role: "",
};

const AUTH_ALERT_TITLE = "Авторизация";
const AUTH_ALERT_MESSAGE = "Сессия истекла. Войдите снова.";
const REFRESH_URL = `${API_PREFIX}/secure/refresh`;

let refreshPromise = null;
let authModalShown = false;

const roleBadge = document.getElementById("roleBadge");
const statusLog = document.getElementById("statusLog");
const usersTableBody = document.getElementById("usersTableBody");
const btnReloadUsers = document.getElementById("btnReloadUsers");
const formCreateUser = document.getElementById("formCreateUser");
const formUpdateUser = document.getElementById("formUpdateUser");
const formDeleteUser = document.getElementById("formDeleteUser");
const logoutBtn = document.getElementById("logoutBtn");
const bodyEl = document.body;

function showAlert(message, title = "Сообщение") {
  return new Promise((resolve) => {
    const modal = document.getElementById("alertModal");
    const msgEl = document.getElementById("alertMessage");
    const titleEl = document.getElementById("alertTitle");
    const btnOk = document.getElementById("alertOk");

    if (!modal || !msgEl || !titleEl || !btnOk) {
      alert(message);
      resolve();
      return;
    }

    msgEl.textContent = message;
    titleEl.textContent = title;
    modal.style.display = "flex";

    const close = () => {
      modal.style.display = "none";
      btnOk.removeEventListener("click", close);
      document.removeEventListener("keydown", keyHandler);
      resolve();
    };

    const keyHandler = (event) => {
      if (event.key === "Enter" || event.key === "Escape") {
        event.preventDefault();
        close();
      }
    };

    btnOk.addEventListener("click", close);
    document.addEventListener("keydown", keyHandler);
  });
}

function isRefreshRequest(pathOrUrl) {
  if (!pathOrUrl) return false;
  if (pathOrUrl === REFRESH_URL) return true;
  return String(pathOrUrl).includes("/api/v1/secure/refresh");
}

async function requestNewAccessToken() {
  const res = await fetch(REFRESH_URL, {
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
    state.token = token;

    const payload = decodeTokenPayload(token);
    if (payload?.role) {
      state.role = payload.role;
      if (roleBadge) {
        roleBadge.textContent = `role: ${payload.role}`;
      }
    }
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

async function handleAuthFailure(message = AUTH_ALERT_MESSAGE) {
  if (authModalShown) return false;
  authModalShown = true;

  try {
    await showAlert(message, AUTH_ALERT_TITLE);
  } catch {}

  localStorage.removeItem("access_token_v1");
  document.cookie = "refresh_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
  window.location.assign("/");
  return false;
}
async function redirectToLogin(message = "") {
  if (message) {
    await showAlert(message, "Авторизация");
  }
  localStorage.removeItem("access_token_v1");
  document.cookie = "refresh_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
  window.location.assign("/");
  return false;
}

function decodeTokenPayload(token) {
  try {
    const payload = token.split(".")[1];
    return JSON.parse(atob(payload));
  } catch {
    return null;
  }
}

async function ensureAdminAccess() {
  const token = localStorage.getItem("access_token_v1");
  if (!token) {
    await redirectToLogin("Требуется авторизация.");
    return false;
  }

  const payload = decodeTokenPayload(token);
  if (!payload || payload.role !== "admin") {
    await redirectToLogin("Недостаточно прав");
    return false;
  }

  state.token = token;
  state.role = payload.role;

  if (roleBadge) {
    roleBadge.textContent = `role: ${payload.role}`;
  }
  return true;
}

function logStatus(message, kind = "info") {
  if (!statusLog) return;
  const entry = document.createElement("div");
  entry.className = `status-entry ${kind === "error" ? "error" : "success"}`;
  entry.innerHTML = `<span>${message}</span><time>${new Date().toLocaleTimeString()}</time>`;
  statusLog.prepend(entry);
  const maxEntries = 12;
  while (statusLog.children.length > maxEntries) {
    statusLog.removeChild(statusLog.lastChild);
  }
}

async function fetchJSON(path, options = {}) {
  let token = state.token || localStorage.getItem("access_token_v1");
  if (!token) {
    await redirectToLogin("Требуется авторизация.");
    return null;
  }

  const buildOpts = () => {
    const opts = { ...options };
    opts.headers = { "Content-Type": "application/json", ...(options.headers || {}) };
    opts.credentials = "include";
    if (token) {
      opts.headers["Authorization"] = `Bearer ${token}`;
    }
    return opts;
  };

  const doFetch = () => fetch(`${API_PREFIX}${path}`, buildOpts());

  let res = await doFetch();
  if (res.status === 401) {
    if (isRefreshRequest(path)) {
      await handleAuthFailure();
      return null;
    }

    const refresh = await refreshAccessToken();
    if (refresh && refresh.ok) {
      token = state.token || localStorage.getItem("access_token_v1") || refresh.token;
      res = await doFetch();
    } else {
      await handleAuthFailure();
      return null;
    }
  }

  if (res.status === 401) {
    await handleAuthFailure();
    return null;
  }

  if (res.status === 204) {
    return {};
  }
  const text = await res.text();
  if (!res.ok) {
    let errorMessage = text || `${res.status} ${res.statusText}`;
    try {
      const parsed = JSON.parse(text);
      errorMessage = parsed?.err || parsed?.msg || errorMessage;
    } catch {
      /* noop */
    }
    throw new Error(errorMessage);
  }

  if (!text) return {};
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

function renderUsers(users = []) {
  if (!usersTableBody) return;
  if (!users.length) {
    usersTableBody.innerHTML = `<tr><td colspan="3">Нет данных</td></tr>`;
    return;
  }

  usersTableBody.innerHTML = users
    .map(
      (user) => `
      <tr>
        <td>${user.user_id}</td>
        <td>${user.username}</td>
        <td>${user.role}</td>
      </tr>`
    )
    .join("");
}

async function loadUsers() {
  try {
    const data = await fetchJSON(ENDPOINTS.users);
    const rawUsers = Array.isArray(data) ? data : (data?.users || []);

    // нормализуем поля под то, что ждёт renderUsers
    const users = rawUsers.map((u) => ({
      user_id: u.user_id ?? u.UserID,
      username: u.username ?? u.Username,
      role: u.role ?? u.Role,
    }));

    renderUsers(users);
    logStatus("Список пользователей обновлён", "success");
  } catch (err) {
    console.error(err);
    logStatus(`Ошибка загрузки пользователей: ${err.message}`, "error");
  }
}

function openModal(id) {
  const modal = document.getElementById(id);
  if (!modal) return;
  modal.classList.add("show");
  bodyEl?.classList.add("modal-open");
  if (id === "modalUsers") {
    loadUsers();
  }
}

function closeModal(modal) {
  if (!modal) return;
  modal.classList.remove("show");
  if (!document.querySelector(".admin-modal.show")) {
    bodyEl?.classList.remove("modal-open");
  }
}

function initModalTriggers() {
  document.querySelectorAll("[data-modal]").forEach((button) => {
    button.addEventListener("click", () => openModal(button.dataset.modal));
  });

  document.querySelectorAll("[data-close-modal]").forEach((btn) => {
    btn.addEventListener("click", () => closeModal(btn.closest(".admin-modal")));
  });

  document.querySelectorAll(".admin-modal").forEach((modal) => {
    modal.addEventListener("click", (event) => {
      if (event.target === modal) {
        closeModal(modal);
      }
    });
  });
}

function initForms() {
  formCreateUser?.addEventListener("submit", async (event) => {
    event.preventDefault();
    const formEl = event.currentTarget;
    const formData = new FormData(event.currentTarget);
    const payload = {
      username: (formData.get("username") || "").trim(),
      password: formData.get("password") || "",
      role: formData.get("role") || "user",
    };
    try {
      await fetchJSON(ENDPOINTS.users, {
        method: "POST",
        body: JSON.stringify(payload),
      });
      logStatus(`Пользователь ${payload.username} создан`, "success");
      formEl?.reset();
      await loadUsers();
    } catch (err) {
      logStatus(`Ошибка создания: ${err.message}`, "error");
    }
  });

  formUpdateUser?.addEventListener("submit", async (event) => {
    event.preventDefault();
    const formEl = event.currentTarget;
    const formData = new FormData(event.currentTarget);
    const body = {
      user_id: String(formData.get("userId") || "").trim(),
      username: (formData.get("username") || "").trim(),
      password: formData.get("password") || "",
      role: formData.get("role") || "",
    };

    if (!body.username) delete body.username;
    if (!body.password) delete body.password;
    if (!body.role) delete body.role;

    try {
      await fetchJSON(ENDPOINTS.users, {
        method: "PUT",
        body: JSON.stringify(body),
      });
      logStatus(`Пользователь #${body.user_id} обновлён`, "success");
      formEl?.reset();
      await loadUsers();
    } catch (err) {
      logStatus(`Ошибка обновления: ${err.message}`, "error");
    }
  });

  formDeleteUser?.addEventListener("submit", async (event) => {
    event.preventDefault();
    const formEl = event.currentTarget;
    const formData = new FormData(event.currentTarget);
    const userId = String(formData.get("userId") || "").trim();
    if (!userId) return;
    try {
      await fetchJSON(ENDPOINTS.deleteUser(userId), { method: "DELETE" });
      logStatus(`Пользователь ${userId} удалён`, "success");
      formEl?.reset();
      await loadUsers();
    } catch (err) {
      logStatus(`Ошибка удаления: ${err.message}`, "error");
    }
  });
}

function initLogout() {
  logoutBtn?.addEventListener("click", async () => {
    await redirectToLogin();
  });
}

async function bootstrap() {
  const allowed = await ensureAdminAccess();
  if (!allowed) return;
  initModalTriggers();
  initForms();
  initLogout();
  btnReloadUsers?.addEventListener("click", (e) => {
    e.preventDefault();
    loadUsers();
  });
}

document.addEventListener("DOMContentLoaded", bootstrap);
