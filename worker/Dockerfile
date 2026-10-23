# ---- Build Stage ----
FROM ubuntu:26.04 AS builder

# Install native dependencies required by govips and ffmpeg-go
RUN apt-get update && apt-get install -y --fix-broken \
    && apt-get install -y --no-install-recommends \
    build-essential \
    golang-go \
    libvips-dev \
    ffmpeg \
    libvpx12 \
    libopus0 \
    ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

WORKDIR /app

# Copy source and build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o worker ./main.go

# ---- Runtime Stage ----
FROM ubuntu:26.04
# Install runtime native libraries (libvips, ffmpeg)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libvips \
    ffmpeg \
    libvpx12 \
    libopus0 \
    ca-certificates 

WORKDIR /app

COPY --from=builder /app/worker .

ENTRYPOINT ["/app/worker"]
