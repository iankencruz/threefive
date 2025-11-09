<script lang="ts">
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
		media?: Media | null; // ✨ NEW: Pre-loaded media from server
	}

	let { data, media = null }: Props = $props();

	// --- utils ---
	const isVideo = (mime: string) => mime?.startsWith("video/");
	const isImage = (mime: string) => mime?.startsWith("image/");
	const getImageUrl = (m: Media) => m.large_url || m.url;
	const getVideoUrl = (m: Media) => m.url;
	const getVideoPoster = (m: Media) => m.thumbnail_url || "";
</script>

<!-- ✨ No more loading state! Media is pre-loaded -->
{#if media}
	<section class="relative overflow-hidden">
		<Navbar />

		<div class="relative min-h-screen flex items-center">
			<div class="absolute inset-0 bg-black">
				{#if isVideo(media.mime_type)}
					<video
						src={getVideoUrl(media)}
						poster={getVideoPoster(media)}
						autoplay
						muted
						loop
						playsinline
						class="w-full h-full object-cover mask-radial-at-right mask-radial-from-100% "
					>
						<track kind="captions" />
					</video>
				{:else if isImage(media.mime_type)}
					<img
						src={getImageUrl(media)}
						alt={data.title}
						class="w-full h-full object-cover mask-radial-at-right mask-radial-from-100% "
					/>
				{/if}

				<div class="absolute inset-0 bg-linear-to-b from-black/85 via-transparent to-black/40"></div>
				<div class="absolute inset-0 bg-linear-to-r from-black/70 via-black/5 to-transparent"></div>
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
