# sv

Everything you need to build a Svelte project, powered by [`sv`](https://github.com/sveltejs/cli).

## Creating a project

If you're seeing this, you've probably already done this step. Congrats!

```sh
# create a new project
npx sv create my-app
```

To recreate this project with the same configuration:

```sh
# recreate this project
bun x sv@0.17.1 create --template minimal --types ts --add prettier eslint tailwindcss="plugins:typography,forms" better-auth="demo:none" mdsvex drizzle="database:sqlite+sqlite:libsql" --install bun my-app
```

## Developing

Install dependencies with Bun, then start a development server:

```sh
bun install
bun run --bun dev

# or start the server and open the app in a new browser tab
bun run --bun dev --open
```

## Building

To create a production version of your app:

```sh
bun run --bun build
```

You can preview the production build with `bun run --bun preview`.

The project uses the Bun adapter for self-hosted deployments. Start the generated server with:

```sh
bun run start
```

The production server listens on port `3000` by default. Set `HOST` and `PORT` to change its bind address.

Bun applies the patches in `patches/` during installation. They work around missing
`node:module` APIs in the PWA dependencies; review them when upgrading those packages.
