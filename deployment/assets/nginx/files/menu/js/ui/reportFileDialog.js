import { getReportComment, getReportInfo } from '../api/userReports.js';

const STATUS_LABELS = {
    CREATED: 'В очереди',
    RUNNING: 'Формируется',
    COMPLETED: 'Готов',
    FAILED: 'Ошибка формирования',
    ARCHIVED: 'Недоступен',
};
const REPORT_STATUS_POLL_MS = 2000;

export function reportStatusLabel(status) {
    return Object.hasOwn(STATUS_LABELS, status) ? STATUS_LABELS[status] : 'Проверить статус';
}

export function isReportPending(status) {
    return status === 'CREATED' || status === 'RUNNING';
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
    const format = document.getElementById('reportFileFormat');
    const status = document.getElementById('reportFileStatus');
    const closeButton = document.getElementById('reportFileClose');
    const download = document.getElementById('reportFileDownload');
    let controller;
    let pollTimer;
    let opener;
    let link;

    function showMetadata(report) {
        title.textContent = report.name || 'Файл отчёта';
        comment.textContent = getReportComment(report) || '';
        comment.hidden = !comment.textContent;
        format.textContent = report.format ? `Формат: ${report.format}` : '';
        format.hidden = !format.textContent;
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
        if (opener?.isConnected) opener.focus();
    }

    async function refresh(reportId, signal) {
        try {
            const { report } = await getReportInfo(reportId, { signal, cache: 'no-store' });
            if (signal.aborted) return;
            showMetadata(report);

            if (isReportPending(report.status)) {
                showStatus(reportStatusLabel(report.status));
                pollTimer = setTimeout(() => refresh(reportId, signal), REPORT_STATUS_POLL_MS);
                return;
            }

            if (report.status === 'COMPLETED') {
                const fileLink = parseDownloadLink(report.file_path);
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
            const first = closeButton;
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

    return (button, name) => {
        stopPolling();
        controller = new AbortController();
        opener = button;
        showMetadata({ name });
        showStatus('Проверяем статус…');
        dialog.hidden = false;
        closeButton.focus();
        refresh(button.dataset.reportId, controller.signal);
    };
}
