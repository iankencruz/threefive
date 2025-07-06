<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { uploadMediaAndLinkToContext } from '$lib/api/media';
	import type { MediaItem } from '$lib/types';

	let { context, onuploaded } = $props<{
		context: { type: 'project' | 'block'; id: string };
		onuploaded: (media: MediaItem) => void;
	}>();

	let file: File | null = null;
	let uploading = $state(false);

	async function handleUpload() {
		if (!file) {
			toast.error('Please select a file');
			return;
		}

		try {
			uploading = true;
			const uploaded = await uploadMediaAndLinkToContext(file, context);
			onuploaded(uploaded);
			file = null;
		} catch (err) {
			console.error(err);
			toast.error('Upload failed');
		} finally {
			uploading = false;
		}
	}
</script>

<div class="flex flex-col gap-4">
	<input
		type="file"
		accept="image/*"
		class="block w-full rounded border px-3 py-2"
		onchange={(e) => (file = (e.target as HTMLInputElement).files?.[0] || null)}
	/>

	<button
		onclick={handleUpload}
		disabled={uploading}
		class="inline-flex items-center justify-center rounded bg-black px-4 py-2 text-white hover:bg-gray-800 disabled:opacity-50"
	>
		{uploading ? 'Uploading...' : 'Upload & Link'}
	</button>
</div>
