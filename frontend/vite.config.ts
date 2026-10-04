import { mdsvex } from 'mdsvex';
import tailwindcss from '@tailwindcss/vite';
import adapter from '@sveltejs/adapter-auto';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import openapiPlugin from 'sveltekit-openapi-generator';

export default defineConfig({
	plugins: [
		tailwindcss(),
		openapiPlugin({
			// OpenAPI info section
			info: {
				title: 'Albums2 API',
				version: '1.0.0',
				description: 'Albums2 API documentation',
			},
			// OpenAPI servers configuration
			servers: [
				{ url: 'http://localhost:5173', description: 'Development' }
			],
			// Path to shared schema definitions
			baseSchemasPath: 'src/lib/schemas.js',
			// Additional YAML files to include
			yamlFiles: ['src/lib/extra-specs.yaml'],
			// Path prefix for all routes
			prependPath: '/api',
			// Glob patterns to include
			include: ['src/routes/**/{+server,+page.server}.{js,ts}'],
			// Glob patterns to exclude
			exclude: ['**/node_modules/**', '**/.svelte-kit/**'],
			// Whether to fail on JSDoc parsing errors
			failOnErrors: false,
			// Output path for the spec file during build
			outputPath: 'static/openapi.json',
			// Debounce delay in milliseconds for file watching
			debounceMs: 200
		}),
		sveltekit({
			compilerOptions: {
				// Force runes mode for the project, except for libraries. Can be removed in svelte 6.
				runes: ({ filename }) =>
					filename.split(/[/\\]/).includes('node_modules') ? undefined : true
			},

			// adapter-auto only supports some environments, see https://svelte.dev/docs/kit/adapter-auto for a list.
			// If your environment is not supported, or you settled on a specific environment, switch out the adapter.
			// See https://svelte.dev/docs/kit/adapters for more information about adapters.
			adapter: adapter(),
			preprocess: [mdsvex({ extensions: ['.svx', '.md'] })],
			extensions: ['.svelte', '.svx', '.md'],
			typescript: {
				config: (config) => {
					// config.include.push('../drizzle.config.ts');
				}
			}
		})
	]
});
