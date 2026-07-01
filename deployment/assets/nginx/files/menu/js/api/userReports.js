import { fetchWithToken } from './fetchWithToken.js';
import { API_BASE } from '../config/index.js';

const REPORTS_API = `${API_BASE}/api/v1/reports`;

export async function getUserReports(options = {}) {
  return fetchWithToken(`${REPORTS_API}`, options);
}

export async function getReportStatus(reportId, options = {}) {
  const id = String(reportId || '').trim();
  if (!id) throw new Error('reportId is required');
  return fetchWithToken(`${REPORTS_API}/${encodeURIComponent(id)}/status`, options);
}
