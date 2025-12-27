<script lang="ts">
	import type { Media } from '$api/media';
	import DynamicForm, { type FormConfig } from '$components/ui/DynamicForm.svelte';

	export interface FeatureBlockData {
		title: string;
		description?: string;
		heading?: string;
		subheading?: string;
		image_id?: string | null;
		media?: Media;
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
			image_id: null
		}),
		onchange
	}: Props = $props();

	const formConfig: FormConfig = {
		fields: [
			{
				name: 'title',
				label: 'Title',
				type: 'text',
				required: true
			},
			{
				name: 'description',
				label: 'Description',
				type: 'textarea'
			},
			{
				name: 'heading',
				label: 'Heading',
				type: 'text'
			},
			{
				name: 'subheading',
				label: 'Subheading',
				type: 'text'
			},
			{
				name: 'image_id',
				type: 'media',
				label: 'Featured Image',
				colSpan: 12
			},
			{
				name: 'status',
				type: 'checkbox',
				label: 'Active'
			}
		]
	};

	// Build initialMediaCache from data.media
	const initialMediaCache = $derived<Record<string, Media>>(
		data?.media && data.image_id ? { [data.image_id]: data.media } : {}
	);

	const handleChange = (updatedData: Record<string, any>) => {
		data = updatedData as FeatureBlockData;
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
