<!-- frontend/src/lib/components/blocks/forms/RichtextBlockForm.svelte -->
<script lang="ts">
import DynamicForm, {
	type FormConfig,
} from "$components/ui/DynamicForm.svelte";

interface RichtextBlockData {
	content: string;
}

interface Props {
	data: RichtextBlockData;
	onchange?: (data: RichtextBlockData) => void;
}

let { data = $bindable({ content: "" }), onchange }: Props = $props();

const formConfig: FormConfig = {
	fields: [
		{
			name: "content",
			type: "textarea",
			label: "Content",
			placeholder: "Enter your content here...",
			required: true,
			rows: 8,
			colSpan: 12,
			helperText: "You can use Markdown formatting (future: WYSIWYG editor)",
		},
	],
	hideSubmitButton: true,
};

const handleChange = (updatedData: Record<string, any>) => {
	data = updatedData as RichtextBlockData;
	onchange?.(data);
};
</script>

<DynamicForm 
	config={formConfig} 
	formData={data} 
	onchange={handleChange}
/>
