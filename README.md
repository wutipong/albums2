# albums2

albums2 is an application for browsing collections of images and videos stored in
S3-compatible object storage. It combines a SvelteKit frontend, a Go processing worker,
and a Python service that generates image and text embeddings for semantic search.

The application is under development. The frontend's album and asset administration
pages are currently placeholders.

## Repository

| Directory                            | Contents                                                                 |
| ------------------------------------ | ------------------------------------------------------------------------ |
| [`frontend/`](./frontend/)           | SvelteKit application, API routes, and PWA assets.                       |
| [`worker/`](./worker/)               | Go worker for collection and album synchronization and media processing. |
| [`mlservice/`](./mlservice/)         | Python gRPC service for image and text embeddings.                       |
| [`grpc/`](./grpc/)                   | Shared protocol definitions.                                             |
| [`.devcontainer/`](./.devcontainer/) | Development container and supporting services.                           |

See [Architecture](./ARCHITECTURE.md) for the service interactions and configuration.
The [frontend README](./frontend/README.md) covers its environment variables, commands,
container build, and PWA patches.

## Getting started

**Development should be done inside the provided
[dev container](./.devcontainer/devcontainer.json).** It supplies Bun, Go, Python,
`uv`, worker dependencies, and supporting services including MongoDB Atlas Local,
Redis, S3-compatible storage, imgproxy, and Dex.

1. Configure `.devcontainer/.env` for the Compose services and application. The file is
   referenced by [docker-compose.yml](./.devcontainer/docker-compose.yml); there is no
   checked-in example. Do not commit credentials.
2. In VS Code, open the repository and select **Dev Containers: Reopen in Container**.
   Run development commands from the container's terminal. Configure the frontend,
   worker, and ML service as described in [Architecture](./ARCHITECTURE.md).
3. Start the ML service and worker separately. They are not included in the dev-container
   Compose services.
4. Install and start the frontend:

   ```sh
   cd frontend
   bun install
   bun run --bun dev
   ```

Vite listens on `127.0.0.1:5173` by default. MongoDB and Redis must be reachable when
the frontend starts. Sign in through the configured OIDC provider; the first user created
is assigned the administrator role.

## Development and builds

Frontend checks, build instructions, and container usage are documented in
[`frontend/README.md`](./frontend/README.md). Each service has its own Dockerfile;
the dev-container Compose file is for supporting services, not a complete production stack.

The generated gRPC clients are based on [`grpc/mlservice.proto`](./grpc/mlservice.proto).
See [Architecture](./ARCHITECTURE.md) for regeneration commands.
