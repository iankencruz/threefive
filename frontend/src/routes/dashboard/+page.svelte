<script lang="ts">
  import type { PageProps } from "./$types";

  // SvelteKit automatically passes the result of load() into the 'data' prop
  let { data }: PageProps = $props();
</script>

<div class="p-6">
  <h1 class="text-2xl font-bold">Admin Dashboard</h1>

  {#if data.dashboard}
    <pre class="bg-gray-100 p-4 rounded mt-4">
      {JSON.stringify(data.dashboard, null, 4)}
    </pre>
  {:else}
    <p class="text-red-500">{"No data found."}</p>
  {/if}
</div>
