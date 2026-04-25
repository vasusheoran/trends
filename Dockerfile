FROM python:3.12-slim

# Install uv
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

WORKDIR /app

# Install dependencies first (layer-cached until pyproject/lockfile changes)
COPY trends-py/pyproject.toml trends-py/uv.lock ./
RUN uv sync --frozen --no-dev

# Copy application code
COPY trends-py/app ./app

# Copy Excel seed data
# COPY data/Final-bullish-ce.xlsx /data/Final-bullish-ce.xlsx

EXPOSE 5001

CMD ["uv", "run", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "5001"]
