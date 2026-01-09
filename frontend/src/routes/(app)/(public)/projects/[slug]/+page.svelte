<!-- frontend/src/routes/(app)/(public)/projects/[slug]/+page.svelte -->
<script lang="ts">
	import type { PageData } from './$types';
	import { Calendar, ExternalLink, User } from 'lucide-svelte';
	import MediaLightbox from '$lib/components/ui/MediaLightbox.svelte';

	let { data }: { data: PageData } = $props();

	// Lightbox state
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
</script>

<svelte:head>
	<!-- SEO Meta Tags -->
	<title>{data.project.title || data.project.seo?.meta_title}</title>
	<meta
		name="description"
		content={data.project.seo?.meta_description || data.project.description || ''}
	/>

	<!-- Open Graph -->
	{#if data.project.seo?.og_title}
		<meta property="og:title" content={data.project.seo.og_title} />
	{/if}
	{#if data.project.seo?.og_description}
		<meta property="og:description" content={data.project.seo.og_description} />
	{/if}

	<!-- Robots -->
	{#if data.project.seo}
		<meta
			name="robots"
			content="{data.project.seo.robots_index ? 'index' : 'noindex'}, {data.project.seo
				.robots_follow
				? 'follow'
				: 'nofollow'}"
		/>
	{/if}

	<!-- Canonical URL -->
	{#if data.project.seo?.canonical_url}
		<link rel="canonical" href={data.project.seo.canonical_url} />
	{/if}
</svelte:head>

<!-- Page Content -->
<div class="min-h-screen bg-background">
	<!-- Project Header -->
	<section class="border-b border-gray-200 bg-background py-16">
		<div class="container mx-auto max-w-6xl px-4">
			<h1 class="mb-4 text-4xl font-bold text-gray-100 md:text-5xl">{data.project.title}</h1>
			{#if data.project.description}
				<p class="text-xl text-gray-200">{data.project.description}</p>
			{/if}
		</div>
	</section>

	<!-- Project Media Gallery -->
	{#if data.projectMedia && data.projectMedia.length > 0}
		<section class="py-12">
			<div class="container mx-auto max-w-7xl px-4">
				<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
					{#each data.projectMedia as media, index}
						<button
							onclick={() => openLightbox(index)}
							class="group relative aspect-square overflow-hidden rounded-lg bg-gray-100 transition-transform hover:scale-[1.02]"
						>
							{#if isImage(media.mime_type)}
								<img
									src={media.medium_url || media.url}
									alt={media.original_filename}
									class="h-full w-full object-cover transition-opacity group-hover:opacity-90"
								/>
							{:else if isVideo(media.mime_type)}
								<video
									src={media.url}
									poster={media.thumbnail_url}
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
					{/each}
				</div>
			</div>
		</section>
	{/if}

	<!-- Project Metadata Section -->
	<section class=" bg-surface py-16">
		<div class="container mx-auto max-w-4xl px-4">
			<h2 class="mb-8 text-center text-2xl font-bold text-gray-100">Project Details</h2>

			<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
				<!-- Client Name -->
				{#if data.project.client_name}
					<div class="rounded-lg bg-white p-6 shadow-sm">
						<div class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-500">
							<User size={16} />
							<h3>Client</h3>
						</div>
						<p class="text-lg font-semibold text-gray-900">{data.project.client_name}</p>
					</div>
				{/if}

				<!-- Project Year -->
				{#if data.project.project_year}
					<div class="rounded-lg bg-white p-6 shadow-sm">
						<div class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-500">
							<Calendar size={16} />
							<h3>Year</h3>
						</div>
						<p class="text-lg font-semibold text-gray-900">{data.project.project_year}</p>
					</div>
				{/if}

				<!-- Technologies -->
				{#if data.project.technologies && data.project.technologies.length > 0}
					<div class="rounded-lg bg-white p-6 shadow-sm md:col-span-2">
						<h3 class="mb-3 text-sm font-medium text-gray-500">Technologies</h3>
						<div class="flex flex-wrap gap-2">
							{#each data.project.technologies as tech}
								<span class="rounded-full bg-blue-100 px-3 py-1 text-sm font-medium text-blue-800">
									{tech}
								</span>
							{/each}
						</div>
					</div>
				{/if}

				<!-- Project URL -->
				{#if data.project.project_url}
					<div class="rounded-lg bg-white p-6 shadow-sm md:col-span-2">
						<div class="mb-2 flex items-center gap-2 text-sm font-medium text-gray-500">
							<ExternalLink size={16} />
							<h3>Live Project</h3>
						</div>
						<a
							href={data.project.project_url}
							target="_blank"
							rel="noopener noreferrer"
							class="inline-flex items-center gap-2 text-lg font-semibold text-blue-600 transition-colors hover:text-blue-800"
						>
							Visit Website
							<ExternalLink size={18} />
						</a>
					</div>
				{/if}

				<!-- Published Date -->
				{#if data.project.published_at}
					<div class="rounded-lg bg-white p-6 shadow-sm md:col-span-2">
						<h3 class="mb-2 text-sm font-medium text-gray-500">Published</h3>
						<p class="text-lg font-semibold text-gray-900">
							{formatDate(data.project.published_at)}
						</p>
					</div>
				{/if}
			</div>
		</div>
	</section>
</div>

<!-- Lightbox Component -->
<MediaLightbox
	media={data.projectMedia || []}
	currentIndex={lightboxIndex}
	open={lightboxOpen}
	onclose={closeLightbox}
	onnext={nextImage}
	onprev={prevImage}
/>
