<!-- frontend/src/lib/components/blocks/display/HeaderBlock.svelte -->
<script lang="ts">
	interface HeaderBlockData {
		heading: string;
		subheading?: string;
		level: string; // h1, h2, h3, h4, h5, h6
	}

	interface Props {
		data: HeaderBlockData;
	}

	let { data }: Props = $props();

	// Get the appropriate heading styles based on level
	const getHeadingClasses = (level: string) => {
		switch (level) {
			case "h1":
				return "text-4xl md:text-5xl lg:text-6xl font-bold";
			case "h2":
				return "text-3xl md:text-4xl lg:text-5xl font-bold";
			case "h3":
				return "text-2xl md:text-3xl lg:text-4xl font-semibold";
			case "h4":
				return "text-xl md:text-2xl lg:text-3xl font-semibold";
			case "h5":
				return "text-lg md:text-xl lg:text-2xl font-semibold";
			case "h6":
				return "text-base md:text-lg lg:text-xl font-semibold";
			default:
				return "text-2xl md:text-3xl lg:text-4xl font-semibold";
		}
	};

	const headingClasses = getHeadingClasses(data.level);
</script>

<section class="py-12 md:py-16 bg-background">
	<div class="container mx-auto px-4 max-w-7xl">
		<div class="text-center">
			{#if data.level === 'h1'}
				<h1 class="{headingClasses}  mb-4">
					{data.heading}
				</h1>
			{:else if data.level === 'h2'}
				<h2 class="{headingClasses}  mb-4">
					{data.heading}
				</h2>
			{:else if data.level === 'h3'}
				<h3 class="{headingClasses}  mb-4">
					{data.heading}
				</h3>
			{:else if data.level === 'h4'}
				<h4 class="{headingClasses}  mb-4">
					{data.heading}
				</h4>
			{:else if data.level === 'h5'}
				<h5 class="{headingClasses}  mb-4">
					{data.heading}
				</h5>
			{:else}
				<h6 class="{headingClasses}  mb-4">
					{data.heading}
				</h6>
			{/if}

			{#if data.subheading}
				<p class="text-lg md:text-xl max-w-2xl mx-auto">
					{data.subheading}
				</p>
			{/if}
		</div>
	</div>
</section>
