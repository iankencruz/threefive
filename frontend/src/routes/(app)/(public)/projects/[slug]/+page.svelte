<!-- frontend/src/routes/(app)/(public)/projects/[slug]/+page.svelte -->
<script lang="ts">
	import type { PageData } from './$types';
	import { Calendar, ChevronLeft, ChevronRight, ExternalLink, User } from 'lucide-svelte';
	import MediaLightbox from '$lib/components/ui/MediaLightbox.svelte';

	let { data }: { data: PageData } = $props();

	// Slider state
	let currentSlide = $state(0);
	let lightboxOpen = $state(false);
	let lightboxIndex = $state(0);

	// Format date helper
	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	};

	function isVideo(mimeType: string): boolean {
		return mimeType.startsWith('video/');
	}

	function isImage(mimeType: string): boolean {
		return mimeType.startsWith('image/');
	}

	function nextSlide() {
		if (data.projectMedia) {
			currentSlide = (currentSlide + 1) % data.projectMedia.length;
		}
	}

	function prevSlide() {
		if (data.projectMedia) {
			currentSlide = currentSlide === 0 ? data.projectMedia.length - 1 : currentSlide - 1;
		}
	}

	function goToSlide(index: number) {
		currentSlide = index;
	}

	function openLightbox(index: number) {
		lightboxIndex = index;
		lightboxOpen = true;
	}

	function closeLightbox() {
		lightboxOpen = false;
	}

	function nextImage() {
		if (data.projectMedia && lightboxIndex < data.projectMedia.length - 1) {
			lightboxIndex++;
		}
	}

	function prevImage() {
		if (lightboxIndex > 0) {
			lightboxIndex--;
		}
	}

	// Handle keyboard navigation for slider
	function handleKeydown(e: KeyboardEvent) {
		if (lightboxOpen) return; // Don't interfere with lightbox

		if (e.key === 'ArrowRight') {
			nextSlide();
		} else if (e.key === 'ArrowLeft') {
			prevSlide();
		}
	}

	$effect(() => {
		window.addEventListener('keydown', handleKeydown);
		return () => {
			window.removeEventListener('keydown', handleKeydown);
		};
	});
</script>

<svelte:head>
	<title>TFP - {data.project.title || data.project.seo?.meta_title}</title>
	{#if data.project.seo?.meta_description || data.project.description}
		<meta
			name="description"
			content={data.project.seo?.meta_description || data.project.description}
		/>
	{/if}
	{#if data.project.seo?.og_title}
		<meta property="og:title" content={data.project.seo.og_title} />
	{/if}
	{#if data.project.seo?.og_description}
		<meta property="og:description" content={data.project.seo.og_description} />
	{/if}
	{#if data.project.seo}
		<meta
			name="robots"
			content="{data.project.seo.robots_index ? 'index' : 'noindex'}, {data.project.seo
				.robots_follow
				? 'follow'
				: 'nofollow'}"
		/>
	{/if}
	{#if data.project.seo?.canonical_url}
		<link rel="canonical" href={data.project.seo.canonical_url} />
	{/if}
</svelte:head>

<!-- Project Media Slider -->
{#if data.projectMedia && data.projectMedia.length > 0}
	{@const currentMedia = data.projectMedia[currentSlide]}
	{@const secondaryMedia = data.projectMedia[(currentSlide + 1) % data.projectMedia.length]}
	{@const tertiaryMedia = data.projectMedia[(currentSlide + 2) % data.projectMedia.length]}
	<section class="-mx-4 h-[calc(100vh-5rem)] bg-background px-4 sm:-mx-6 lg:-mx-8">
		<div class="flex h-full gap-4 p-4">
			<!-- Left: Primary Image (2 columns) -->
			<div class="relative flex flex-[2]">
				<button
					onclick={() => openLightbox(currentSlide)}
					class="group h-full w-full overflow-hidden rounded-lg bg-gray-800"
				>
					{#if isImage(currentMedia.mime_type)}
						<img
							src={currentMedia.large_url || currentMedia.url}
							alt={currentMedia.original_filename}
							class="h-full w-full object-cover transition-transform group-hover:scale-105"
						/>
					{:else if isVideo(currentMedia.mime_type)}
						<video
							src={currentMedia.url}
							poster={currentMedia.thumbnail_url}
							class="h-full w-full object-cover"
							muted
						>
							<track kind="captions" />
						</video>
						<!-- Video play icon overlay -->
						<div class="pointer-events-none absolute inset-0 flex items-center justify-center">
							<div class="rounded-full bg-black/50 p-6 text-white backdrop-blur-sm">
								<svg class="h-12 w-12" fill="currentColor" viewBox="0 0 24 24">
									<path d="M8 5v14l11-7z" />
								</svg>
							</div>
						</div>
					{/if}
				</button>

				<!-- Navigation Arrows -->
				{#if data.projectMedia.length > 1}
					<button
						onclick={prevSlide}
						class="absolute top-1/2 left-4 z-10 -translate-y-1/2 rounded-full bg-white/10 p-3 text-white backdrop-blur-sm transition-all hover:bg-white/20"
						aria-label="Previous image"
					>
						<ChevronLeft size={24} />
					</button>

					<button
						onclick={nextSlide}
						class="absolute top-1/2 right-4 z-10 -translate-y-1/2 rounded-full bg-white/10 p-3 text-white backdrop-blur-sm transition-all hover:bg-white/20"
						aria-label="Next image"
					>
						<ChevronRight size={24} />
					</button>
				{/if}
			</div>

			<!-- Right: Project Details + Secondary + Tertiary Images (1 column) -->
			<div class="flex min-h-0 flex-1 flex-col gap-4">
				<!-- Project Details Card -->
				<div
					class="max-h-[40%] flex-shrink-0 overflow-y-auto rounded-lg border border-gray-800 bg-gray-950 p-6"
				>
					<h2 class="mb-4 text-2xl font-bold text-white">{data.project.title}</h2>

					<div class="space-y-4">
						{#if data.project.description}
							<p class="text-sm leading-relaxed text-gray-400">{data.project.description}</p>
						{/if}

						{#if data.project.client_name}
							<div>
								<div class="mb-1 text-xs font-medium tracking-wide text-gray-500 uppercase">
									Client
								</div>
								<div class="text-sm text-white">{data.project.client_name}</div>
							</div>
						{/if}

						{#if data.project.project_year}
							<div>
								<div class="mb-1 text-xs font-medium tracking-wide text-gray-500 uppercase">
									Year
								</div>
								<div class="text-sm text-white">{data.project.project_year}</div>
							</div>
						{/if}

						{#if data.project.technologies && data.project.technologies.length > 0}
							<div>
								<div class="mb-2 text-xs font-medium tracking-wide text-gray-500 uppercase">
									Technologies
								</div>
								<div class="flex flex-wrap gap-2">
									{#each data.project.technologies as tech}
										<span class="rounded bg-gray-800 px-2 py-1 text-xs text-gray-300">
											{tech}
										</span>
									{/each}
								</div>
							</div>
						{/if}

						{#if data.project.project_url}
							<div>
								<div class="mb-1 text-xs font-medium tracking-wide text-gray-500 uppercase">
									Live Project
								</div>
								<a
									href={data.project.project_url}
									target="_blank"
									rel="noopener noreferrer"
									class="inline-flex items-center gap-1 text-sm text-blue-400 transition-colors hover:text-blue-300"
								>
									Visit Website
									<ExternalLink size={14} />
								</a>
							</div>
						{/if}
					</div>

					<!-- Image Counter -->
					<div class="mt-6 border-t border-gray-800 pt-4 text-center">
						<div class="text-sm text-gray-400">
							{currentSlide + 1} / {data.projectMedia.length}
						</div>
					</div>
				</div>

				<!-- Secondary and Tertiary Images in 2 rows -->
				{#if data.projectMedia.length > 1}
					<div class="flex min-h-0 flex-1 flex-col gap-4">
						<!-- Secondary Image (Top) -->
						<button
							onclick={nextSlide}
							class="group relative min-h-0 flex-1 overflow-hidden rounded-lg bg-gray-800"
						>
							{#if isImage(secondaryMedia.mime_type)}
								<img
									src={secondaryMedia.large_url || secondaryMedia.url}
									alt={secondaryMedia.original_filename}
									class="h-full w-full object-cover transition-transform group-hover:scale-105"
								/>
							{:else if isVideo(secondaryMedia.mime_type)}
								<video
									src={secondaryMedia.url}
									poster={secondaryMedia.thumbnail_url}
									class="h-full w-full object-cover"
									muted
								>
									<track kind="captions" />
								</video>
								<!-- Video play icon overlay -->
								<div class="pointer-events-none absolute inset-0 flex items-center justify-center">
									<div class="rounded-full bg-black/50 p-4 text-white backdrop-blur-sm">
										<svg class="h-8 w-8" fill="currentColor" viewBox="0 0 24 24">
											<path d="M8 5v14l11-7z" />
										</svg>
									</div>
								</div>
							{/if}
						</button>

						<!-- Tertiary Image (Bottom) -->
						{#if data.projectMedia.length > 2}
							<button
								onclick={() => {
									currentSlide = (currentSlide + 2) % data.projectMedia.length;
								}}
								class="group relative min-h-0 flex-1 overflow-hidden rounded-lg bg-gray-800"
							>
								{#if isImage(tertiaryMedia.mime_type)}
									<img
										src={tertiaryMedia.large_url || tertiaryMedia.url}
										alt={tertiaryMedia.original_filename}
										class="h-full w-full object-cover transition-transform group-hover:scale-105"
									/>
								{:else if isVideo(tertiaryMedia.mime_type)}
									<video
										src={tertiaryMedia.url}
										poster={tertiaryMedia.thumbnail_url}
										class="h-full w-full object-cover"
										muted
									>
										<track kind="captions" />
									</video>
									<!-- Video play icon overlay -->
									<div
										class="pointer-events-none absolute inset-0 flex items-center justify-center"
									>
										<div class="rounded-full bg-black/50 p-4 text-white backdrop-blur-sm">
											<svg class="h-8 w-8" fill="currentColor" viewBox="0 0 24 24">
												<path d="M8 5v14l11-7z" />
											</svg>
										</div>
									</div>
								{/if}
							</button>
						{/if}
					</div>
				{/if}
			</div>
		</div>
	</section>
{/if}

<!-- Lightbox Component -->
<MediaLightbox
	media={data.projectMedia || []}
	currentIndex={lightboxIndex}
	open={lightboxOpen}
	onclose={closeLightbox}
	onnext={nextImage}
	onprev={prevImage}
/>
