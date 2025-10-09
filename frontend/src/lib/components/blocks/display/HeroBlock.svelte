<!-- frontend/src/lib/components/blocks/display/HeroBlock.svelte -->
<script lang="ts">
import { onMount } from "svelte";
import { getMediaUrl, type Media } from "$api/media";
import { PUBLIC_API_URL } from "$env/static/public";

interface HeroBlockData {
	title: string;
	subtitle?: string | null;
	image_id?: string | null;
	cta_text?: string | null;
	cta_url?: string | null;
}

interface Props {
	data: HeroBlockData;
}

let { data }: Props = $props();

let media = $state<Media | null>(null);
let loading = $state(false);
let error = $state(false);

onMount(async () => {
	if (data.image_id) {
		await loadMedia(data.image_id);
	}
});

// Watch for changes to image_id
$effect(() => {
	if (data.image_id && (!media || media.id !== data.image_id)) {
		loadMedia(data.image_id);
	} else if (!data.image_id) {
		media = null;
	}
});

async function loadMedia(imageId: string) {
	loading = true;
	error = false;

	try {
		const response = await fetch(`${PUBLIC_API_URL}/api/v1/media/${imageId}`, {
			credentials: "include",
		});

		if (!response.ok) {
			throw new Error("Failed to load media");
		}

		media = await response.json();
	} catch (err) {
		console.error("Failed to load hero media:", err);
		error = true;
		media = null;
	} finally {
		loading = false;
	}
}
</script>

<section class="relative overflow-hidden">
	{#if media && !loading && !error}
		<!-- Hero with Background Image -->
		<div class="relative min-h-[500px] md:min-h-[600px] flex items-center">
			<!-- Background Image -->
			<div class="absolute inset-0">
				<img
					src={getMediaUrl(media)}
					alt={data.title}
					class="w-full h-full object-cover"
				/>
				<!-- Gradient Overlay -->
				<div class="absolute inset-0 bg-gradient-to-r from-black/70 via-black/50 to-transparent"></div>
			</div>
			
			<!-- Content -->
			<div class="relative z-10 container mx-auto px-4 py-20 md:py-32 max-w-6xl">
				<div class="max-w-2xl text-white">
					<h1 class="text-4xl md:text-5xl lg:text-6xl font-bold mb-6 leading-tight">
						{data.title}
					</h1>
					
					{#if data.subtitle}
						<p class="text-xl md:text-2xl mb-8 text-gray-100 leading-relaxed">
							{data.subtitle}
						</p>
					{/if}
					
					{#if data.cta_text && data.cta_url}
						<a
							href={data.cta_url}
							class="inline-flex items-center gap-2 bg-white text-gray-900 px-8 py-4 rounded-lg font-semibold text-lg hover:bg-gray-100 transition-all shadow-lg hover:shadow-xl hover:scale-105"
						>
							{data.cta_text}
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
							</svg>
						</a>
					{/if}
				</div>
			</div>
		</div>
	{:else}
		<!-- Hero without Image (Gradient Background) -->
		<div class="relative bg-gradient-to-br from-blue-600 via-blue-700 to-blue-800">
			<div class="container mx-auto px-4 py-20 md:py-32 max-w-6xl">
				<div class="text-center text-white">
					<h1 class="text-4xl md:text-5xl lg:text-6xl font-bold mb-6 leading-tight">
						{data.title}
					</h1>
					
					{#if data.subtitle}
						<p class="text-xl md:text-2xl text-blue-100 mb-8 max-w-3xl mx-auto leading-relaxed">
							{data.subtitle}
						</p>
					{/if}
					
					{#if data.cta_text && data.cta_url}
						<a
							href={data.cta_url}
							class="inline-flex items-center gap-2 bg-white text-blue-600 px-8 py-4 rounded-lg font-semibold text-lg hover:bg-blue-50 transition-all shadow-lg hover:shadow-xl hover:scale-105"
						>
							{data.cta_text}
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
							</svg>
						</a>
					{/if}
				</div>
			</div>
			
			<!-- Decorative gradient overlay -->
			<div class="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent pointer-events-none"></div>
		</div>
	{/if}
	
	{#if loading}
		<div class="absolute inset-0 flex items-center justify-center bg-gray-100">
			<div class="flex flex-col items-center gap-3">
				<div class="w-12 h-12 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin"></div>
				<p class="text-sm text-gray-600">Loading image...</p>
			</div>
		</div>
	{/if}
</section>
