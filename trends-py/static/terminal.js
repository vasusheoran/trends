// ── State ────────────────────────────────────────────────
let activeId = null;
let tickers = {};           // { tickerName: snapshot }
let eventSources = {};      // { tickerName: EventSource }
let connectionCount = 0;

let chart = null;
let candleSeries = null;
let closeSeries  = null;
let ema5Series   = null;
let ema20Series  = null;
let bullishSeries = null;
let supportSeries = null;
let historyView = 'chart';   // 'chart' | 'table'
let historyMode = 'daily';   // 'daily' | 'live'
let historyLiveTF = '1m';    // '1m' | '5m'
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
                ema50:   s.ema50   || 0,
                rsi:     s.rsi     || 0,
                hold:    s.hold    || 0,
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
        { id: `e50-${tickerName}`, key: 'ema50',   fmt: v => v.toFixed(2) },
        { id: `r-${tickerName}`,   key: 'rsi',     fmt: v => v.toFixed(0) },
        { id: `hld-${tickerName}`, key: 'hold',    fmt: v => v.toFixed(2) },
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

    // Push live tick to open chart or table
    if (activeId === tickerName) {
        if (historyMode === 'daily' && historyView === 'chart') {
            const ts = dateToTimestamp(snapshot.date);
            if (ts && closeSeries) {
                closeSeries.update({ time: ts, value: snapshot.close });
                if (ema5Series  && snapshot.ema5    !== null) ema5Series.update({ time: ts, value: snapshot.ema5 });
                if (ema20Series && snapshot.ema20   !== null) ema20Series.update({ time: ts, value: snapshot.ema20 });
                if (bullishSeries && snapshot.bullish !== null) bullishSeries.update({ time: ts, value: snapshot.bullish });
                if (supportSeries && snapshot.support !== null) supportSeries.update({ time: ts, value: snapshot.support });
            }
        } else if (historyMode === 'live') {
            const dayPicker = document.getElementById('day-picker');
            if (!dayPicker || dayPicker.value === '') {
                if (historyView === 'chart' && candleSeries) {
                    const period = historyLiveTF === '1m' ? 60 : 300;
                    const ts = snapshot.timestamp || Math.floor(Date.now() / 1000);
                    const periodTs = Math.floor(ts / period) * period;
                    candleSeries.update({ time: periodTs, open: snapshot.open, high: snapshot.high, low: snapshot.low, close: snapshot.close });
                } else if (historyView === 'table') {
                    if (!currentHistoryData) currentHistoryData = { history: [] };
                    const ts = snapshot.timestamp || Math.floor(Date.now() / 1000);
                    const tick = {
                        time: ts,
                        open: snapshot.open, high: snapshot.high, low: snapshot.low, close: snapshot.close,
                        hl: snapshot.hl, avg: snapshot.avg, ema5: snapshot.ema5, ema20: snapshot.ema20,
                        ema50: snapshot.ema50, rsi: snapshot.rsi, support: snapshot.support,
                    };
                    currentHistoryData.history.push(tick);
                    insertLiveTableRow(tick);
                }
            }
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
            <td id="e50-${tickerName}">${fmt(t.ema50)}</td>
            <td id="r-${tickerName}"   class="fade-out">${fmt(t.rsi, 0)}</td>
            <td id="hld-${tickerName}" class="fade-out">${fmt(t.hold)}</td>
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

    // Reset to Daily mode
    historyMode = 'daily';
    historyLiveTF = '1m';
    document.getElementById('btn-mode-daily').classList.add('active');
    document.getElementById('btn-mode-live').classList.remove('active');
    document.getElementById('daily-controls').style.display = 'flex';
    document.getElementById('live-controls').style.display = 'none';
    document.getElementById('chart-series-controls').style.display = 'flex';
    document.getElementById('month-picker').value = '';
    document.getElementById('day-nav-picker').innerHTML = '<option value="">All</option>';

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

        // Reset month/day filters and populate day picker
        document.getElementById('month-picker').value = '';
        populateDayNavPicker(data.history);

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
    if (historyMode === 'live') {
        const dayPicker = document.getElementById('day-picker');
        loadIntraday(dayPicker?.value || null);
    } else {
        renderCurrentHistoryView();
    }
}

async function setHistoryMode(mode) {
    historyMode = mode;
    const isDaily = mode === 'daily';
    document.getElementById('btn-mode-daily').classList.toggle('active', isDaily);
    document.getElementById('btn-mode-live').classList.toggle('active', !isDaily);
    document.getElementById('daily-controls').style.display = isDaily ? 'flex' : 'none';
    document.getElementById('live-controls').style.display  = isDaily ? 'none' : 'flex';
    document.getElementById('chart-series-controls').style.display = isDaily ? 'flex' : 'none';

    if (isDaily) {
        await loadHistory();
    } else {
        await loadIntraday();
    }
}

async function setLiveTF(tf) {
    historyLiveTF = tf;
    document.getElementById('btn-tf-1m').classList.toggle('active', tf === '1m');
    document.getElementById('btn-tf-5m').classList.toggle('active', tf === '5m');
    await loadIntraday();
}

function onYearChange() {
    document.getElementById('month-picker').value = '';
    document.getElementById('day-nav-picker').innerHTML = '<option value="">All</option>';
    loadHistory();
}

function onMonthChange() {
    if (!currentHistoryData) return;
    const month = parseInt(document.getElementById('month-picker').value) || 0;
    let bars = currentHistoryData.history;
    if (month) {
        bars = bars.filter(b => parseDateParts(b.date).month === month);
    }
    document.getElementById('day-nav-picker').value = '';
    populateDayNavPicker(bars);
    if (historyView === 'chart') renderHistoryChart(bars);
    else renderHistoryTable(bars);
}

function onDayNav(dateStr) {
    if (!dateStr || !chart) return;
    const ts = dateToTimestamp(dateStr);
    if (!ts) return;
    chart.timeScale().setVisibleRange({ from: ts - 7 * 86400, to: ts + 7 * 86400 });
}

function populateDayNavPicker(bars) {
    const picker = document.getElementById('day-nav-picker');
    if (!picker) return;
    const prev = picker.value;
    picker.innerHTML = '<option value="">All</option>';
    bars.forEach(b => {
        if (!b.date) return;
        const opt = document.createElement('option');
        opt.value = b.date;
        opt.textContent = b.date;
        picker.appendChild(opt);
    });
    if (prev) picker.value = prev;
}

async function loadIntraday(date = null) {
    if (!activeId) return;
    const content = document.getElementById('history-content');
    content.className = '';
    content.innerHTML = '<div class="history-loading">LOADING INTRADAY...</div>';

    try {
        // Raw per-second ticks for table view (backend fetches most recent day if no date given)
        const useRaw = historyView === 'table';
        let url = `/api/intraday/${activeId}?tf=${historyLiveTF}`;
        if (date) url += `&date=${encodeURIComponent(date)}`;
        if (useRaw) url += `&raw=true`;

        const res = await fetch(url);
        const raw = await res.json();

        const bars = Array.isArray(raw) ? raw : (raw.bars || []);
        const days = Array.isArray(raw) ? [] : (raw.days || []);

        // Populate day picker
        const picker = document.getElementById('day-picker');
        picker.innerHTML = '<option value="">Live (today)</option>';
        days.forEach(d => {
            const opt = document.createElement('option');
            opt.value = d;
            opt.textContent = d;
            if (d === date) opt.selected = true;
            picker.appendChild(opt);
        });

        // Pre-load existing ticks; SSE will append new ones
        currentHistoryData = { history: bars };
        renderCurrentHistoryView();
    } catch (e) {
        console.error('Intraday load failed:', e);
        content.className = 'idle';
        content.innerHTML = `[ LOAD FAILED: ${e.message} ]`;
    }
}

function onDayChange(date) {
    loadIntraday(date || null);
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

    const formatted = bars
        .map(b => ({
            ...b,
            time: b.timestamp
               || (typeof b.time === 'number' ? b.time : null)
               || dateToTimestamp(b.date)
        }))
        .filter(b => b.time !== null && b.time !== 0)
        .sort((a, b) => a.time - b.time);

    if (historyMode === 'daily') {
        closeSeries = chart.addLineSeries({ color: '#d1d1d1', lineWidth: 1, title: 'Close' });
        ema5Series  = chart.addLineSeries({ color: '#f5a623', lineWidth: 1, title: 'EMA5' });
        ema20Series = chart.addLineSeries({ color: '#4a90e2', lineWidth: 1, title: 'EMA20' });
        bullishSeries = chart.addLineSeries({
            color: '#00ff7f', lineWidth: 1, lineStyle: LightweightCharts.LineStyle.Dashed,
            title: 'Bullish',
        });
        supportSeries = chart.addLineSeries({
            color: '#ff453a', lineWidth: 1, lineStyle: LightweightCharts.LineStyle.Dotted,
            title: 'Support',
        });
        updateSeriesVisibility();

        if (formatted.length) {
            closeSeries.setData(formatted.map(b => ({ time: b.time, value: b.close })));
            ema5Series.setData(formatted.filter(b => b.ema5 != null).map(b => ({ time: b.time, value: b.ema5 })));
            ema20Series.setData(formatted.filter(b => b.ema20 != null).map(b => ({ time: b.time, value: b.ema20 })));
            bullishSeries.setData(formatted.filter(b => b.bullish != null).map(b => ({ time: b.time, value: b.bullish })));
            supportSeries.setData(formatted.filter(b => b.support != null).map(b => ({ time: b.time, value: b.support })));
            chart.timeScale().fitContent();
        }
    } else {
        candleSeries = chart.addCandlestickSeries({
            upColor: '#00ff7f', downColor: '#ff453a', borderVisible: false,
            wickUpColor: '#00ff7f', wickDownColor: '#ff453a',
        });

        if (formatted.length) {
            candleSeries.setData(formatted);
            chart.timeScale().fitContent();
        }
    }
}

function updateSeriesVisibility() {
    if (closeSeries)   closeSeries.applyOptions({ visible: document.getElementById('show-close').checked });
    if (ema5Series)    ema5Series.applyOptions({ visible: document.getElementById('show-ema5').checked });
    if (ema20Series)   ema20Series.applyOptions({ visible: document.getElementById('show-ema20').checked });
    if (bullishSeries) bullishSeries.applyOptions({ visible: document.getElementById('show-bullish').checked });
    if (supportSeries) supportSeries.applyOptions({ visible: document.getElementById('show-support').checked });
}

function renderHistoryTable(bars) {
    destroyChart();
    const content = document.getElementById('history-content');
    content.className = '';

    const isLive = historyMode === 'live';
    const fmt = (v, d=2) => v !== null && v !== undefined ? Number(v).toFixed(d) : '-';

    const barTs = b => b.timestamp || (typeof b.time === 'number' ? b.time : null) || dateToTimestamp(b.date) || 0;
    const barTimeStr = b => {
        const ts = b.timestamp || (typeof b.time === 'number' ? b.time : null);
        if (ts) {
            return isLive
                ? new Date(ts * 1000).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
                : new Date(ts * 1000).toLocaleString();
        }
        return b.date || '-';
    };

    const sorted = [...bars].sort((a, b) => {
        return historySortDir === 'desc' ? barTs(b) - barTs(a) : barTs(a) - barTs(b);
    });

    const icon = historySortDir === 'desc' ? '↓' : '↑';

    let thead, rows;
    if (isLive) {
        thead = `<th style="text-align:left;cursor:pointer" onclick="toggleHistorySort()">Time ${icon}</th>
            <th>Close</th><th>Open</th><th>High</th><th>Low</th>
            <th>H/L</th><th>AVG</th><th>EMA-5</th><th>EMA-20</th><th>EMA-50</th>
            <th>RSI</th><th>Support</th>`;
        rows = sorted.map(b => `
            <tr>
                <td style="text-align:left;color:var(--text-muted)">${barTimeStr(b)}</td>
                <td>${fmt(b.close)}</td><td>${fmt(b.open)}</td>
                <td>${fmt(b.high)}</td><td>${fmt(b.low)}</td>
                <td>${fmt(b.hl)}</td><td>${fmt(b.avg)}</td>
                <td>${fmt(b.ema5)}</td><td>${fmt(b.ema20)}</td><td>${fmt(b.ema50)}</td>
                <td>${fmt(b.rsi, 1)}</td><td>${fmt(b.support)}</td>
            </tr>`).join('');
    } else {
        thead = `<th style="text-align:left;cursor:pointer" onclick="toggleHistorySort()">Date ${icon}</th>
            <th>Close</th><th>Open</th><th>High</th><th>Low</th>
            <th>H/L</th><th>AVG</th><th>EMA-5</th><th>EMA-20</th><th>EMA-50</th>
            <th>RSI</th><th>Hold</th><th>Support</th><th>Bullish</th>`;
        rows = sorted.map(b => `
            <tr>
                <td style="text-align:left;color:var(--text-muted)">${barTimeStr(b)}</td>
                <td>${fmt(b.close)}</td><td>${fmt(b.open)}</td>
                <td>${fmt(b.high)}</td><td>${fmt(b.low)}</td>
                <td>${fmt(b.hl)}</td><td>${fmt(b.avg)}</td>
                <td>${fmt(b.ema5)}</td><td>${fmt(b.ema20)}</td><td>${fmt(b.ema50)}</td>
                <td>${fmt(b.rsi, 1)}</td><td>${fmt(b.hold)}</td>
                <td>${fmt(b.support)}</td>
                <td style="color:var(--accent)">${fmt(b.bullish)}</td>
            </tr>`).join('');
    }

    content.innerHTML = `
        <table class="history-table">
            <thead><tr>${thead}</tr></thead>
            <tbody>${rows}</tbody>
        </table>`;
}

function insertLiveTableRow(tick) {
    const tbody = document.querySelector('#history-content table tbody');
    if (!tbody) return;
    const fmt = (v, d=2) => v !== null && v !== undefined ? Number(v).toFixed(d) : '-';
    const ts = tick.time || tick.timestamp;
    const timeStr = ts
        ? new Date(ts * 1000).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
        : '-';
    const tr = document.createElement('tr');
    tr.innerHTML = `
        <td style="text-align:left;color:var(--text-muted)">${timeStr}</td>
        <td>${fmt(tick.close)}</td><td>${fmt(tick.open)}</td>
        <td>${fmt(tick.high)}</td><td>${fmt(tick.low)}</td>
        <td>${fmt(tick.hl)}</td><td>${fmt(tick.avg)}</td>
        <td>${fmt(tick.ema5)}</td><td>${fmt(tick.ema20)}</td><td>${fmt(tick.ema50)}</td>
        <td>${fmt(tick.rsi, 1)}</td><td>${fmt(tick.support)}</td>`;
    if (historySortDir === 'desc') {
        tbody.insertBefore(tr, tbody.firstChild);
    } else {
        tbody.appendChild(tr);
    }
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
    if (chart) { chart.remove(); chart = null; }
    candleSeries = closeSeries = ema5Series = ema20Series = bullishSeries = supportSeries = null;
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
const _MONTHS = { Jan:1,Feb:2,Mar:3,Apr:4,May:5,Jun:6,
                  Jul:7,Aug:8,Sep:9,Oct:10,Nov:11,Dec:12 };

function parseDateParts(dateStr) {
    if (!dateStr) return { day: 0, month: 0, year: 0 };
    const [d, mon, y] = dateStr.split('-');
    return { day: +d, month: _MONTHS[mon] || 0, year: +y };
}

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
