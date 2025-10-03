<script lang="ts">
import { mediaApi, type ErrorResponse } from "$api/media";

interface Props {
	onClose: () => void;
	onComplete: () => void;
}

let { onClose, onComplete }: Props = $props();

let dragOver = $state(false);
let uploading = $state(false);
let uploadProgress = $state(0);
let selectedFiles = $state<File[]>([]);
let error = $state("");

function handleDragOver(e: DragEvent) {
	e.preventDefault();
	dragOver = true;
}

function handleDragLeave() {
	dragOver = false;
}

function handleDrop(e: DragEvent) {
	e.preventDefault();
	dragOver = false;

	const files = Array.from(e.dataTransfer?.files || []);
	selectedFiles = files;
}

function handleFileSelect(e: Event) {
	const input = e.target as HTMLInputElement;
	const files = Array.from(input.files || []);
	selectedFiles = files;
}

async function uploadFiles() {
	if (selectedFiles.length === 0) return;

	uploading = true;
	error = "";

	try {
		for (let i = 0; i < selectedFiles.length; i++) {
			await mediaApi.uploadMedia(selectedFiles[i]);
			uploadProgress = ((i + 1) / selectedFiles.length) * 100;
		}

		onComplete();
	} catch (err) {
		const errorResponse = err as ErrorResponse;
		error = errorResponse.message || "Upload failed";
	} finally {
		uploading = false;
		uploadProgress = 0;
	}
}

function removeFile(index: number) {
	selectedFiles.splice(index, 1);
	selectedFiles = selectedFiles;
}

function formatFileSize(bytes: number): string {
	if (bytes === 0) return "0 Bytes";
	const k = 1024;
	const sizes = ["Bytes", "KB", "MB", "GB"];
	const i = Math.floor(Math.log(bytes) / Math.log(k));
	return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + " " + sizes[i];
}
</script>

<!-- Modal Container -->
<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
	<!-- Backdrop - Separate element with opacity -->
	<div
		class="absolute inset-0 bg-black/50"
		onclick={onClose}
		role="button"
		tabindex="-1"
		aria-label="Close modal"
	></div>

	<!-- Modal Content - No opacity applied -->
	<div
		class="relative bg-white rounded-xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto"
		onclick={(e) => e.stopPropagation()}
		role="dialog"
		aria-modal="true"
		tabindex="-1"
	>
		<!-- Header -->
		<div class="flex items-center justify-between pt-6 px-6 ">
			<h2 class="text-2xl font-bold text-gray-900">Upload Media</h2>
			<button
				onclick={onClose}
				class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
				aria-label="Close"
			>
				<svg class="w-6 h-6 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
				</svg>
			</button>
		</div>

		<!-- Content -->
		<div class="p-6">
			{#if selectedFiles.length === 0}
				<!-- Drop Zone -->
				<div
					ondragover={handleDragOver}
					ondragleave={handleDragLeave}
					ondrop={handleDrop}
					class="border-2 border-dashed rounded-xl p-12 text-center transition-colors"
					class:border-blue-500={dragOver}
					class:bg-blue-50={dragOver}
					class:border-gray-300={!dragOver}
					class:hover:border-gray-400={!dragOver}
				>
					<svg class="w-16 h-16 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
					</svg>
					
					<h3 class="text-lg font-semibold text-gray-700 mb-2">
						Drop files here or click to browse
					</h3>
					<p class="text-sm text-gray-500 mb-4">
						Support for images, videos, and documents (Max 50MB)
					</p>
					
					<input
						type="file"
						multiple
						accept="image/*,video/*,.pdf,.doc,.docx"
						onchange={handleFileSelect}
						class="hidden"
						id="file-input"
					/>
					<label
						for="file-input"
						class="inline-block px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 cursor-pointer transition-colors"
					>
						Select Files
					</label>
				</div>
			{:else}
				<!-- File List -->
				<div class="space-y-3 mb-6">
					{#each selectedFiles as file, i}
						<div class="flex items-center gap-4 p-4 bg-gray-50 rounded-lg">
							<div class="flex-shrink-0 w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
								<svg class="w-6 h-6 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
									<path fill-rule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clip-rule="evenodd" />
								</svg>
							</div>
							
							<div class="flex-1 min-w-0">
								<p class="font-medium text-gray-900 truncate">{file.name}</p>
								<p class="text-sm text-gray-500">{formatFileSize(file.size)}</p>
							</div>
							
							{#if !uploading}
								<button
									onclick={() => removeFile(i)}
									class="p-2 hover:bg-gray-200 rounded-lg transition-colors"
									aria-label="Remove file"
								>
									<svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
									</svg>
								</button>
							{/if}
						</div>
					{/each}
				</div>

				<!-- Progress Bar -->
				{#if uploading}
					<div class="mb-6">
						<div class="flex justify-between text-sm text-gray-600 mb-2">
							<span>Uploading...</span>
							<span>{Math.round(uploadProgress)}%</span>
						</div>
						<div class="h-2 bg-gray-200 rounded-full overflow-hidden">
							<div
								class="h-full bg-blue-600 transition-all duration-300"
								style="width: {uploadProgress}%"
							></div>
						</div>
					</div>
				{/if}

				<!-- Error Message -->
				{#if error}
					<div class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
						<p class="text-sm text-red-700">{error}</p>
					</div>
				{/if}

				<!-- Actions -->
				<div class="flex gap-3 justify-end">
					<button
						onclick={() => selectedFiles = []}
						disabled={uploading}
						class="px-6 py-2.5 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
					>
						Clear
					</button>
					<button
						onclick={uploadFiles}
						disabled={uploading}
						class="px-6 py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
					>
						{#if uploading}
							<div class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
						{/if}
						{uploading ? 'Uploading...' : `Upload ${selectedFiles.length} file${selectedFiles.length > 1 ? 's' : ''}`}
					</button>
				</div>
			{/if}
		</div>
	</div>
</div>
