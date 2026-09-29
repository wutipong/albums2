FROM mcr.microsoft.com/devcontainers/base:dev-ubuntu26.04

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl debian-keyring debian-archive-keyring apt-transport-https gnupg wget ca-certificates \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/apt/keyrings \
    && curl -fsSL https://www.mongodb.org/static/pgp/server-8.0.asc | \
       gpg --dearmor -o /etc/apt/keyrings/mongodb-server-8.0.gpg

RUN echo "deb [ arch=amd64,arm64 signed-by=/etc/apt/keyrings/mongodb-server-8.0.gpg ] https://repo.mongodb.org/apt/ubuntu noble/mongodb-org/8.0 multiverse" \
    > /etc/apt/sources.list.d/mongodb-org-8.0.list

RUN apt-get update && apt-get install -y --no-install-recommends \
    mongodb-mongosh \
    caddy \
    ffmpeg \
    libopus0 \
    libvips-dev \
    libvips-tools \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

USER vscode

RUN mkdir -p /home/vscode/.mc
RUN mkdir -p /home/vscode/bin
RUN mkdir -p /home/vscode/.cache/huggingface