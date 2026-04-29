let activeId = null;
let tickers = {}; // { tickerName: snapshot }
let eventSources = {}; // { tickerName: EventSource }
let chart = null;
let candleSeries = null;

async function init() {
    await fetchTickers();
    render();
    setupSeedUpload();
}

async function fetchTickers() {
    try {
        const response = await fetch('/api/tickers');
        const data = await response.json();
        const tickerList = data.tickers || [];
        
        for (const t of tickerList) {
            if (!tickers[t]) {
                await addTickerToState(t);
            }
        }
    } catch (e) {
        console.error("Failed to fetch tickers:", e);
    }
}

async function addTickerToState(tickerName) {
    tickerName = tickerName.toLowerCase();
    try {
        const res = await fetch(`/api/state/${tickerName}`);
        if (res.ok) {
            const state = await res.json();
            tickers[tickerName] = {
                ticker: tickerName.toUpperCase(),
                date: state.date || 'Seeded',
                close: state.close || 0,
                open: state.open || 0,
                high: state.high || 0,
                low: state.low || 0,
                hl: state.hl || 0,
                avg: state.avg || 0,
                ema5: state.ema5 || 0,
                ema20: state.ema20 || 0,
                rsi: state.rsi || 0,
                support: state.support || 0,
                bullish: state.bullish || 0,
                warning: null
            };
        } else {
            tickers[tickerName] = { ticker: tickerName.toUpperCase(), date: 'Pending...' };
        }
        subscribeTicker(tickerName);
    } catch (e) {
        console.error(`Error adding ticker ${tickerName}:`, e);
    }
}

function subscribeTicker(tickerName) {
    if (eventSources[tickerName]) return;

    const source = new EventSource(`/api/stream/${tickerName}`);
    source.onmessage = (event) => {
        const data = JSON.parse(event.data);
        updateTickerUI(data);
    };
    source.onerror = (err) => {
        source.close();
        delete eventSources[tickerName];
        setTimeout(() => subscribeTicker(tickerName), 5000);
    };
    eventSources[tickerName] = source;
}

function updateTickerUI(snapshot) {
    const tickerName = snapshot.ticker.toLowerCase();
    const old = tickers[tickerName] || {};
    tickers[tickerName] = snapshot;

    const row = document.getElementById(`row-${tickerName}`);
    if (!row) {
        render();
        return;
    }

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
                setTimeout(() => {
                    cell.classList.add('fade-out');
                    cell.classList.remove(cls);
                }, 50);
            }
        }
    });

    document.getElementById(`d-${tickerName}`).innerText = snapshot.date || '-';
    
    // Update chart if this is the active ticker
    if (activeId === tickerName && candleSeries) {
        candleSeries.update({
            time: snapshot.date,
            open: snapshot.open,
            high: snapshot.high,
            low: snapshot.low,
            close: snapshot.close
        });
    }
}

function render() {
    const tbody = document.getElementById('ticker-body');
    if (!tbody) return;
    tbody.innerHTML = '';

    Object.keys(tickers).sort().forEach(tickerName => {
        const t = tickers[tickerName];
        const tr = document.createElement('tr');
        tr.id = `row-${tickerName}`;
        if (activeId === tickerName) tr.className = 'selected';
        
        tr.onclick = (e) => {
            e.stopPropagation();
            toggleSelection(tickerName);
        };

        tr.innerHTML = `
            <td class="ticker-label">${t.ticker}</td>
            <td id="d-${tickerName}" style="color:var(--text-muted)">${t.date || '-'}</td>
            <td id="c-${tickerName}" class="fade-out">${t.close?.toFixed(2) || '-'}</td>
            <td id="o-${tickerName}">${t.open?.toFixed(2) || '-'}</td>
            <td id="h-${tickerName}">${t.high?.toFixed(2) || '-'}</td>
            <td id="l-${tickerName}">${t.low?.toFixed(2) || '-'}</td>
            <td id="hl-${tickerName}">${t.hl?.toFixed(2) || '-'}</td>
            <td id="avg-${tickerName}">${t.avg?.toFixed(2) || '-'}</td>
            <td id="e5-${tickerName}" class="fade-out">${t.ema5?.toFixed(2) || '-'}</td>
            <td id="e20-${tickerName}">${t.ema20?.toFixed(2) || '-'}</td>
            <td id="r-${tickerName}" class="fade-out">${t.rsi?.toFixed(0) || '-'}</td>
            <td id="s-${tickerName}" class="fade-out">${t.support?.toFixed(2) || '-'}</td>
            <td id="b-${tickerName}" class="fade-out" style="color:var(--accent)">${t.bullish?.toFixed(2) || '-'}</td>
            <td id="a-${tickerName}">${t.warning ? '⚠️' : '-'}</td>
        `;
        tbody.appendChild(tr);
    });
    updateButtons();
}

function toggleSelection(id) {
    destroyChart();
    activeId = (activeId === id) ? null : id;
    render();
}

function handleGlobalClick(e) {
    const table = document.getElementById('ticker-table');
    if (!table) return;
    if (!table.contains(e.target) && !e.target.closest('.toolbar')) {
        activeId = null;
        destroyChart();
        render();
    }
}

function updateButtons() {
    const isSelected = activeId !== null;
    const btnHist = document.getElementById('btn-hist');
    const btnDel = document.getElementById('btn-del');
    if (btnHist) btnHist.disabled = !isSelected;
    if (btnDel) btnDel.disabled = !isSelected;
}

async function triggerChart() {
    if (!activeId) return;

    // Always destroy previous chart instance before creating a new one.
    destroyChart();

    const container = document.getElementById('chart-area');
    container.innerHTML = '';

    chart = LightweightCharts.createChart(container, {
        layout: { backgroundColor: '#080808', textColor: '#d1d1d1' },
        grid: { vertLines: { color: '#111' }, horzLines: { color: '#111' } },
        timeScale: { borderColor: '#222' },
    });

    candleSeries = chart.addCandlestickSeries({
        upColor: '#00ff7f', downColor: '#ff453a', borderVisible: false,
        wickUpColor: '#00ff7f', wickDownColor: '#ff453a',
    });

    try {
        const res = await fetch(`/api/history/${activeId}`);
        const data = await res.json();
        const formatted = data.history
            .map(b => ({ ...b, time: formatTimeForChart(b.time) }))
            .filter(b => b.time !== null)
            .sort((a, b) => a.time.localeCompare(b.time));

        candleSeries.setData(formatted);
        chart.timeScale().fitContent();
    } catch (e) {
        console.error("Failed to load chart data:", e);
        destroyChart();
        container.innerText = '[ HISTORY LOAD FAILED ]';
    }
}

function formatTimeForChart(t) {
    if (!t) return null;
    // Already YYYY-MM-DD
    if (/^\d{4}-\d{2}-\d{2}$/.test(t)) return t;
    // DD-Mon-YYYY (e.g. 20-Dec-2024)
    const months = { Jan:'01', Feb:'02', Mar:'03', Apr:'04', May:'05', Jun:'06',
                     Jul:'07', Aug:'08', Sep:'09', Oct:'10', Nov:'11', Dec:'12' };
    const parts = t.split('-');
    if (parts.length === 3 && months[parts[1]]) {
        return `${parts[2]}-${months[parts[1]]}-${parts[0].padStart(2, '0')}`;
    }
    return null;
}

function destroyChart() {
    if (chart) {
        chart.remove();
        chart = null;
        candleSeries = null;
    }
    document.getElementById('chart-area').innerText = '[ CHART_ENGINE_IDLE ]';
}

function triggerSeedUpload() {
    // Resolve target ticker: use selected row, or prompt for a new name.
    let target = activeId;
    if (!target) {
        const name = prompt("Ticker name to seed (e.g. NIFTY):");
        if (!name || !name.trim()) return;
        target = name.trim().toLowerCase();
    }
    // Store target on the input so the onchange handler can read it.
    const input = document.getElementById('fileInput');
    input.dataset.seedTarget = target;
    // Reset value so re-selecting the same file triggers onchange again.
    input.value = '';
    input.click();
}

function setupSeedUpload() {
    const input = document.getElementById('fileInput');
    input.onchange = async (e) => {
        const target = input.dataset.seedTarget;
        const file = e.target.files[0];
        if (!file || !target) return;

        const formData = new FormData();
        formData.append('file', file);

        try {
            const res = await fetch(`/api/seed/${target}`, {
                method: 'POST',
                body: formData
            });
            const result = await res.json();

            if (res.ok) {
                alert(`SUCCESS: Seeded ${result.bars_loaded} bars for ${target.toUpperCase()}\nDetection: ${result.column_detection}`);
                await addTickerToState(target);
                render();
            } else {
                alert(`UPLOAD FAILED: ${result.detail || 'Unknown error'}`);
            }
        } catch (err) {
            console.error("Seed upload failed:", err);
            alert(`NETWORK ERROR: Could not connect to server.`);
        }
        // Reset so the same file can be re-uploaded later.
        input.value = '';
    };
}

async function removeTicker() {
    if (!activeId) return;
    const tickerToRemove = activeId;
    if (!confirm(`Delete ${tickerToRemove.toUpperCase()}?`)) return;

    try {
        const res = await fetch(`/api/tickers/${tickerToRemove}`, { method: 'DELETE' });
        if (res.ok) {
            if (eventSources[tickerToRemove]) {
                eventSources[tickerToRemove].close();
                delete eventSources[tickerToRemove];
            }
            delete tickers[tickerToRemove];
            activeId = null;
            destroyChart();
            render();
        }
    } catch (e) {
        console.error(`Failed to delete ticker ${tickerToRemove}:`, e);
    }
}

function addTicker() {
    const n = prompt("Ticker Name:");
    if(n) {
        addTickerToState(n.toLowerCase());
    }
}

document.addEventListener('DOMContentLoaded', init);
