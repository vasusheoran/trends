// ── State ────────────────────────────────────────────────
let activeId = null;
let tickers = {};           // { tickerName: snapshot }
let eventSources = {};      // { tickerName: EventSource }
let connectionCount = 0;

let chart = null;
let candleSeries = null;
let bullishSeries = null;
let supportSeries = null;
let historyView = 'chart';  // 'chart' | 'table'
let historyTF = 'D';        // 'D' | '1m' | '5m'
let historySortDir = 'desc'; // 'asc' | 'desc'
let currentHistoryData = null;

let resizeObserver = null;

// ── Init ─────────────────────────────────────────────────
async function init() {
    await loadPreferences();
    await fetchTickers();
    render();
    setupSeedUpload();
    setupSplitter();
}

// ── Ticker fetching ───────────────────────────────────────
async function fetchTickers() {
    try {
        const res = await fetch('/api/tickers');
        const data = await res.json();
        for (const t of (data.tickers || [])) {
            if (!tickers[t]) await addTickerToState(t);
        }
    } catch (e) {
        console.error('Failed to fetch tickers:', e);
    }
}

async function addTickerToState(tickerName) {
    tickerName = tickerName.toLowerCase();
    try {
        const res = await fetch(`/api/state/${tickerName}`);
        if (res.ok) {
            const s = await res.json();
            tickers[tickerName] = {
                ticker:  tickerName.toUpperCase(),
                date:    s.date    || 'Seeded',
                close:   s.close   || 0,
                open:    s.open    || 0,
                high:    s.high    || 0,
                low:     s.low     || 0,
                hl:      s.hl      || 0,
                avg:     s.avg     || 0,
                ema5:    s.ema5    || 0,
                ema20:   s.ema20   || 0,
                rsi:     s.rsi     || 0,
                support: s.support || 0,
                bullish: s.bullish || 0,
                warning: null,
            };
        } else {
            tickers[tickerName] = { ticker: tickerName.toUpperCase(), date: 'Pending...' };
        }
        subscribeTicker(tickerName);
    } catch (e) {
        console.error(`Error adding ticker ${tickerName}:`, e);
    }
}

// ── SSE subscription ──────────────────────────────────────
function subscribeTicker(tickerName) {
    if (eventSources[tickerName]) return;
    const source = new EventSource(`/api/stream/${tickerName}`);

    source.onopen = () => {
        connectionCount++;
        updateConnectionStatus();
    };
    source.onmessage = (event) => {
        updateTickerUI(JSON.parse(event.data));
    };
    source.onerror = () => {
        connectionCount = Math.max(0, connectionCount - 1);
        updateConnectionStatus();
        source.close();
        delete eventSources[tickerName];
        setTimeout(() => subscribeTicker(tickerName), 5000);
    };
    eventSources[tickerName] = source;
}

function updateConnectionStatus() {
    const el = document.getElementById('connection-status');
    if (!el) return;
    if (connectionCount > 0) {
        el.textContent = '● CONNECTED';
        el.style.color = 'var(--accent)';
    } else {
        el.textContent = '● DISCONNECTED';
        el.style.color = 'var(--danger)';
    }
}

// ── Real-time UI update ───────────────────────────────────
function updateTickerUI(snapshot) {
    const tickerName = snapshot.ticker.toLowerCase();
    const old = tickers[tickerName] || {};
    tickers[tickerName] = snapshot;

    const row = document.getElementById(`row-${tickerName}`);
    if (!row) { render(); return; }

    const fields = [
        { id: `c-${tickerName}`,   key: 'close',   fmt: v => v.toFixed(2) },
        { id: `o-${tickerName}`,   key: 'open',    fmt: v => v.toFixed(2) },
        { id: `h-${tickerName}`,   key: 'high',    fmt: v => v.toFixed(2) },
        { id: `l-${tickerName}`,   key: 'low',     fmt: v => v.toFixed(2) },
        { id: `hl-${tickerName}`,  key: 'hl',      fmt: v => v.toFixed(2) },
        { id: `avg-${tickerName}`, key: 'avg',     fmt: v => v.toFixed(2) },
        { id: `e5-${tickerName}`,  key: 'ema5',    fmt: v => v.toFixed(2) },
        { id: `e20-${tickerName}`, key: 'ema20',   fmt: v => v.toFixed(2) },
        { id: `r-${tickerName}`,   key: 'rsi',     fmt: v => v.toFixed(0) },
        { id: `s-${tickerName}`,   key: 'support', fmt: v => v.toFixed(2) },
        { id: `b-${tickerName}`,   key: 'bullish', fmt: v => v.toFixed(2) },
    ];

    fields.forEach(f => {
        const cell = document.getElementById(f.id);
        if (!cell) return;
        const newVal = snapshot[f.key];
        const oldVal = old[f.key];
        if (newVal !== oldVal) {
            cell.innerText = typeof newVal === 'number' ? f.fmt(newVal) : (newVal || '-');
            if (oldVal !== undefined && typeof newVal === 'number') {
                const cls = newVal > oldVal ? 'flash-up' : 'flash-down';
                cell.classList.remove('fade-out');
                cell.classList.add(cls);
                setTimeout(() => { cell.classList.add('fade-out'); cell.classList.remove(cls); }, 50);
            }
        }
    });

    document.getElementById(`d-${tickerName}`).innerText = snapshot.date || '-';

    // Push live tick to open chart
    if (activeId === tickerName && candleSeries && historyView === 'chart') {
        if (historyTF === 'D') {
            const ts = dateToTimestamp(snapshot.date);
            if (ts) {
                candleSeries.update({ time: ts, open: snapshot.open, high: snapshot.high, low: snapshot.low, close: snapshot.close });
                if (bullishSeries && snapshot.bullish !== null) bullishSeries.update({ time: ts, value: snapshot.bullish });
                if (supportSeries && snapshot.support !== null) supportSeries.update({ time: ts, value: snapshot.support });
            }
        } else {
            const period = historyTF === '1m' ? 60 : 300;
            const ts = snapshot.timestamp || Math.floor(Date.now() / 1000);
            const periodTs = Math.floor(ts / period) * period;
            candleSeries.update({ 
                time: periodTs, 
                open: snapshot.open, 
                high: snapshot.high, 
                low: snapshot.low, 
                close: snapshot.close 
            });
        }
    }
}

// ── Grid render ───────────────────────────────────────────
function render() {
    const tbody = document.getElementById('ticker-body');
    if (!tbody) return;
    tbody.innerHTML = '';

    Object.keys(tickers).sort().forEach(tickerName => {
        const t = tickers[tickerName];
        const tr = document.createElement('tr');
        tr.id = `row-${tickerName}`;
        if (activeId === tickerName) tr.className = 'selected';
        tr.onclick = (e) => { e.stopPropagation(); toggleSelection(tickerName); };

        const fmt = (v, d=2) => typeof v === 'number' ? v.toFixed(d) : '-';
        tr.innerHTML = `
            <td class="ticker-label">${t.ticker}</td>
            <td id="d-${tickerName}"   style="color:var(--text-muted)">${t.date || '-'}</td>
            <td id="c-${tickerName}"   class="fade-out">${fmt(t.close)}</td>
            <td id="o-${tickerName}">${fmt(t.open)}</td>
            <td id="h-${tickerName}">${fmt(t.high)}</td>
            <td id="l-${tickerName}">${fmt(t.low)}</td>
            <td id="hl-${tickerName}">${fmt(t.hl)}</td>
            <td id="avg-${tickerName}">${fmt(t.avg)}</td>
            <td id="e5-${tickerName}"  class="fade-out">${fmt(t.ema5)}</td>
            <td id="e20-${tickerName}">${fmt(t.ema20)}</td>
            <td id="r-${tickerName}"   class="fade-out">${fmt(t.rsi, 0)}</td>
            <td id="s-${tickerName}"   class="fade-out">${fmt(t.support)}</td>
            <td id="b-${tickerName}"   class="fade-out" style="color:var(--accent)">${fmt(t.bullish)}</td>
            <td id="a-${tickerName}">${t.warning ? '⚠' : '-'}</td>
        `;
        tbody.appendChild(tr);
    });
    updateButtons();
}

function toggleSelection(id) {
    closeHistory();
    activeId = (activeId === id) ? null : id;
    render();
}

function handleGlobalClick(e) {
    // Selection logic moved entirely to toggleSelection (row click)
    // Global click no longer deselects to prevent accidental closing.
}

function updateButtons() {
    const on = activeId !== null;
    document.getElementById('btn-hist').disabled = !on;
    document.getElementById('btn-del').disabled  = !on;
}

function setupSplitter() {
    const splitter = document.getElementById('splitter');
    const history = document.getElementById('chart-area');
    let isDragging = false;

    splitter.onmousedown = (e) => {
        isDragging = true;
        document.body.style.cursor = 'ns-resize';
        e.preventDefault();
    };

    document.onmousemove = (e) => {
        if (!isDragging) return;
        const totalHeight = window.innerHeight;
        const newHistoryHeight = totalHeight - e.clientY;
        
        if (newHistoryHeight > 40 && newHistoryHeight < totalHeight - 150) {
            history.style.height = `${newHistoryHeight}px`;
        }
    };

    document.onmouseup = () => {
        if (isDragging) {
            isDragging = false;
            document.body.style.cursor = 'default';
        }
    };
}

async function toggleTheme() {
    const isLight = document.body.classList.toggle('light-mode');
    const theme = isLight ? 'light' : 'dark';
    try {
        await fetch('/api/preferences', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({theme})
        });
    } catch (e) { console.error('Failed to save theme preference', e); }
    if (chart) renderCurrentHistoryView();
}

async function loadPreferences() {
    try {
        const res = await fetch('/api/preferences');
        if (res.ok) {
            const prefs = await res.json();
            if (prefs.theme === 'light') {
                document.body.classList.add('light-mode');
            }
        }
    } catch (e) { console.error('Failed to load preferences', e); }
}

// ── History panel ─────────────────────────────────────────
async function openHistory() {
    if (!activeId) return;
    destroyChart();

    document.getElementById('splitter').style.display = 'block';
    const panel = document.getElementById('chart-area');
    panel.style.height = '220px'; 
    panel.classList.add('expanded');
    document.getElementById('history-header').style.display = 'flex';
    document.getElementById('history-title').textContent = `${activeId.toUpperCase()} · History`;

    // Reset TF to Daily
    historyTF = 'D';
    document.getElementById('tf-picker').value = 'D';
    document.getElementById('year-picker-container').style.display = 'flex';
    document.getElementById('chart-series-controls').style.display = 'flex';

    const defaultYear = extractYear(tickers[activeId]?.date) || new Date().getFullYear();
    const picker = document.getElementById('year-picker');
    picker.innerHTML = `<option value="${defaultYear}">${defaultYear}</option>`;

    await loadHistory();
}

async function loadHistory() {
    if (!activeId) return;
    const picker = document.getElementById('year-picker');
    let year = parseInt(picker.value);
    
    const content = document.getElementById('history-content');
    content.className = '';
    content.innerHTML = '<div class="history-loading">LOADING...</div>';

    try {
        const res = await fetch(`/api/history/${activeId}${year ? `?year=${year}` : ''}`);
        if (!res.ok) {
            const txt = await res.text();
            throw new Error(`HTTP ${res.status}: ${txt.slice(0, 300)}`);
        }
        const data = await res.json();
        currentHistoryData = data;
        
        const years = data.years || [];
        populateYearPicker(years, year);

        // If we requested a year but got no bars, and there ARE years available, 
        // retry with the latest year.
        if (data.history.length === 0 && years.length > 0 && year !== years[0]) {
            picker.value = years[0];
            return await loadHistory();
        }

        renderCurrentHistoryView();
    } catch (e) {
        console.error('History load failed:', e);
        content.className = 'idle';
        content.innerHTML = `[ LOAD FAILED: ${e.message} ]`;
    }
}

function setHistoryView(mode) {
    historyView = mode;
    document.getElementById('btn-view-chart').classList.toggle('active', mode === 'chart');
    document.getElementById('btn-view-table').classList.toggle('active', mode === 'table');
    renderCurrentHistoryView();
}

async function setHistoryTF(tf) {
    historyTF = tf;
    // Show/hide year picker and chart series controls based on TF
    document.getElementById('year-picker-container').style.display = tf === 'D' ? 'flex' : 'none';
    document.getElementById('chart-series-controls').style.display = tf === 'D' ? 'flex' : 'none';
    
    if (tf === 'D') {
        await loadHistory();
    } else {
        await loadIntraday();
    }
}

async function loadIntraday() {
    if (!activeId) return;
    const content = document.getElementById('history-content');
    content.className = '';
    content.innerHTML = '<div class="history-loading">LOADING INTRADAY...</div>';

    try {
        const res = await fetch(`/api/intraday/${activeId}?tf=${historyTF}`);
        const data = await res.json();
        currentHistoryData = { history: data };
        renderCurrentHistoryView();
    } catch (e) {
        console.error('Intraday load failed:', e);
        content.className = 'idle';
        content.innerHTML = `[ LOAD FAILED: ${e.message} ]`;
    }
}

function toggleHistorySort() {
    historySortDir = historySortDir === 'desc' ? 'asc' : 'desc';
    renderCurrentHistoryView();
}

function renderCurrentHistoryView() {
    if (!currentHistoryData) return;
    if (historyView === 'chart') {
        renderHistoryChart(currentHistoryData.history);
    } else {
        renderHistoryTable(currentHistoryData.history);
    }
}

function populateYearPicker(years, selectedYear) {
    const picker = document.getElementById('year-picker');
    picker.innerHTML = years.map(y =>
        `<option value="${y}"${y === selectedYear ? ' selected' : ''}>${y}</option>`
    ).join('');
}

function renderHistoryChart(bars) {
    destroyChart();
    const content = document.getElementById('history-content');
    content.className = '';
    content.innerHTML = '';

    const isLight = document.body.classList.contains('light-mode');
    const colors = isLight ?
        { bg: '#f5f5f7', text: '#1d1d1f', grid: '#d2d2d7' } :
        { bg: '#080808', text: '#d1d1d1', grid: '#111' };

    chart = LightweightCharts.createChart(content, {
        layout: { backgroundColor: colors.bg, textColor: colors.text },
        grid: { vertLines: { color: colors.grid }, horzLines: { color: colors.grid } },
        timeScale: { borderColor: colors.grid, timeVisible: true, secondsVisible: false },
        crosshair: { mode: LightweightCharts.CrosshairMode.Normal },
        width: content.clientWidth,
        height: content.clientHeight || 220,
    });

    resizeObserver = new ResizeObserver(entries => {
        if (chart && entries[0].contentRect) {
            const { width, height } = entries[0].contentRect;
            chart.resize(width, height);
        }
    });
    resizeObserver.observe(content);

    candleSeries = chart.addCandlestickSeries({
        upColor: '#00ff7f', downColor: '#ff453a', borderVisible: false,
        wickUpColor: '#00ff7f', wickDownColor: '#ff453a',
    });

    if (historyTF === 'D') {
        bullishSeries = chart.addLineSeries({
            color: '#00ff7f', lineWidth: 1, lineStyle: LightweightCharts.LineStyle.Dashed,
            title: 'Bullish',
        });
        supportSeries = chart.addLineSeries({
            color: '#ff453a', lineWidth: 1, lineStyle: LightweightCharts.LineStyle.Dotted,
            title: 'Support',
        });
        updateSeriesVisibility();
    }

    const formatted = bars
        .map(b => ({ ...b, time: b.timestamp || dateToTimestamp(b.date) }))
        .filter(b => b.time !== null)
        .sort((a, b) => a.time - b.time);

    if (formatted.length) {
        candleSeries.setData(formatted);

        if (historyTF === 'D') {
            const bullishData = formatted
                .filter(b => b.bullish !== null)
                .map(b => ({ time: b.time, value: b.bullish }));
            bullishSeries.setData(bullishData);

            const supportData = formatted
                .filter(b => b.support !== null)
                .map(b => ({ time: b.time, value: b.support }));
            supportSeries.setData(supportData);
        }

        chart.timeScale().fitContent();
    }
}

function updateSeriesVisibility() {
    if (!bullishSeries || !supportSeries) return;
    const showBullish = document.getElementById('show-bullish').checked;
    const showSupport = document.getElementById('show-support').checked;
    bullishSeries.applyOptions({ visible: showBullish });
    supportSeries.applyOptions({ visible: showSupport });
}

function renderHistoryTable(bars) {
    destroyChart();
    const content = document.getElementById('history-content');
    content.className = '';

    const fmt = (v, d=2) => v !== null && v !== undefined ? Number(v).toFixed(d) : '-';

    // Sort logic
    const sorted = [...bars].sort((a, b) => {
        const tA = a.timestamp || dateToTimestamp(a.date) || 0;
        const tB = b.timestamp || dateToTimestamp(b.date) || 0;
        return historySortDir === 'desc' ? tB - tA : tA - tB;
    });

    const rows = sorted.map(b => {
        const dateStr = b.timestamp ? new Date(b.timestamp * 1000).toLocaleString() : b.date;
        return `
        <tr>
            <td style="text-align:left;color:var(--text-muted)">${dateStr}</td>
            <td>${fmt(b.close)}</td>
            <td>${fmt(b.open)}</td>
            <td>${fmt(b.high)}</td>
            <td>${fmt(b.low)}</td>
            <td>${fmt(b.hl)}</td>
            <td>${fmt(b.avg)}</td>
            <td>${fmt(b.ema5)}</td>
            <td>${fmt(b.ema20)}</td>
            <td>${fmt(b.rsi, 1)}</td>
            <td>${fmt(b.support)}</td>
            <td style="color:var(--accent)">${fmt(b.bullish)}</td>
        </tr>`}).join('');

    const icon = historySortDir === 'desc' ? '↓' : '↑';
    content.innerHTML = `
        <table class="history-table">
            <thead>
                <tr>
                    <th style="text-align:left; cursor:pointer" onclick="toggleHistorySort()">Date ${icon}</th>
                    <th>Close</th><th>Open</th><th>High</th><th>Low</th>
                    <th>H/L</th><th>AVG</th><th>EMA-5</th><th>EMA-20</th>
                    <th>RSI</th><th>Support</th><th>Bullish</th>
                </tr>
            </thead>
            <tbody>${rows}</tbody>
        </table>`;
}

function closeHistory() {
    destroyChart();
    document.getElementById('splitter').style.display = 'none';
    const panel = document.getElementById('chart-area');
    panel.classList.remove('expanded');
    panel.style.height = '0';
    document.getElementById('history-header').style.display = 'none';
    const content = document.getElementById('history-content');
    content.className = 'idle';
    content.innerHTML = '[ CHART_ENGINE_IDLE ]';
    currentHistoryData = null;
}

function destroyChart() {
    if (resizeObserver) { resizeObserver.disconnect(); resizeObserver = null; }
    if (chart) { chart.remove(); chart = null; candleSeries = null; }
}

// ── Seed upload ───────────────────────────────────────────
function triggerSeedUpload() {
    let target = activeId;
    if (!target) {
        const name = prompt('Ticker name to seed (e.g. NIFTY):');
        if (!name || !name.trim()) return;
        target = name.trim().toLowerCase();
    }
    const input = document.getElementById('fileInput');
    input.dataset.seedTarget = target;
    input.value = '';
    input.click();
}

function setupSeedUpload() {
    const input = document.getElementById('fileInput');
    input.onchange = async () => {
        const target = input.dataset.seedTarget;
        const file = input.files[0];
        if (!file || !target) return;

        const formData = new FormData();
        formData.append('file', file);
        try {
            const res = await fetch(`/api/seed/${target}`, { method: 'POST', body: formData });
            const result = await res.json();
            if (res.ok) {
                alert(`SUCCESS: Seeded ${result.bars_loaded} bars for ${target.toUpperCase()}\nDetection: ${result.column_detection}`);
                await addTickerToState(target);
                render();
            } else {
                alert(`UPLOAD FAILED: ${result.detail || 'Unknown error'}`);
            }
        } catch (err) {
            console.error('Seed upload failed:', err);
            alert('NETWORK ERROR: Could not connect to server.');
        }
        input.value = '';
    };
}

// ── Ticker management ─────────────────────────────────────
async function removeTicker() {
    if (!activeId) return;
    const t = activeId;
    if (!confirm(`Delete ${t.toUpperCase()}?`)) return;
    try {
        const res = await fetch(`/api/tickers/${t}`, { method: 'DELETE' });
        if (res.ok) {
            if (eventSources[t]) { eventSources[t].close(); delete eventSources[t]; }
            delete tickers[t];
            activeId = null;
            closeHistory();
            render();
        }
    } catch (e) {
        console.error(`Failed to delete ticker ${t}:`, e);
    }
}

function addTicker() {
    const n = prompt('Ticker Name:');
    if (n) addTickerToState(n.toLowerCase());
}

// ── Helpers ───────────────────────────────────────────────
function dateToTimestamp(t) {
    if (!t) return null;
    let iso = null;
    if (/^\d{4}-\d{2}-\d{2}$/.test(t)) {
        iso = t;
    } else {
        const months = { Jan:'01',Feb:'02',Mar:'03',Apr:'04',May:'05',Jun:'06',
                         Jul:'07',Aug:'08',Sep:'09',Oct:'10',Nov:'11',Dec:'12' };
        const parts = t.split('-');
        if (parts.length === 3 && months[parts[1]]) {
            iso = `${parts[2]}-${months[parts[1]]}-${parts[0].padStart(2, '0')}`;
        }
    }
    if (!iso) return null;
    return Math.floor(new Date(iso + 'T00:00:00Z').getTime() / 1000);
}

function extractYear(dateStr) {
    if (!dateStr) return null;
    // "20-Dec-2024" → 2024
    const parts = dateStr.split('-');
    if (parts.length === 3) return parseInt(parts[2]);
    return null;
}

document.addEventListener('DOMContentLoaded', init);
