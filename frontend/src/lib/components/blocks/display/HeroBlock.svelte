<!-- frontend/src/lib/components/blocks/display/HeroBlock.svelte -->
<script lang="ts">
interface HeroBlockData {
	title: string;
	subtitle?: string;
	cta_text?: string;
	cta_url?: string;
	image_id?: string;
}

interface Props {
	data: HeroBlockData;
}

let { data }: Props = $props();
</script>

<section class="relative bg-gradient-to-br from-blue-600 to-blue-800 text-white py-20 md:py-32">
	<div class="container mx-auto px-4 max-w-6xl">
		<div class="text-center">
			<h1 class="text-4xl md:text-5xl lg:text-6xl font-bold mb-6 leading-tight">
				{data.title}
			</h1>
			
			{#if data.subtitle}
				<p class="text-xl md:text-2xl text-blue-100 mb-8 max-w-3xl mx-auto">
					{data.subtitle}
				</p>
			{/if}
			
			{#if data.cta_text && data.cta_url}
				<a
					href={data.cta_url}
					class="inline-block bg-white text-blue-600 px-8 py-4 rounded-lg font-semibold text-lg hover:bg-blue-50 transition-colors shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-transform"
				>
					{data.cta_text}
				</a>
			{/if}
		</div>
	</div>
	
	<!-- Decorative gradient overlay -->
	<div class="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent pointer-events-none"></div>
</section>
