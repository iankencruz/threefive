<script lang="ts">
	import type { Media } from '$api/media';
	import Navbar from '$components/ui/Navbar.svelte';

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
	const isVideo = (mime: string) => mime?.startsWith('video/');
	const isImage = (mime: string) => mime?.startsWith('image/');
	const getImageUrl = (m: Media) => m.large_url || m.url;
	const getVideoUrl = (m: Media) => m.url;
	const getVideoPoster = (m: Media) => m.thumbnail_url || '';
</script>

<!-- ✨ No more loading state! Media is pre-loaded -->
{#if media}
	<section class="relative overflow-hidden">
		<Navbar variant="ghost" />

		<div class="relative flex min-h-screen items-center">
			<div class="absolute inset-0 bg-black">
				{#if isVideo(media.mime_type)}
					<video
						src={getVideoUrl(media)}
						poster={getVideoPoster(media)}
						autoplay
						muted
						loop
						playsinline
						class="h-full w-full mask-radial-from-100% mask-radial-at-right object-cover"
					>
						<track kind="captions" />
					</video>
				{:else if isImage(media.mime_type)}
					<img
						src={getImageUrl(media)}
						alt={data.title}
						class="h-full w-full mask-radial-from-100% mask-radial-at-right object-cover"
					/>
				{/if}

				<div
					class="absolute inset-0 bg-linear-to-b from-black/85 via-transparent to-black/40"
				></div>
				<div class="absolute inset-0 bg-linear-to-r from-black/70 via-black/5 to-transparent"></div>
			</div>
			<div
				class="@container absolute bottom-16 left-12 z-10 container mx-auto max-w-6xl px-4 text-left"
			>
				<div class=" leading-tight font-normal text-white @max-lg:text-2xl @lg:text-4xl">
					<div class="flex flex-col space-y-2">
						<span>Not just photography.</span>
						<span>Not just branding. </span>
						<span>We builds visual identity. </span>
					</div>
					<a
						href={data.cta_url || '/contact'}
						class="mt-12 inline-flex max-w-max items-center gap-2 rounded-lg bg-white px-6 py-3 text-lg font-semibold text-gray-900 shadow-lg transition-all hover:scale-105 hover:bg-gray-100 hover:shadow-xl"
					>
						{data.cta_text || 'Book Now'}
					</a>
				</div>
			</div>

			<div
				class="@container absolute right-12 bottom-16 z-10 container mx-auto max-w-6xl px-4 text-right"
			>
				<div class="w-full text-white">
					<h1
						class=" leading-tight font-bold text-primary @max-lg:text-2xl @lg:@max-2xl:text-5xl @2xl:text-7xl @5xl:text-9xl"
					>
						{data.title}
					</h1>

					{#if data.subtitle}
						<p
							class=" leading-tight font-bold @max-lg:text-2xl @lg:@max-2xl:text-4xl @2xl:text-6xl @5xl:text-8xl"
						>
							{data.subtitle}
						</p>
					{/if}
				</div>
			</div>
		</div>
	</section>
{:else}
	<section class="relative flex min-h-screen items-center justify-center bg-gray-900">
		<p class="text-xl text-white">{data.title} (No background media specified)</p>
	</section>
{/if}
