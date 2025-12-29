<!-- frontend/src/lib/components/blocks/forms/HeroBlockForm.svelte -->
<script lang="ts">
	import type { Media } from '$api/media';
	import DynamicForm from '$components/ui/form/DynamicForm.svelte';

	export interface HeroBlockData {
		title: string;
		subtitle?: string;
		image_id?: string | null;
		cta_text?: string;
		cta_url?: string;
		media?: Media;
	}

	interface Props {
		data: HeroBlockData;
		onchange?: (data: HeroBlockData) => void;
	}

	let {
		data = $bindable({
			title: '',
			subtitle: '',
			image_id: null, // Use null instead of empty string
			cta_text: '',
			cta_url: ''
		}),
		onchange
	}: Props = $props();

	// Move formConfig outside of reactive context
	const formConfig: FormConfig = {
		fields: [
			{
				name: 'title',
				type: 'text',
				label: 'Title',
				placeholder: 'Enter hero title',
				required: true,
				colSpan: 12
			},
			{
				name: 'subtitle',
				type: 'text',
				label: 'Subtitle',
				placeholder: 'Enter subtitle (optional)',
				colSpan: 12
			},

			{
				name: 'cta_text',
				type: 'text',
				label: 'CTA Text',
				placeholder: 'Button text (optional)',
				colSpan: 6
			},
			{
				name: 'cta_url',
				type: 'text',
				label: 'CTA URL',
				placeholder: 'Button link (optional)',
				colSpan: 6
			},
			{
				name: 'image_id',
				type: 'media',
				label: 'Background Image',
				colSpan: 12
			}
		]
	};

	const initialMediaCache = $derived<Record<string, Media>>(
		data?.media && data.image_id ? { [data.image_id]: data.media } : {}
	);

	const handleChange = (updatedData: Record<string, any>) => {
		data = updatedData as HeroBlockData;
		onchange?.(data);
	};
</script>

<DynamicForm
	config={formConfig}
	formData={data}
	onchange={handleChange}
	asForm={false}
	{initialMediaCache}
/>
