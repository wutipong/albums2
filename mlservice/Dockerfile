FROM ghcr.io/astral-sh/uv:python3.11-bookworm-slim

WORKDIR /app

# Enable bytecode compilation and set link mode for faster installs
ENV UV_COMPILE_BYTECODE=1
ENV UV_LINK_MODE=copy
ENV HF_HOME=/app/hf_cache

# Install dependencies using the lockfile
RUN --mount=type=cache,target=/root/.cache/uv \
    --mount=type=bind,source=uv.lock,target=uv.lock \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    uv sync --frozen --no-install-project

# Copy the rest of the application source
COPY . .

# Sync the full project
RUN --mount=type=cache,target=/root/.cache/uv \
    uv sync --frozen

EXPOSE 8173

RUN uv add huggingface_hub[cli]
RUN uv run huggingface-cli download google/siglip2-base-patch16-256 --local-dir /app/model

CMD ["uv", "run", "server.py"]
