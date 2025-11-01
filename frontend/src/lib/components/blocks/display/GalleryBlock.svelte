<script lang="ts">
	import { onMount } from "svelte";
	import { X } from "lucide-svelte";
	import type { Media } from "$api/media";

	interface GalleryBlockData {
		title?: string;
		media?: Media[];
	}

	interface Props {
		data: GalleryBlockData;
	}

	let { data }: Props = $props();

	function isVideo(mimeType: string): boolean {
		return mimeType.startsWith("video/");
	}

	function isImage(mimeType: string): boolean {
		return mimeType.startsWith("image/");
	}

	// Get the video URL - use url which contains the optimized video
	function getVideoUrl(media: Media): string {
		return media.url;
	}

	// Get poster image for video
	function getVideoPoster(media: Media): string {
		return media.thumbnail_url || "";
	}

	let lightboxOpen = $state(false);
	let lightboxIndex = $state(0);

	function openLightbox(index: number) {
		lightboxIndex = index;
		lightboxOpen = true;
		// Prevent body scroll when lightbox is open
		document.body.style.overflow = "hidden";
	}

	function closeLightbox() {
		lightboxOpen = false;
		document.body.style.overflow = "";
	}

	function nextImage() {
		if (data.media && lightboxIndex < data.media.length - 1) {
			lightboxIndex++;
		}
	}

	function prevImage() {
		if (lightboxIndex > 0) {
			lightboxIndex--;
		}
	}

	// Handle keyboard navigation
	function handleKeydown(e: KeyboardEvent) {
		if (!lightboxOpen) return;

		if (e.key === "Escape") {
			closeLightbox();
		} else if (e.key === "ArrowRight") {
			nextImage();
		} else if (e.key === "ArrowLeft") {
			prevImage();
		}
	}

	onMount(() => {
		window.addEventListener("keydown", handleKeydown);
		return () => {
			window.removeEventListener("keydown", handleKeydown);
			document.body.style.overflow = "";
		};
	});

	function getImageUrl(
		media: Media,
		size: "thumbnail" | "medium" | "large" | "original" = "medium",
	): string {
		switch (size) {
			case "thumbnail":
				return media.thumbnail_url || media.medium_url || media.url;
			case "medium":
				return media.medium_url || media.url;
			case "large":
				return media.large_url || media.url;
			case "original":
				return media.original_url || media.url;
			default:
				return media.url;
		}
	}
</script>

<section class="py-12 md:py-16 bg-background">
	<div class="container mx-auto px-4 max-w-7xl">
		{#if data.title}
			<h2 class="text-3xl md:text-4xl font-bold text-center mb-8 md:mb-12">
				{data.title}
			</h2>
		{/if}

		{#if data.media && data.media.length > 0}
			<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
				{#each data.media as media, index (media.id)}
					<button
						type="button"
						onclick={() => openLightbox(index)}
						class="group relative aspect-square overflow-hidden rounded-lg bg-gray-200 hover:opacity-90 transition-opacity cursor-pointer"
					>
          {#if isVideo(media.mime_type)}
					<!-- Video Background - Optimized MP4 -->
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
							src={getImageUrl(media, 'medium')}
							alt={media.original_filename}
							class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
							loading="lazy"
						/>
          {/if}
						<div class="absolute inset-0  group-hover:opacity-20 transition-all duration-300"></div>
					</button>
				{/each}
			</div>
		{:else}
			<div class="text-center py-12">
				<p class="text-gray-500">No images in this gallery</p>
			</div>
		{/if}
	</div>
</section>

<!-- Lightbox Modal -->
{#if lightboxOpen && data.media && data.media[lightboxIndex]}
	<!-- Lightbox container -->
	<div class="fixed inset-0 z-50" role="dialog" aria-modal="true">
		<!-- Background overlay (clickable to close) -->
		<button
      type="button"
			class="absolute inset-0 bg-black opacity-85"
			onclick={closeLightbox}
			aria-label="Close lightbox"
		></button>

		<!-- Content layer -->
		<div class="relative w-full h-full flex items-center justify-center pointer-events-none">
			<!-- Close button -->
			<button
				type="button"
				onclick={closeLightbox}
				class="absolute top-4 right-4 z-10 w-10 h-10 flex items-center justify-center  bg-opacity-10 hover:opacity-20 rounded-full transition-all pointer-events-auto"
				aria-label="Close lightbox"
			>
				<X class="w-6 h-6 text-white" />
			</button>

			<!-- Previous button -->
			{#if lightboxIndex > 0}
				<button
					type="button"
					onclick={(e) => { e.stopPropagation(); prevImage(); }}
					class="absolute left-4 z-10 w-12 h-12 flex items-center justify-center  bg-opacity-10 hover:opacity-20 rounded-full transition-all pointer-events-auto"
					aria-label="Previous image"
				>
					<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
					</svg>
				</button>
			{/if}

			<!-- Next button -->
			{#if data.media && lightboxIndex < data.media.length - 1}
				<button
					type="button"
					onclick={(e) => { e.stopPropagation(); nextImage(); }}
					class="absolute right-4 z-10 w-12 h-12 flex items-center justify-center  bg-opacity-10 hover:opacity-20 rounded-full transition-all pointer-events-auto"
					aria-label="Next image"
				>
					<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
					</svg>
				</button>
			{/if}

			<!-- Media container -->
			<div class="relative max-w-7xl max-h-full w-full h-full flex items-center justify-center p-4 ">
				{#if isVideo(data.media[lightboxIndex].mime_type)}
					<!-- Video in lightbox -->
					<video
						src={getVideoUrl(data.media[lightboxIndex])}
						poster={getVideoPoster(data.media[lightboxIndex])}
						controls
						autoplay
						class=" max-h-full object-contain pointer-events-auto"
					>
						<track kind="captions" />
					</video>
				{:else if isImage(data.media[lightboxIndex].mime_type)}
					<!-- Image in lightbox -->
					<img
						src={getImageUrl(data.media[lightboxIndex], 'large')}
						alt={data.media[lightboxIndex].original_filename}
						class="max-w-[70vw] max-h-[85vh] object-contain pointer-events-auto"
					/>
				{/if}

				<!-- Media counter -->
				<div class="absolute bottom-4 left-1/2 transform -translate-x-1/2 bg-black bg-opacity-50 text-white px-4 py-2 rounded-full text-sm">
					{lightboxIndex + 1} / {data.media.length}
				</div>
			</div>
		</div>
	</div>
{/if}
