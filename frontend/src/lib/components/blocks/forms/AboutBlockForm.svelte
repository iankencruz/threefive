<script lang="ts">
	import DynamicForm, { type FormConfig } from '$components/ui/DynamicForm.svelte';

	export interface AboutBlockData {
		title: string;
		description?: string;
		heading?: string;
		subheading?: string;
	}

	interface Props {
		data: AboutBlockData;
		onchange?: (data: AboutBlockData) => void;
	}

	let {
		data = $bindable({ title: '', description: '', heading: '', subheading: '' }),
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
			}
		]
	};

	const handleChange = (updatedData: Record<string, any>) => {
		data = updatedData as AboutBlockData;
		onchange?.(data);
	};
</script>

<DynamicForm asForm={false} config={formConfig} formData={data} onchange={handleChange} />
