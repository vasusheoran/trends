"""
Global ticker state registry and SSE pub/sub.
"""

import asyncio
from typing import Dict, Set
from app.engine.state import TickerState
from app.models import TickerSnapshot

# ticker -> TickerState
_states: Dict[str, TickerState] = {}

# ticker -> set of subscriber queues
_subscribers: Dict[str, Set[asyncio.Queue]] = {}


def get_state(ticker: str) -> TickerState:
    if ticker not in _states:
        _states[ticker] = TickerState(ticker=ticker)
    return _states[ticker]


async def publish(ticker: str, snapshot: TickerSnapshot) -> None:
    for queue in _subscribers.get(ticker, set()):
        await queue.put(snapshot)


def subscribe(ticker: str, queue: asyncio.Queue) -> None:
    _subscribers.setdefault(ticker, set()).add(queue)


def unsubscribe(ticker: str, queue: asyncio.Queue) -> None:
    _subscribers.get(ticker, set()).discard(queue)
