# UI Specification: "trends" Terminal

## 1. Visual Architecture
* **Theme Support:** Must support dynamic switching between Dark (default) and Light modes via CSS variables.
* **State Management:** Local selection state for tickers; real-time DOM updates for SSE streams.
* **Layout:** 
    * Top: Fixed Horizontal Toolbar.
    * Center: High-density data grid (flex-grow).
    * Bottom: Fixed Chart Placeholder (min-height: 220px).
* **Visuals:** Dark mode (pure blacks/deep grays). Monospaced fonts for numerical data. 500ms flash indicators for price changes.
* **Typography:** Use monospaced fonts (e.g., JetBrains Mono) for all numerical cells to prevent layout shifting during updates.

## 2. Interaction Model
* **Selection:** * Single-click on a row selects that ticker.
    * Clicking an already selected row or the "dead space" (grid background) deselects it.
* **Contextual Actions:** "View History" and "Remove Ticker" buttons are disabled until a ticker is selected.
* **Real-time Feedback:** * 500ms flash effect on numeric changes (Green for up, Red for down).
    * Flash should be instantaneous; transition back to base color should be linear (0.5s).

## 3. Technical Integration
* **Discovery:** Consult `http://192.168.29.204:5001/docs` for API routing.
* **Streaming:** Use native `EventSource` for SSE. Do not poll.
* **Performance:** Use `document.getElementById` or targeted refs to update specific `<td>` nodes directly; avoid full table re-renders on every tick.
