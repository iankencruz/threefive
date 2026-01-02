<script lang="ts">
	const { children } = $props();
</script>

<div class=" min-h-screen bg-background">
	<div class="max-w-8xl mx-auto px-4 sm:px-6 lg:px-8">
		{@render children?.()}
	</div>
</div>
