import { getUserReports } from '../api/index.js';

const REPORTS_POLL_MS = 2000;

function escapeHtml(value) {
  return String(value ?? '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}

function formatDateTime(value) {
  if (!value) return '';
  const date = value instanceof Date ? value : new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  return date.toLocaleString();
}

function normalizeFilePathUrl(path) {
  const raw = String(path || '').trim();
  if (!raw) return raw;
  try {
    const u = new URL(raw);
    if (u.hostname === 'report-minio' || u.hostname.includes('minio')) {
      u.hostname = window.location.hostname;
    }
    if (window.location.protocol && u.protocol !== window.location.protocol) {
      u.protocol = window.location.protocol;
    }
    return u.toString();
  } catch {
    return raw;
  }
}

function statusClass(status) {
  const normalized = String(status || '').trim().toLowerCase();
  if (!normalized) return '';

  const ok = ['done', 'success', 'completed', 'ready', 'finished', 'ok'];
  const run = ['pending', 'queued', 'running', 'in_progress', 'processing', 'created'];
  const err = ['failed', 'error', 'err', 'canceled', 'cancelled'];

  if (ok.includes(normalized)) return 'report-status--ok';
  if (run.includes(normalized)) return 'report-status--run';
  if (err.includes(normalized)) return 'report-status--err';
  return '';
}

function normalizeReportsResponse(payload) {
  if (!payload) return [];
  if (Array.isArray(payload)) return payload;
  if (payload && typeof payload === 'object') {
    if (Array.isArray(payload.reports)) return payload.reports;
    if (Array.isArray(payload.data)) return payload.data;
    if (payload.result && Array.isArray(payload.result.reports)) return payload.result.reports;
  }
  return [];
}

function isSidebarVisible(sidebarEl) {
  if (!sidebarEl || !document.documentElement.contains(sidebarEl)) return false;

  const style = window.getComputedStyle(sidebarEl);
  if (style.display === 'none' || style.visibility === 'hidden') return false;

  const rect = sidebarEl.getBoundingClientRect();
  if (rect.width <= 0 || rect.height <= 0) return false;

  if (rect.right <= 0) return false;
  if (rect.left >= window.innerWidth) return false;

  return true;
}

export function initRightPanelTabs() {
  const sidebar = document.querySelector('.sidebar');
  const tabsRoot = document.getElementById('rightPanelTabs');
  const titleEl = document.getElementById('rightPanelTitle');

  const cacheActions = document.getElementById('cacheActions');
  const reportsActions = document.getElementById('reportsActions');
  const btnReloadReports = document.getElementById('btnReloadReports');

  const historyList = document.getElementById('historyList');
  const reportsPane = document.getElementById('reportsPane');
  const reportsList = document.getElementById('reportsList');

  const detailsOverlay = document.getElementById('reportDetailsModal');
  const detailsTitle = document.getElementById('reportDetailsTitle');
  const detailsBody = document.getElementById('reportDetailsBody');
  const detailsClose = document.getElementById('reportDetailsClose');

  if (!sidebar || !tabsRoot || !titleEl || !cacheActions || !reportsActions || !historyList || !reportsPane || !reportsList) {
    return;
  }

  let activeTab = 'cache';
  let pollTimer = null;
  let abortController = null;
  let reports = [];
  let loading = false;
  let lastError = null;
  let openReportId = null;

  function setActiveTab(next) {
    const tab = next === 'reports' ? 'reports' : 'cache';
    activeTab = tab;

    const buttons = [...tabsRoot.querySelectorAll('.sidebar-tab[data-tab]')];
    buttons.forEach(btn => {
      const isActive = btn.dataset.tab === tab;
      btn.classList.toggle('active', isActive);
      btn.setAttribute('aria-selected', String(isActive));
    });

    titleEl.textContent = tab === 'reports' ? 'Отчёты' : 'Кэш запросов';
    cacheActions.style.display = tab === 'reports' ? 'none' : '';
    reportsActions.style.display = tab === 'reports' ? '' : 'none';

    historyList.style.display = tab === 'reports' ? 'none' : '';
    reportsPane.style.display = tab === 'reports' ? '' : 'none';

    ensurePolling();
  }

  function renderReports() {
    if (loading && !reports.length) {
      reportsList.innerHTML = '<li class="empty-state"><small>Загрузка…</small></li>';
      return;
    }

    if (lastError && !reports.length) {
      reportsList.innerHTML = `<li class="empty-state"><small>${escapeHtml(lastError)}</small></li>`;
      return;
    }

    if (!reports.length) {
      reportsList.innerHTML = '<li class="empty-state"><small>Отчётов пока нет</small></li>';
      return;
    }

    reportsList.innerHTML = reports
      .map(r => {
        const id = r.report_id || r.reportId || r.uuid || r.id || '';
        const name = r.name || (id ? `Отчёт ${id}` : 'Отчёт');
        const status = r.status || '';
        const createdAt = r.created_at || r.createdAt || '';
        const updatedAt = r.updated_at || r.updatedAt || '';
        const timeLabel = updatedAt ? `Обновлён: ${formatDateTime(updatedAt)}` : (createdAt ? `Создан: ${formatDateTime(createdAt)}` : '');

        return `
          <li data-report-id="${escapeHtml(id)}">
            <div class="report-card">
              <div class="report-card__top">
                <div class="report-card__title">${escapeHtml(name)}</div>
                <div class="report-status ${statusClass(status)}">${escapeHtml(status || 'unknown')}</div>
              </div>
              <div class="report-card__meta">
                ${timeLabel ? `<div class="report-card__time">${escapeHtml(timeLabel)}</div>` : ''}
              </div>
            </div>
          </li>
        `;
      })
      .join('');
  }

  function closeDetailsModal() {
    if (!detailsOverlay) return;
    detailsOverlay.style.display = 'none';
    openReportId = null;
  }

  function renderDetailsModal(report) {
    if (!detailsOverlay || !detailsTitle || !detailsBody) return;

    const id = report?.report_id || report?.reportId || report?.uuid || report?.id || '';
    const name = report?.name || (id ? `Отчёт ${id}` : 'Отчёт');

    detailsTitle.textContent = name;

    const fileUrl = normalizeFilePathUrl(report?.file_path);
    const rows = [
      ['report_id', report?.report_id],
      ['author_id', report?.author_id],
      ['status', report?.status],
      ['format', report?.format],
      ['csv_sep', report?.csv_sep],
      ['created_at', report?.created_at ? formatDateTime(report.created_at) : report?.created_at],
      ['updated_at', report?.updated_at ? formatDateTime(report.updated_at) : report?.updated_at],
      ['file_path', fileUrl],
      ['err_msg', report?.err_msg],
      ['comment', report?.comm],
    ].filter(([, v]) => v !== undefined && v !== null && String(v).trim() !== '');

    const grid = rows
      .map(([k, v]) => {
        if (k === 'file_path') {
          const href = String(v || '').trim();
          const isLink = /^https?:\/\//i.test(href)
          const valueHtml = isLink
            ? `<a href="${escapeHtml(href)}" target="_blank" rel="noopener noreferrer">${escapeHtml(href)}</a>`
            : escapeHtml(href);
          return `
            <div class="report-details__key">${escapeHtml(k)}</div>
            <div class="report-details__value">${valueHtml}</div>
          `;
        }
        return `
          <div class="report-details__key">${escapeHtml(k)}</div>
          <div class="report-details__value">${escapeHtml(v)}</div>
        `;
      })
      .join('');

    const query = report?.query;

    detailsBody.innerHTML = `
      <div class="report-details">
        <div class="report-details__grid">${grid || '<div class="report-details__value">Нет данных</div>'}</div>
        ${query ? `<div><div class="report-details__key">query</div><pre class="report-details__pre">${escapeHtml(query)}</pre></div>` : ''}
      </div>
    `;
  }

  function openDetailsModal(report) {
    if (!detailsOverlay) return;
    const id = report?.report_id || report?.reportId || report?.uuid || report?.id || '';
    openReportId = String(id || '');
    renderDetailsModal(report);
    detailsOverlay.style.display = 'flex';
  }

  async function fetchReports({ showLoading = false } = {}) {
    if (activeTab !== 'reports') return;
    if (!isSidebarVisible(sidebar)) return;

    abortController?.abort();
    abortController = new AbortController();

    if (showLoading) {
      loading = true;
      lastError = null;
      renderReports();
    }

    try {
      const payload = await getUserReports({ signal: abortController.signal });
      reports = normalizeReportsResponse(payload);
      reports.sort((a, b) => {
        const ta = new Date(a.updated_at || a.created_at || 0).getTime();
        const tb = new Date(b.updated_at || b.created_at || 0).getTime();
        return (Number.isFinite(tb) ? tb : 0) - (Number.isFinite(ta) ? ta : 0);
      });
      lastError = null;

      if (typeof window.syncHistoryWithReports === 'function') {
        window.syncHistoryWithReports(reports);
      }

      if (openReportId) {
        const updated = reports.find(r => String(r.report_id || r.reportId || r.uuid || r.id || '') === openReportId);
        if (updated && detailsOverlay?.style?.display === 'flex') {
          renderDetailsModal(updated);
        }
      }
    } catch (err) {
      if (err?.name === 'AbortError') return;
      lastError = err?.message ? `Ошибка загрузки отчётов: ${err.message}` : 'Ошибка загрузки отчётов';
    } finally {
      loading = false;
      renderReports();
    }
  }


  function stopPolling() {
    if (pollTimer) {
      clearInterval(pollTimer);
      pollTimer = null;
    }
    abortController?.abort();
    abortController = null;
  }

  function startPolling() {
    if (pollTimer) return;
    pollTimer = setInterval(() => fetchReports({ showLoading: false }), REPORTS_POLL_MS);
  }

  function ensurePolling() {
    const shouldPoll =
      activeTab === 'reports' &&
      document.visibilityState === 'visible' &&
      isSidebarVisible(sidebar);

    if (!shouldPoll) {
      stopPolling();
      return;
    }

    const wasRunning = !!pollTimer;
    startPolling();
    if (!wasRunning) {
      fetchReports({ showLoading: !reports.length });
    }
  }

  tabsRoot.addEventListener('click', (e) => {
    const btn = e.target.closest('.sidebar-tab[data-tab]');
    if (!btn) return;
    setActiveTab(btn.dataset.tab);
  });

  btnReloadReports?.addEventListener('click', () => fetchReports({ showLoading: true }));

  reportsList.addEventListener('click', (e) => {
    const li = e.target.closest('li[data-report-id]');
    if (!li) return;
    const id = li.dataset.reportId || '';
    const report = reports.find(r => String(r.report_id || r.reportId || r.uuid || r.id || '') === id);
    if (!report) return;
    openDetailsModal(report);
  });

  detailsClose?.addEventListener('click', closeDetailsModal);
  detailsOverlay?.addEventListener('click', (e) => {
    if (e.target === detailsOverlay) closeDetailsModal();
  });
  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && detailsOverlay?.style?.display === 'flex') closeDetailsModal();
  });

  const dropdownBtn = document.querySelector('.dropdown-btn');
  dropdownBtn?.addEventListener('click', (e) => {
    e.preventDefault();
    e.stopPropagation();
    sidebar.classList.toggle('show');
    ensurePolling();
  });
  document.addEventListener('click', (e) => {
    if (!sidebar.classList.contains('show')) return;
    if (sidebar.contains(e.target) || dropdownBtn?.contains(e.target)) return;
    sidebar.classList.remove('show');
    ensurePolling();
  });

  window.addEventListener('resize', ensurePolling);
  document.addEventListener('visibilitychange', ensurePolling);

  window.addEventListener('rightPanel:setTab', (e) => {
    setActiveTab(e?.detail?.tab);
  });
  window.addEventListener('rightPanel:refreshReports', () => {
    fetchReports({ showLoading: !reports.length });
  });

  setActiveTab('cache');
}
