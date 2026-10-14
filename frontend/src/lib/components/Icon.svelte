<script lang="ts">
	// Based on https://github.com/AlexxNB/mdi-svelte/blob/master/src/Index.svelte
	// MIT License

	// Copyright (c) 2019 Alexey Schebelev

	// Permission is hereby granted, free of charge, to any person obtaining a copy
	// of this software and associated documentation files (the "Software"), to deal
	// in the Software without restriction, including without limitation the rights
	// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	// copies of the Software, and to permit persons to whom the Software is
	// furnished to do so, subject to the following conditions:

	// The above copyright notice and this permission notice shall be included in all
	// copies or substantial portions of the Software.

	// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	// SOFTWARE.

	interface Props {
		path: string;
		size?: number | string;
		color?: string;
		flip?: string | boolean;
		rotate?: number;
		spin?: boolean;
		title?: string;
	}

	const {
		path,
		size = 1,
		color = 'currentColor',
		flip = false,
		rotate = 0,
		spin = false,
		title = ''
	}: Props = $props();

	let inverse = $derived(typeof spin !== 'boolean' && spin < 0 ? true : false);
	let spintime = $derived.by(() => {
		if (typeof spin === 'boolean') {
			if (spin) {
				return 2;
			} else {
				return false;
			}
		}
		return Math.abs(spin);
	});
	let spinfunc = $derived(inverse ? 'spin-inverse' : 'spin');
	let style = $derived(getStyles(size, color, flip, rotate));

	function getStyles(size: number | string, color: string, flip: boolean | string, rotate: number) {
		const transform = [];
		const styles = [];

		if (size !== null) {
			const width = typeof size === 'string' ? size : `${size * 1.5}rem`;
			styles.push(['width', width]);
			styles.push(['height', width]);
		}

		styles.push(['fill', color]);

		if (flip === true || flip === 'h') {
			transform.push('scaleX(-1)');
		}

		if (flip === true || flip === 'v') {
			transform.push('scaleY(-1)');
		}

		if (rotate != 0) {
			transform.push(`rotate(${rotate}deg)`);
		}

		if (transform.length > 0) {
			styles.push(['transform', transform.join(' ')]);
			styles.push(['transform-origin', 'center']);
		}

		return styles.reduce((cur, item) => {
			return `${cur} ${item[0]}:${item[1]};`;
		}, '');
	}
</script>

<svg viewBox="0 0 24 24" {style}>
	{#if title}<title>{title}</title>{/if}
	{#if spin !== false}
		{#if inverse}
			<style>
				@keyframes spin-inverse {
					to {
						transform: rotate(-360deg);
					}
				}
			</style>
		{:else}
			<style>
				@keyframes spin {
					to {
						transform: rotate(360deg);
					}
				}
			</style>
		{/if}
		<g style={`animation: ${spinfunc} linear ${spintime}s infinite; transform-origin: center`}>
			<path d={path}></path>
		</g>
	{:else}
		<path d={path}></path>
	{/if}
</svg>

<style>
	svg {
		vertical-align: middle;
	}
</style>
