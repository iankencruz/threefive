<script lang="ts">
	import type { Media } from '$api/media';
	import DynamicForm, { type FormConfig } from '$components/ui/DynamicForm.svelte';

	export interface FeatureBlockData {
		title: string;
		description: string;
		heading: string;
		subheading: string;
		media_ids?: string[];
		media?: Media[];
	}

	interface Props {
		data: FeatureBlockData;
		onchange?: (data: FeatureBlockData) => void;
	}

	let {
		data = $bindable({
			title: '',
			description: '',
			heading: '',
			subheading: '',
			media_ids: []
		}),
		onchange
	}: Props = $props();

	const formConfig: FormConfig = {
		fields: [
			{
				name: 'title',
				label: 'Title',
				type: 'text',
				required: true,
				colSpan: 12
			},
			{
				name: 'description',
				label: 'Description',
				type: 'textarea',
				required: true,
				colSpan: 12
			},
			{
				name: 'heading',
				label: 'Heading',
				type: 'text',
				required: true,
				colSpan: 12
			},
			{
				name: 'subheading',
				label: 'Subheading',
				type: 'text',
				colSpan: 12
			},
			{
				name: 'media_ids',
				type: 'media',
				label: 'Feature Images',
				colSpan: 12,
				multiple: true, // NEW: Enable multi-select
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
		data = {
			...updatedData,
			media: updatedData.media_ids?.map((id: string) => initialMediaCache[id]).filter(Boolean) || []
		} as FeatureBlockData;
		onchange?.(data);
	};
</script>

<DynamicForm
	asForm={false}
	config={formConfig}
	formData={data}
	onchange={handleChange}
	{initialMediaCache}
/>
