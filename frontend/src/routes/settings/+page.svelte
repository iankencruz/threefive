<script lang="ts">
	import type { PageProps } from './$types';

	// SvelteKit automatically passes the result of load() into the 'data' prop
	let { data }: PageProps = $props();
</script>

<div class="p-6">
	<h1 class="text-2xl font-bold">Settings</h1>

	{#if data.config}
		<pre class="mt-4 rounded bg-gray-100 p-4">
      {JSON.stringify(data.config, null, 4)}
    </pre>
	{:else}
		<p class="text-red-500">{'No data found.'}</p>
	{/if}
</div>
