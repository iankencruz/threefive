<!-- frontend/src/lib/components/blocks/forms/HeroBlockForm.svelte -->
<script lang="ts">
import DynamicForm, {
	type FormConfig,
} from "$components/ui/DynamicForm.svelte";

interface HeroBlockData {
	title: string;
	subtitle?: string;
	cta_text?: string;
	cta_url?: string;
}

interface Props {
	data: HeroBlockData;
	onchange?: (data: HeroBlockData) => void;
}

let {
	data = $bindable({ title: "", subtitle: "", cta_text: "", cta_url: "" }),
	onchange,
}: Props = $props();

const formConfig: FormConfig = {
	fields: [
		{
			name: "title",
			type: "text",
			label: "Title",
			placeholder: "Enter hero title",
			required: true,
			colSpan: 12,
		},
		{
			name: "subtitle",
			type: "text",
			label: "Subtitle",
			placeholder: "Enter subtitle (optional)",
			colSpan: 12,
		},
		{
			name: "media",
			type: "media",
			label: "Image",
			colSpan: 12,
		},
		{
			name: "cta_text",
			type: "text",
			label: "CTA Text",
			placeholder: "Button text (optional)",
			colSpan: 6,
		},
		{
			name: "cta_url",
			type: "url",
			label: "CTA URL",
			placeholder: "Button link (optional)",
			colSpan: 6,
		},
	],
	hideSubmitButton: true,
};

const handleChange = (updatedData: Record<string, any>) => {
	data = updatedData as HeroBlockData;
	onchange?.(data);
};
</script>

<DynamicForm 
	config={formConfig} 
	formData={data} 
	onchange={handleChange}
/>
