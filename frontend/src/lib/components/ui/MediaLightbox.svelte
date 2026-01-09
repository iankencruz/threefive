<!-- frontend/src/lib/components/ui/MediaLightbox.svelte -->
<script lang="ts">
	import { X } from 'lucide-svelte';

	interface Media {
		id: string;
		mime_type: string;
		url: string;
		large_url?: string;
		thumbnail_url?: string;
		original_filename: string;
	}

	interface Props {
		media: Media[];
		currentIndex: number;
		open: boolean;
		onclose: () => void;
		onnext?: () => void;
		onprev?: () => void;
	}

	let { media, currentIndex, open, onclose, onnext, onprev }: Props = $props();

	const currentMedia = $derived(media[currentIndex]);

	function isVideo(mimeType: string): boolean {
		return mimeType.startsWith('video/');
	}

	function isImage(mimeType: string): boolean {
		return mimeType.startsWith('image/');
	}

	// Handle keyboard navigation
	function handleKeydown(e: KeyboardEvent) {
		if (!open) return;

		if (e.key === 'Escape') {
			onclose();
		} else if (e.key === 'ArrowRight' && onnext) {
			onnext();
		} else if (e.key === 'ArrowLeft' && onprev) {
			onprev();
		}
	}

	$effect(() => {
		if (open) {
			document.body.style.overflow = 'hidden';
			window.addEventListener('keydown', handleKeydown);
			return () => {
				document.body.style.overflow = '';
				window.removeEventListener('keydown', handleKeydown);
			};
		}
	});
</script>

{#if open && currentMedia}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/90"
		onclick={onclose}
		role="button"
		tabindex="0"
	>
		<!-- Close button -->
		<button
			onclick={onclose}
			class="absolute top-4 right-4 rounded-full bg-white/10 p-2 text-white backdrop-blur-sm transition-colors hover:bg-white/20"
			aria-label="Close lightbox"
		>
			<X size={24} />
		</button>

		<!-- Previous button -->
		{#if onprev && currentIndex > 0}
			<button
				onclick={(e) => {
					e.stopPropagation();
					onprev();
				}}
				class="absolute left-4 rounded-full bg-white/10 p-3 text-white backdrop-blur-sm transition-colors hover:bg-white/20"
				aria-label="Previous image"
			>
				<svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 19l-7-7 7-7"
					/>
				</svg>
			</button>
		{/if}

		<!-- Next button -->
		{#if onnext && currentIndex < media.length - 1}
			<button
				onclick={(e) => {
					e.stopPropagation();
					onnext();
				}}
				class="absolute right-4 rounded-full bg-white/10 p-3 text-white backdrop-blur-sm transition-colors hover:bg-white/20"
				aria-label="Next image"
			>
				<svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
				</svg>
			</button>
		{/if}

		<!-- Media content -->
		<div
			class="relative max-h-[90vh] max-w-[90vw]"
			onclick={(e) => e.stopPropagation()}
			role="button"
			tabindex="0"
		>
			{#if isImage(currentMedia.mime_type)}
				<img
					src={currentMedia.large_url || currentMedia.url}
					alt={currentMedia.original_filename}
					class="max-h-[90vh] max-w-[90vw] object-contain"
				/>
			{:else if isVideo(currentMedia.mime_type)}
				<video src={currentMedia.url} controls autoplay class="max-h-[90vh] max-w-[90vw]">
					<track kind="captions" />
				</video>
			{/if}

			<!-- Media counter -->
			<div
				class="absolute bottom-4 left-1/2 -translate-x-1/2 rounded-full bg-black/50 px-4 py-2 text-sm text-white backdrop-blur-sm"
			>
				{currentIndex + 1} / {media.length}
			</div>
		</div>
	</div>
{/if}
