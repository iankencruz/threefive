<!-- frontend/src/lib/components/blocks/forms/GalleryBlockForm.svelte -->
<script lang="ts">
	import type { Media } from '$api/media';
	import type { FormConfig } from '$components/ui/form';
	import DynamicForm from '$components/ui/form/DynamicForm.svelte';

	export interface GalleryBlockData {
		title?: string;
		media_ids?: string[];
		media?: Media[];
	}

	interface Props {
		data: GalleryBlockData;
		onchange?: (data: GalleryBlockData) => void;
	}

	let { data, onchange }: Props = $props();

	const formConfig: FormConfig = {
		fields: [
			{
				name: 'title',
				type: 'text',
				label: 'Gallery Title',
				placeholder: 'e.g., Summer Vacation Photos',
				colSpan: 12
			},
			{
				name: 'media_ids',
				type: 'media',
				label: 'Images',
				colSpan: 12,
				multiple: true,
				required: true
			}
		]
	};

	// Build initialMediaCache from data.media array
	const initialMediaCache = $derived<Record<string, Media>>(
		data?.media?.reduce(
			(acc, media) => {
				acc[media.id] = media;
				return acc;
			},
			{} as Record<string, Media>
		) || {}
	);

	const handleChange = (updatedData: Record<string, any>) => {
		const newData = {
			title: updatedData.title,
			media_ids: updatedData.media_ids || []
			// Backend will populate media array when fetching
		};
		onchange?.(newData);
	};
</script>

<DynamicForm
	asForm={false}
	config={formConfig}
	formData={data}
	onchange={handleChange}
	{initialMediaCache}
/>
