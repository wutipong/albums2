# Architecture and service setup

## Components

- **Frontend:** the [SvelteKit app](./frontend/) serves collection and album pages,
  media views, administrative pages, and API routes. It uses Better Auth with an OIDC
  provider. MongoDB stores application and authentication data; Redis holds processing
  queues; imgproxy serves signed image URLs for objects in S3-compatible storage.
- **Worker:** the [Go worker](./worker/) consumes Redis tasks and processes collection
  and album synchronization, covers, and media embeddings. It accesses MongoDB and
  S3-compatible storage and calls the ML service over gRPC.
- **ML service:** the [Python service](./mlservice/) exposes image and text encoding
  through the [shared gRPC protocol](./grpc/mlservice.proto). It uses the
  `google/siglip2-base-patch16-256` model.
- **Search:** the frontend sends text to the ML service and queries MongoDB Vector
  Search using the resulting embedding. The media schema defines a 768-dimensional
  vector index; MongoDB must support Search and Vector Search.

The gRPC service currently uses an insecure connection; keep it on a trusted network
or provide transport protection externally.

## Configuration and startup

Run development commands inside the provided
[dev container](./.devcontainer/devcontainer.json).
Configure addresses for the network where each process runs. For example, container
service names and `localhost` are not interchangeable. Do not commit secrets.

### Supporting services

The [dev-container Compose file](./.devcontainer/docker-compose.yml) runs MongoDB Atlas
Local, Redis, Versity Gateway, imgproxy, Dex, and their administrative tools. It expects
`.devcontainer/.env` and matching storage and imgproxy settings. Create the S3 bucket and
configure the OIDC client and application URLs for your environment. The Compose file does
not start the worker or ML service.

### Frontend

Configure the frontend environment as described in the
[frontend README](./frontend/README.md), then start it with:

```sh
cd frontend
bun install
bun run --bun dev
```

### ML service

Python 3.11 or newer and `uv` are required. The model and processor must be available
from Hugging Face or its local cache; the first startup can download model files.

```sh
cd mlservice
uv sync --frozen
uv run --frozen server.py
```

`CLIP_ADDRESS` controls the gRPC listen address and defaults to `0.0.0.0:8173`.
Set `MLSERVICE_ADDRESS` in the frontend and worker to an address reachable from those
processes, such as `localhost:8173` when they share a host.

### Worker

The worker requires the Go version specified in [`worker/go.mod`](./worker/go.mod),
CGO, libvips, and FFmpeg. Configure its environment before starting it:

| Variables                                                                                           | Purpose                                                            |
| --------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------ |
| `DB_CONNECTION`, `REDIS_URL`                                                                        | MongoDB URI, including the database name, and Redis URL.           |
| `AWS_ENDPOINT_URL`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_DEFAULT_REGION`, `S3_BUCKET` | S3-compatible storage settings.                                    |
| `MLSERVICE_ADDRESS`                                                                                 | Reachable address of the ML gRPC service.                          |
| `WORKER_ID`                                                                                         | Unique integer for each worker sharing the queue; defaults to `0`. |
| `APP_ENV`, `LOG_LEVEL`                                                                              | Optional environment and logging settings.                         |

```sh
cd worker
go run .
```

The worker also accepts command-line flags; `go run . -help` lists them.

## Protocol generation

The gRPC schema is maintained in [`grpc/mlservice.proto`](./grpc/mlservice.proto).
The frontend's `bun run gen:protoc` and worker's `go generate` directives expect the
repository at `/workspaces` and require `protoc` and their respective language plugins.
The Python client can be regenerated from `mlservice/` using `./codegen.sh`.
