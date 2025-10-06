<!-- frontend/src/lib/components/blocks/forms/HeaderBlockForm.svelte -->
<script lang="ts">
import DynamicForm, {
	type FormConfig,
} from "$components/ui/DynamicForm.svelte";

interface HeaderBlockData {
	heading: string;
	subheading?: string;
	level: string;
}

interface Props {
	data: HeaderBlockData;
	onchange?: (data: HeaderBlockData) => void;
}

let {
	data = $bindable({ heading: "", subheading: "", level: "h2" }),
	onchange,
}: Props = $props();

const formConfig: FormConfig = {
	fields: [
		{
			name: "heading",
			type: "text",
			label: "Heading",
			placeholder: "Enter heading text",
			required: true,
			colSpan: 12,
		},
		{
			name: "subheading",
			type: "text",
			label: "Subheading",
			placeholder: "Enter subheading (optional)",
			colSpan: 12,
		},
		{
			name: "level",
			type: "select",
			label: "Heading Level",
			options: [
				{ value: "h1", label: "H1 - Largest" },
				{ value: "h2", label: "H2 - Large" },
				{ value: "h3", label: "H3 - Medium" },
				{ value: "h4", label: "H4 - Small" },
				{ value: "h5", label: "H5 - Smaller" },
				{ value: "h6", label: "H6 - Smallest" },
			],
			colSpan: 12,
		},
	],
	hideSubmitButton: true,
};

const handleChange = (updatedData: Record<string, any>) => {
	data = updatedData as HeaderBlockData;
	onchange?.(data);
};
</script>

<DynamicForm 
	config={formConfig} 
	formData={data} 
	onchange={handleChange}
/>
