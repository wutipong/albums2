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

Once you've created a project and installed dependencies with `npm install` (or `pnpm install` or `yarn`), start a development server:

```sh
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```sh
npm run build
```

You can preview the production build with `npm run preview`.

This project uses the Node adapter for self-hosted deployments. Build and start the generated server with:

```sh
bun run build
bun run start
```

The production server listens on port `3000` by default. Set `HOST` and `PORT` to change its bind address.
