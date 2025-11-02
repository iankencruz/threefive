
<script lang="ts">
	import { PUBLIC_API_URL } from "$env/static/public";
	import type { Media } from "$api/media";
	import Navbar from "$components/ui/Navbar.svelte";

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

	// --- async fetcher ---
	async function loadMedia(
		id: string | null | undefined,
	): Promise<Media | null> {
		if (!id) return null;

		const res = await fetch(`${PUBLIC_API_URL}/api/v1/media/${id}`, {
			credentials: "include",
		});
		if (!res.ok) {
			throw new Error(`Failed to load media (status ${res.status})`);
		}
		// Optional demo delay:
		// await new Promise((r) => setTimeout(r, 800));
		return (await res.json()) as Media;
	}

	// 🔁 Recomputes whenever data.image_id changes
	let mediaPromise = $derived(loadMedia(data.image_id));

	// --- utils ---
	const isVideo = (mime: string) => mime?.startsWith("video/");
	const isImage = (mime: string) => mime?.startsWith("image/");
	const getImageUrl = (m: Media) => m.large_url || m.url;
	const getVideoUrl = (m: Media) => m.url;
	const getVideoPoster = (m: Media) => m.thumbnail_url || "";
</script>

{#await mediaPromise}
	<!-- pending -->
	<div class="min-h-[500px] flex items-center justify-center bg-gray-100">
		<div class="flex flex-col items-center gap-3">
			<div class="w-12 h-12 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin"></div>
			<p class="text-sm text-gray-600">Loading media...</p>
		</div>
	</div>

{:then media}
	{#if media}
		<section class="relative overflow-hidden">
			<Navbar />

			<div class="relative min-h-screen flex items-center">
				<div class="absolute inset-0">
					{#if isVideo(media.mime_type)}
						<video
							src={getVideoUrl(media)}
							poster={getVideoPoster(media)}
							autoplay
							muted
							loop
							playsinline
							class="w-full h-full object-cover"
						>
							<track kind="captions" />
						</video>
					{:else if isImage(media.mime_type)}
						<img
							src={getImageUrl(media)}
							alt={data.title}
							class="w-full h-full object-cover"
						/>
					{/if}

					<div class="absolute inset-0 bg-gradient-to-r from-black/70 via-black/50 to-transparent"></div>
				</div>

				<div class="@container absolute z-10 container mx-auto px-4 max-w-6xl left-12 bottom-8">
					<div class="w-full text-white">
						<h1 class="@max-lg:text-2xl @lg:@max-2xl:text-5xl @2xl:text-7xl @5xl:text-9xl font-bold mb-6 leading-tight">
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
		</section>
	{:else}
		<section class="relative min-h-screen bg-gray-900 flex items-center justify-center">
			<p class="text-white text-xl">{data.title} (No background media specified)</p>
		</section>
	{/if}

{:catch error}
	<!-- error -->
	<div class="min-h-[500px] flex items-center justify-center bg-gray-100">
		<div class="text-center">
			<svg class="w-16 h-16 mx-auto text-red-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
			</svg>
			<p class="text-gray-700 text-lg">Failed to load media</p>
			{#if error?.message}
				<p class="text-gray-500 text-sm mt-1">{error.message}</p>
			{/if}
		</div>
	</div>
{/await}

