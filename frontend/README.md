# albums2 frontend

albums2 is a SvelteKit application for browsing image and video collections and albums.
It uses Bun, MongoDB, Redis, S3-compatible storage, imgproxy, and an OIDC identity provider.
Text search uses embeddings from the repository's gRPC ML service and MongoDB Vector Search.

The frontend includes collection browsing, media viewing, text search, an administration
area for API keys and the processing queue, and collection and media API routes.
The album and asset administration pages are currently placeholders.

## Requirements

- Bun for installing dependencies, running the development server, and building.
- MongoDB with Search and Vector Search support, Redis, S3-compatible storage, imgproxy,
  and an OIDC provider. The frontend connects to MongoDB and Redis at startup.
- The [ML service](../mlservice/) for text search and the [worker](../worker/) for media
  processing.

The repository's [dev container](../.devcontainer/devcontainer.json) provides MongoDB Atlas
Local, Redis, S3-compatible storage, imgproxy, and Dex through
[Docker Compose](../.devcontainer/docker-compose.yml). It does not start the worker or
ML service; run those separately when needed.

## Configuration

Provide the following environment variables to the frontend process. Use values appropriate
to your deployment and keep credentials out of version control.

| Variables                                                                                           | Purpose                                              |
| --------------------------------------------------------------------------------------------------- | ---------------------------------------------------- |
| `DB_CONNECTION`                                                                                     | MongoDB connection URI, including the database name. |
| `REDIS_URL`                                                                                         | Redis connection URL.                                |
| `ORIGIN`, `BETTER_AUTH_SECRET`                                                                      | Public application origin and Better Auth secret.    |
| `OIDC_PROVIDER_ID`, `OIDC_CLIENT_ID`, `OIDC_SECRET`, `OIDC_DISCOVERY_URL`                           | OIDC provider and client configuration.              |
| `AWS_DEFAULT_REGION`, `AWS_ENDPOINT_URL`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `S3_BUCKET` | S3-compatible storage connection and bucket.         |
| `IMGPROXY_URL`, `IMGPROXY_KEY`, `IMGPROXY_SALT`                                                     | Image proxy endpoint and URL-signing credentials.    |
| `MLSERVICE_ADDRESS`                                                                                 | gRPC address of the text-encoding ML service.        |

The first user created through authentication is assigned the administrator role.
Other users require the appropriate role to access `/admin`.

## Development

From the `frontend` directory:

```sh
bun install
bun run --bun dev
```

Vite listens on `127.0.0.1:5173` by default. To open a browser automatically, use
`bun run --bun dev --open`.

```sh
bun run --bun check   # Svelte and TypeScript diagnostics
bun run --bun lint    # Prettier and ESLint
bun run --bun format  # Format the project
```

`bun run gen:protoc` regenerates the frontend gRPC client from
[`grpc/mlservice.proto`](../grpc/mlservice.proto). It expects `protoc` and the repository
at `/workspaces`, as in the dev container.

## Build and deployment

```sh
bun run --bun build
bun run start
```

The Bun adapter writes the server to `build/`. The server listens on port `3000` by default;
set `HOST` and `PORT` as needed. `bun run --bun preview` starts the Vite preview server.

The [Dockerfile](./Dockerfile) builds the frontend in a Bun build stage and copies its
output and dependencies into a separate Bun runtime image:

```sh
docker build -t albums2-frontend .
docker run --env-file .env -p 3000:3000 albums2-frontend
```

Provide a local `.env` file and make the required services reachable from the container.
The Docker build excludes local environment files and dependencies from its build context.

## PWA assets

The PWA integration generates icons from [`static/icon.svg`](./static/icon.svg), including
the Apple touch icon and manifest icons, and generates a service worker during the build.
Bun applies the patches in [`patches/`](./patches/) during installation. They address
PWA dependency compatibility with Bun and production asset generation; review them when
upgrading the patched packages.
