import { getReportComment, getReportInfo } from '../api/userReports.js';

const STATUS_LABELS = {
    CREATED: 'В очереди',
    RUNNING: 'Формируется',
    COMPLETED: 'Готов',
    FAILED: 'Ошибка формирования',
    ARCHIVED: 'Недоступен',
};
const REPORT_STATUS_POLL_MS = 2000;
const reportDateFormat = new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit',
});

export function reportStatusLabel(status) {
    return Object.hasOwn(STATUS_LABELS, status) ? STATUS_LABELS[status] : 'Проверить статус';
}

export function isReportPending(status) {
    return status === 'CREATED' || status === 'RUNNING';
}

function formatTimestamp(value) {
    const date = new Date(value ?? NaN);
    return date.getUTCFullYear() > 1 ? reportDateFormat.format(date) : '';
}

function parseDownloadLink(filePath) {
    if (!filePath) return null;
    try {
        const url = new URL(filePath);
        if (url.protocol !== 'https:' && url.protocol !== 'http:') return null;

        const signingDate = url.searchParams.get('X-Amz-Date');
        const lifetime = Number(url.searchParams.get('X-Amz-Expires'));
        const parts = signingDate?.match(/^(\d{4})(\d{2})(\d{2})T(\d{2})(\d{2})(\d{2})Z$/);
        const expiresAt = parts && lifetime > 0
            ? Date.parse(`${parts[1]}-${parts[2]}-${parts[3]}T${parts[4]}:${parts[5]}:${parts[6]}Z`) + lifetime * 1000
            : Infinity;
        return { href: filePath, expiresAt };
    } catch {
        return null;
    }
}

export function initReportFileDialog() {
    const dialog = document.getElementById('reportFileDialog');
    const title = document.getElementById('reportFileTitle');
    const comment = document.getElementById('reportFileComment');
    const metadata = document.getElementById('reportFileMetadata');
    const metadataBody = metadata.tBodies[0];
    const details = document.getElementById('reportFileDetails');
    const detailsSummary = details.querySelector('summary');
    const idRow = document.getElementById('reportFileIdRow');
    const idValue = idRow.querySelector('td');
    const queryRow = document.getElementById('reportFileQueryRow');
    const queryCode = queryRow.querySelector('code');
    const status = document.getElementById('reportFileStatus');
    const closeButton = document.getElementById('reportFileClose');
    const download = document.getElementById('reportFileDownload');
    let controller;
    let pollTimer;
    let opener;
    let link;
    let resolveQueuedReport;

    function showMetadata(report, fileLink) {
        title.textContent = report.name || 'Файл отчёта';
        comment.textContent = getReportComment(report) || '';
        comment.hidden = !comment.textContent;

        const rows = [
            ['Статус', Object.hasOwn(STATUS_LABELS, report.status) ? reportStatusLabel(report.status) : report.status],
            ['Формат', report.format],
            ['Разделитель CSV', report.format?.toUpperCase() === 'CSV' ? report.csv_sep : ''],
            ['Создан', formatTimestamp(report.created_at)],
            ['Обновлён', formatTimestamp(report.updated_at)],
            ['Ссылка действует до', formatTimestamp(fileLink?.expiresAt)],
            ['Ошибка', report.status === 'FAILED' ? report.err_msg : ''],
        ];

        const content = document.createDocumentFragment();
        for (const [label, value] of rows) {
            if (!value?.trim()) continue;
            const row = document.createElement('tr');
            const heading = document.createElement('th');
            const cell = document.createElement('td');
            heading.scope = 'row';
            heading.textContent = label;
            cell.textContent = value;
            row.append(heading, cell);
            content.append(row);
        }
        metadata.hidden = !content.childElementCount;
        metadataBody.replaceChildren(content);

        idValue.textContent = report.report_id || '';
        idRow.hidden = !idValue.textContent.trim();
        queryCode.textContent = report.query || '';
        queryRow.hidden = !queryCode.textContent.trim();
        details.hidden = idRow.hidden && queryRow.hidden;
    }

    function showStatus(message, fileLink = null) {
        if (status.textContent !== message) status.textContent = message;
        link = fileLink;
        download.hidden = !link;
        if (link) download.href = link.href;
        else download.removeAttribute('href');
    }

    function stopPolling() {
        clearTimeout(pollTimer);
        controller?.abort();
    }

    function close() {
        stopPolling();
        dialog.hidden = true;
        const button = opener;
        const restoreFocus = () => {
            if (dialog.hidden && button?.isConnected) button.focus();
        };
        if (resolveQueuedReport) {
            resolveQueuedReport();
            resolveQueuedReport = undefined;
            // The submit handler reenables its button after acknowledging the result.
            queueMicrotask(restoreFocus);
        } else {
            restoreFocus();
        }
    }

    async function refresh(reportId, signal) {
        try {
            const { report } = await getReportInfo(reportId, { signal, cache: 'no-store' });
            if (signal.aborted) return;
            const fileLink = report.status === 'COMPLETED' ? parseDownloadLink(report.file_path) : null;
            showMetadata(report, fileLink);

            if (isReportPending(report.status)) {
                showStatus(reportStatusLabel(report.status));
                pollTimer = setTimeout(() => refresh(reportId, signal), REPORT_STATUS_POLL_MS);
                return;
            }

            if (report.status === 'COMPLETED') {
                if (!fileLink) {
                    showStatus('Файл недоступен.');
                } else if (Date.now() >= fileLink.expiresAt) {
                    showStatus('Срок действия ссылки истёк.');
                } else {
                    showStatus('Отчёт готов к скачиванию.', fileLink);
                }
                return;
            }

            showStatus(report.status === 'FAILED' ? 'Не удалось сформировать отчёт.' : 'Отчёт недоступен.');
        } catch {
            if (!signal.aborted) showStatus('Не удалось проверить статус. Закройте окно и попробуйте ещё раз.');
        }
    }

    closeButton.addEventListener('click', close);
    dialog.addEventListener('click', event => {
        if (event.target === dialog) close();
    });
    dialog.addEventListener('keydown', event => {
        if (event.key === 'Escape') {
            close();
        } else if (event.key === 'Tab') {
            const first = details.hidden ? closeButton : detailsSummary;
            const last = download.hidden ? closeButton : download;
            if (document.activeElement === (event.shiftKey ? first : last)) {
                event.preventDefault();
                (event.shiftKey ? last : first).focus();
            }
        }
    });
    download.addEventListener('click', event => {
        if (link && Date.now() < link.expiresAt) return;
        event.preventDefault();
        showStatus('Срок действия ссылки истёк.');
        closeButton.focus();
    });

    function openDialog(button, report, message) {
        if (!dialog.hidden) close();
        opener = button;
        details.open = false;
        showMetadata(report);
        showStatus(message);
        dialog.hidden = false;
        closeButton.focus();
    }

    function openReportFile(button, name) {
        openDialog(button, { name }, 'Проверяем статус…');
        controller = new AbortController();
        refresh(button.dataset.reportId, controller.signal);
    }

    function showQueuedReport(button, report) {
        openDialog(button, report, 'Отчёт поставлен в очередь.');
        return new Promise(resolve => {
            resolveQueuedReport = resolve;
        });
    }

    return { openReportFile, showQueuedReport };
}
