import { fetchWithToken } from './fetchWithToken.js';
import { API_BASE } from '../config/index.js';

const REPORTS_API = `${API_BASE}/api/v1/reports`;

const DONE_STATUSES = new Set(["done", "success", "completed", "ready", "finished", "ok"]);
const FAIL_STATUSES = new Set(["failed", "error", "err", "canceled", "cancelled"]);

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function normalizeStatus(value) {
    return String(value || "").trim().toLowerCase();
}

function normalizeFileUrl(path) {
    const raw = String(path || "").trim();
    if (!raw) return "";
    try {
        const u = new URL(raw, API_BASE);
        if (u.hostname === "report-minio" || u.hostname.includes("minio")) {
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

async function createReportTask(payload) {
    const res = await fetchWithToken(`${REPORTS_API}`, {
        method: "POST",
        headers: { "Content-Type": "application/json", "Accept": "application/json" },
        body: JSON.stringify(payload),
    });

    const reportId = res?.report_id ?? res?.reportId ?? res?.id ?? res?.uuid ?? "";
    const status = res?.status ?? res?.Status ?? "";

    return {
        reportId: reportId ? String(reportId) : "",
        status: status ? String(status) : "",
        raw: res
    };
}

async function waitForReport(reportId, { timeoutMs = 120000, intervalMs = 2000 } = {}) {
    const started = Date.now();
    let lastStatus = "";

    while (Date.now() - started < timeoutMs) {
        const statusResp = await fetchWithToken(`${REPORTS_API}/${encodeURIComponent(reportId)}/status`);
        lastStatus = statusResp?.status ?? statusResp?.Status ?? "";
        const normalized = normalizeStatus(lastStatus);

        if (FAIL_STATUSES.has(normalized)) {
            throw new Error(`Report failed (${lastStatus || "failed"})`);
        }
        if (DONE_STATUSES.has(normalized)) {
            return { status: lastStatus };
        }

        await sleep(intervalMs);
    }

    throw new Error(`Timeout waiting for report ${reportId} (${lastStatus || "pending"})`);
}

async function getReportInfo(reportId) {
    return fetchWithToken(`${REPORTS_API}/${encodeURIComponent(reportId)}`);
}

export async function postReportAndCreateTask({
    format = "PDF",
    sql = "",
    reportName,
    reportComment,
    csvSep,
    createdAt
} = {}) {
    const formatValue = (format || "PDF");
    const fileFormat = String(formatValue).toLowerCase();
    const sepSource = csvSep ?? document.getElementById('csvSeparator')?.value ?? ",";
    const normalizedSep = (typeof sepSource === "string" && sepSource.length) ? sepSource : ",";
    const resolvedName = (reportName ?? document.getElementById('reportName')?.value ?? "").trim();
    const resolvedComment = (reportComment ?? document.getElementById('reportComment')?.value ?? "").trim();

    const payload = {
        name: resolvedName,
        comment: resolvedComment,
        query_sql: (sql || "").trim(),
        format: String(formatValue),
    };
    if (fileFormat === "csv" && normalizedSep) {
        payload.csv_separator = normalizedSep;
    }

    const task = await createReportTask(payload);

    return {
        format: fileFormat,
        task: {
            uuid: task.reportId,
            status: task.status,
        },
        raw: task.raw
    };
}

export async function postReportAndGetBlob({
    format = "PDF",
    sql = "",
    reportName,
    reportComment,
    csvSep,
    createdAt
} = {}) {
    const fileFormat = (format || "PDF").toLowerCase(); // 'pdf' | 'csv' | 'xlsx' | 'json' | 'chart'
    const sepSource = csvSep ?? document.getElementById('csvSeparator')?.value ?? ",";
    const normalizedSep = (typeof sepSource === "string" && sepSource.length) ? sepSource : ",";
    const resolvedName = (reportName ?? document.getElementById('reportName')?.value ?? "").trim();
    const resolvedComment = (reportComment ?? document.getElementById('reportComment')?.value ?? "").trim();

    const payload = {
        name: resolvedName,
        comment: resolvedComment,
        query_sql: (sql || "").trim(),
        format: String(format).toUpperCase(),
    };
    if (fileFormat === "csv" && normalizedSep) {
        payload.csv_separator = normalizedSep;
    }

    const acceptByFormat = {
        pdf:  "application/pdf",
        csv:  "text/csv",
        xlsx: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet, application/zip",
        json: "application/json",
        chart: "application/json",
        docx: "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
    };
    const accept = acceptByFormat[fileFormat] || "*/*";

    const task = await createReportTask(payload);
    if (!task.reportId) {
        throw new Error("Report task creation failed");
    }

    await waitForReport(task.reportId);

    const info = await getReportInfo(task.reportId);
    const report = info?.report ?? info?.Report ?? info;
    const filePath = report?.file_path ?? report?.filePath ?? "";
    if (!filePath) {
        throw new Error("Report file path is missing");
    }

    const fileUrl = normalizeFileUrl(filePath);
    const res = await fetchWithToken(fileUrl, {
        method: "GET",
        headers: { "Accept": accept },
    });

    const expectsJson = fileFormat === "json" || fileFormat === "chart";
    if (expectsJson && res && typeof res !== "string" && typeof res.blob !== "function") {
        const jsonObj = res;
        const blob = new Blob([JSON.stringify(jsonObj, null, 2)], { type: "application/json" });
        const filename = fileFormat === "chart" ? "chart.json" : "report.json";
        return { blob, filename, format: fileFormat, json: jsonObj };
    }

    if (!(res && typeof res.blob === "function")) {
        const details = typeof res === "string" ? res : JSON.stringify(res);
        throw new Error(`Unexpected report response for ${fileFormat}: ${details?.slice?.(0, 300) || ""}`);
    }

    const blob = await res.blob();

    const fallbackNameByFormat = {
        pdf:  "report.pdf",
        csv:  "report.csv",
        xlsx: "report.xlsx",
        json: "report.json",
        chart: "chart.json",
        docx: "report.docx"
    };
    const filename = pickFilename(res.headers, fallbackNameByFormat[fileFormat] || `report.${fileFormat}`);

    return { blob, filename, format: fileFormat };
}


export function saveBlob(blob, filename) {
    const a = document.createElement('a');
    const url = URL.createObjectURL(blob);
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    a.remove();
    setTimeout(() => URL.revokeObjectURL(url), 1000);
}


export function openBlob(blob) {
    const url = URL.createObjectURL(blob);
    window.open(url, '_blank');
    setTimeout(() => URL.revokeObjectURL(url), 60000);
}


export function pickFilename(headers, fallback) {
    const cd = headers.get('Content-Disposition') || headers.get('content-disposition') || '';
    let m = cd.match(/filename\*=(?:UTF-8'')?([^;]+)/i);
    if (m && m[1]) {
        try { return decodeURIComponent(m[1].replace(/(^"|"$)/g, '')); } catch {}
        return m[1].replace(/(^"|"$)/g, '');
    }
    m = cd.match(/filename="?([^"]+)"?/i);
    if (m && m[1]) return m[1];
    return fallback;
}


export async function deleteReport(reportId) {
	if (!reportId) {
		throw new Error('Report ID is required')
	}

	const url = `${REPORTS_API}/${encodeURIComponent(reportId)}`
	const response = await fetchWithToken(url, {
		method: 'DELETE',
	})

	// ✅ Проверяем ошибки более мягко
	if (response && (response.error || response.err)) {
		const errMsg = response.err || response.error || 'Unknown error'
		throw new Error(errMsg)
	}

	return response
}


export async function deleteAllReports() {
    return await fetchWithToken(REPORTS_API, {
        method: 'DELETE',
    });
}
