<script lang="ts">
interface Media {
	id: string;
	filename: string;
	original_filename: string;
	mime_type: string;
	size_bytes: number;
	width?: number;
	height?: number;
	url: string;
	thumbnail_url?: string;
	created_at: string;
}

interface Props {
	item: Media;
	selected: boolean;
	onToggleSelect: () => void;
}

let { item, selected, onToggleSelect }: Props = $props();

function formatFileSize(bytes: number): string {
	if (bytes === 0) return "0 Bytes";
	const k = 1024;
	const sizes = ["Bytes", "KB", "MB", "GB"];
	const i = Math.floor(Math.log(bytes) / Math.log(k));
	return Math.round((bytes / k ** i) * 100) / 100 + " " + sizes[i];
}

function getFileIcon(mimeType: string): string {
	if (mimeType.startsWith("image/")) return "🖼";
	if (mimeType.startsWith("video/")) return "🎥";
	if (mimeType.includes("pdf")) return "📄";
	if (mimeType.includes("document")) return "📝";
	return "📎";
}

function copyUrl() {
	navigator.clipboard.writeText(item.url);
	alert("URL copied to clipboard!");
}
</script>

<div 
	class="group relative bg-white border-2 rounded-lg overflow-hidden transition-all hover:shadow-lg"
	class:border-blue-500={selected}
	class:ring-2={selected}
	class:ring-blue-200={selected}
	class:border-gray-200={!selected}
>	<!-- Selection Checkbox -->
	<div class="absolute top-3 left-3 z-10">
		<input
			type="checkbox"
			checked={selected}
			onchange={onToggleSelect}
			class="w-5 h-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500 cursor-pointer"
		/>
	</div>

	<!-- Preview -->
	<div class="relative h-48 bg-gray-100 flex items-center justify-center overflow-hidden">
		{#if item.mime_type.startsWith('image/')}
			<img
				src={item.url}
				alt={item.original_filename}
				class="w-full h-full object-cover"
			/>
		{:else if item.mime_type.startsWith('video/')}
			<video src={item.url} class="w-full h-full object-cover">
				<track kind="captions" />
			</video>
			<div class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-30">
				<svg class="w-16 h-16 text-white opacity-80" fill="currentColor" viewBox="0 0 20 20">
					<path d="M6.3 2.841A1.5 1.5 0 004 4.11V15.89a1.5 1.5 0 002.3 1.269l9.344-5.89a1.5 1.5 0 000-2.538L6.3 2.84z" />
				</svg>
			</div>
		{:else}
			<div class="text-6xl">{getFileIcon(item.mime_type)}</div>
		{/if}
	</div>

	<!-- Info -->
	<div class="p-4">
		<h3 class="font-semibold text-sm text-gray-900 truncate mb-1" title={item.original_filename}>
			{item.original_filename}
		</h3>
		<div class="flex items-center gap-2 text-xs text-gray-500">
			<span>{formatFileSize(item.size_bytes)}</span>
			{#if item.width && item.height}
				<span>•</span>
				<span>{item.width} × {item.height}</span>
			{/if}
		</div>
	</div>

	<!-- Hover Actions -->
	<div class="absolute bottom-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity">
		<button
			onclick={copyUrl}
			class="p-2 bg-white rounded-lg shadow-lg hover:bg-gray-50 transition-colors"
			title="Copy URL"
		>
			<svg class="w-4 h-4 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
			</svg>
		</button>
	</div>
</div>
