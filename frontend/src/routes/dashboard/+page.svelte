<script lang="ts">
  import type { PageProps } from "./$types";

  // Destructure the props
  let { data }: { data: PageProps } = $props();
</script>

<div class="p-6">
  <h1 class="text-2xl font-bold">Admin Dashboard</h1>

  {#if data}
    <pre class="bg-gray-100 p-4 rounded mt-4">
      {JSON.stringify(data, null, 2)}
    </pre>
  {:else}
    <p>No data found or unauthorized.</p>
  {/if}
</div>
