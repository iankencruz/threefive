<script lang="ts">
	import { uploadMedia } from '$lib/api/media';
	import { Toaster, toast } from 'svelte-sonner';

	let {
		open = false,
		accept = '*/*',
		maxSize = 10 * 1024 * 1024,
		multiple = true,
		onuploaded,
		onclose
	} = $props();

	let files: File[] = $state([]);
	let uploading = $state(false);
	let progress: Record<string, number> = $state({});
	let errors: Record<string, string> = $state({});

	function handleDrop(event: DragEvent) {
		event.preventDefault();
		if (!event.dataTransfer?.files) return;
		addFiles(event.dataTransfer.files);
	}

	function addFiles(fileList: FileList) {
		const valid = Array.from(fileList).filter((f) => f.size <= maxSize);
		files = [...files, ...valid];
	}

	function removeFile(index: number) {
		files = files.toSpliced(index, 1);
	}

	function handleCancel() {
		files = [];
		progress = {};
		errors = {};
		onclose?.();
	}

	async function handleUpload() {
		uploading = true;
		errors = {};
		const uploaded: File[] = [];

		for (const file of files) {
			let success = false;
			progress[file.name] = 0;

			for (let attempt = 1; attempt <= 2 && !success; attempt++) {
				try {
					await uploadMedia(file, (percent: number) => {
						progress[file.name] = percent;
					});
					success = true;
					uploaded.push(file);
				} catch (err) {
					if (attempt === 2) {
						errors[file.name] = 'Upload failed after retry';
						toast.error(`Upload failed: ${file.name}`);
					}
				}
			}
		}

		if (uploaded.length > 0) {
			toast.success(`${uploaded.length} file(s) uploaded successfully`);
			onuploaded?.({ files: uploaded });
		}

		uploading = false;
	}
</script>

{#if open}
	<div
		role="dialog"
		tabindex="0"
		aria-modal="true"
		aria-labelledby="upload_modal"
		aria-describedby="dialog_desc"
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		onclick={() => onclose?.()}
		onkeydown={(e) => {
			if (e.key === 'Escape') onclose?.();
		}}
	>
		<div
			role="button"
			tabindex="0"
			aria-label="Close modal"
			class="relative w-full max-w-lg rounded-lg bg-white p-6 shadow-xl"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => {
				if (e.key === 'Escape' || e.key === 'Enter' || e.key === ' ') onclose?.();
			}}
		>
			<div class="mb-8 flex items-center justify-between">
				<h2 class=" text-lg font-semibold">Upload Files</h2>

				<!-- Close button -->
				<button
					class="cursor-pointer text-2xl text-gray-400 hover:text-red-700"
					onclick={handleCancel}
					aria-label="Close"
				>
					&times;
				</button>
			</div>

			<div
				role="region"
				class="relative mb-4 flex h-48 w-full cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed border-gray-300 text-center text-sm text-gray-500"
				ondrop={handleDrop}
				ondragover={(e) => e.preventDefault()}
			>
				<label class="flex h-full w-full cursor-pointer flex-col items-center justify-center">
					<p>Drop files here or click to select</p>
					<input
						type="file"
						{accept}
						{multiple}
						onchange={(e) => {
							const target = e.target as HTMLInputElement;
							if (target.files) addFiles(target.files);
						}}
						class="absolute inset-0 cursor-pointer opacity-0"
					/>a
				</label>
			</div>

			{#if files.length > 0}
				<ul class="max-h-48 space-y-2 overflow-y-auto">
					{#each files as file, i}
						<li class="flex flex-col gap-1 rounded bg-gray-100 px-3 py-2">
							<div class="flex items-center justify-between gap-4">
								<div class="truncate">
									<strong>{file.name}</strong>
									<span class="text-xs text-gray-600">
										({(file.size / 1024 / 1024).toFixed(2)} MB)
									</span>
								</div>
								{#if uploading}
									<div class="h-2 w-24 rounded bg-gray-300">
										<div
											class="h-2 rounded bg-indigo-600"
											style={`width: ${progress[file.name] ?? 0}%`}
										></div>
									</div>
								{:else}
									<button class="text-sm text-red-600" onclick={() => removeFile(i)}>
										Remove
									</button>
								{/if}
							</div>
							{#if errors[file.name]}
								<p class="text-sm text-red-500">{errors[file.name]}</p>
							{/if}
						</li>
					{/each}
				</ul>
			{/if}

			<div class="mt-6 flex justify-end gap-4">
				<button
					onclick={handleCancel}
					class="rounded bg-gray-200 px-4 py-2 text-sm hover:bg-gray-300"
					disabled={uploading}
				>
					Cancel
				</button>
				<button
					onclick={handleUpload}
					class="rounded bg-indigo-600 px-4 py-2 text-sm text-white hover:bg-indigo-500"
					disabled={uploading || files.length === 0}
				>
					{uploading ? 'Uploading...' : 'Upload'}
				</button>
			</div>
		</div>
	</div>
{/if}
