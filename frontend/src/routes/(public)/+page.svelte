<script lang="ts">
	let { data } = $props();
	console.log('Home data:', data);
</script>

<h1 class="text-2xl font-semibold">Home Page</h1>

<p>{data.Hero}</p>
<!-- <div class="w-full max-w-4xl"> -->
<!-- 	<div class="grid grid-cols-2 gap-2"> -->
<!-- 		{#each data.projects.items as item} -->
<!-- 			<img src={item.url} class="h-auto w-full" /> -->
<!-- 		{/each} -->
<!-- 	</div> -->
<!-- </div> -->
