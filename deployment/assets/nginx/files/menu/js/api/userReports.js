import { fetchWithToken } from './fetchWithToken.js';
import { API_BASE } from '../config/index.js';

const REPORTS_API = `${API_BASE}/api/v1`;

export async function getUserReports(options = {}) {
  return fetchWithToken(`${REPORTS_API}/report`, options);
}

export async function getReportStatus(reportId, options = {}) {
  const id = String(reportId || '').trim();
  if (!id) throw new Error('reportId is required');
  return fetchWithToken(`${REPORTS_API}/report/${encodeURIComponent(id)}`, options);
}
