<script lang="ts">
	import { mdiGithub, mdiLibraryShelves } from '@mdi/js';
	import NavBar from '#lib/components/NavBar.svelte';
	import { version } from '../../../package.json';
	import Icon from '#lib/components/Icon.svelte';

	let { data } = $props();
</script>

<svelte:head>
	<title>Albums2: About</title>
</svelte:head>

<NavBar user={data.user} up={{ description: 'Collections', icon: mdiLibraryShelves, url: '/' }} />

<div class="overflow-auto p-4 pt-8">
	<article class="mx-auto prose w-full md:w-200">
		<h1>About Albums2</h1>
		<p>Version {version}</p>
		<p>
			<a href="https://github.com/wutipong/albums2" target="_blank" rel="noreferrer" class='btn btn-link'>
				<Icon path={mdiGithub}></Icon>GitHub
			</a>
		</p>

		<h2>MIT License</h2>
		<p>Copyright &copy; 2026 Wutipong Wongsakuldej</p>
		<p>
			Permission is hereby granted, free of charge, to any person obtaining a copy of this software
			and associated documentation files (the &quot;Software&quot;), to deal in the Software without
			restriction, including without limitation the rights to use, copy, modify, merge, publish,
			distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the
			Software is furnished to do so, subject to the following conditions:
		</p>
		<p>
			The above copyright notice and this permission notice shall be included in all copies or
			substantial portions of the Software.
		</p>
		<p>
			THE SOFTWARE IS PROVIDED &quot;AS IS&quot;, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
			INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
			PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR
			ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
			ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
			SOFTWARE.
		</p>
	</article>
</div>
