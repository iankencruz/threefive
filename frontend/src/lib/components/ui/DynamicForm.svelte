<script lang="ts">
import Button from "./Button.svelte";
import Input from "./Input.svelte";
import MediaPicker from "./MediaPicker.svelte";

interface FormField {
	name: string;
	label: string;
	type?:
		| "text"
		| "email"
		| "password"
		| "tel"
		| "url"
		| "date"
		| "number"
		| "textarea"
		| "media";
	placeholder?: string;
	required?: boolean;
	value?: string | number;
	colSpan?: number;
	class?: string;
	helperText?: string;
	rows?: number;
	inputClass?: string;
}

interface FormConfig {
	fields: FormField[];
	submitText?: string;
	showSubmit?: boolean;
	columns?: number;
}

interface Props {
	config: FormConfig;
	formData?: Record<string, any>;
	onSubmit?: (data: Record<string, any>) => void | Promise<void>;
	onchange?: (data: Record<string, any>) => void; // ✅ Add this
	errors?: Record<string, string>; // ✅ Add this
	children?: import("svelte").Snippet;
}

let {
	config,
	formData: initialFormData, // ✅ Rename to avoid confusion
	onSubmit,
	onchange,
	errors = {}, // ✅ Accept errors from parent
	children,
}: Props = $props();

// ✅ Watch for changes in formData and call onchange
$effect(() => {
	if (onchange) {
		onchange(formData);
	}
});

// Initialize form data with default values immediately
let formData = $state<Record<string, any>>(
	config.fields.reduce(
		(acc, field) => {
			acc[field.name] = field.value ?? (field.type === "media" ? "" : "");
			return acc;
		},
		{} as Record<string, any>,
	),
);

function getColSpanClass(colSpan?: number): string {
	if (!colSpan) return "col-span-1";
	return `col-span-${colSpan}`;
}

async function handleSubmit(e: SubmitEvent) {
	e.preventDefault();
	if (onSubmit) {
		await onSubmit(formData);
	}
}

function handleMediaChange(fieldName: string, mediaId: string | null) {
	formData[fieldName] = mediaId || "";
}

// Helper to check if field type is media
function isMediaField(type?: string): boolean {
	return type === "media";
}

// Helper to get valid input type (excludes 'media')
function getInputType(
	type?: string,
):
	| "text"
	| "email"
	| "password"
	| "tel"
	| "url"
	| "date"
	| "number"
	| "textarea"
	| undefined {
	if (type === "media") return undefined;
	return type as
		| "text"
		| "email"
		| "password"
		| "tel"
		| "url"
		| "date"
		| "number"
		| "textarea"
		| undefined;
}
</script>

<form onsubmit={handleSubmit} class="space-y-6">
	<div class="grid grid-cols-{config.columns || 1} gap-4">
		{#each config.fields as field}
			<div class="{getColSpanClass(field.colSpan)} {field.class || ''}">
				{#if isMediaField(field.type)}
            <MediaPicker
						bind:value={formData[field.name]}
						label={field.label}
						required={field.required}
						error={errors[field.name]}
						onchange={(mediaId) => handleMediaChange(field.name, mediaId)}
            />
				{:else}
					<Input
						type={getInputType(field.type) || 'text'}
						name={field.name}
						label={field.label}
						placeholder={field.placeholder}
						required={field.required}
						bind:value={formData[field.name]}
						error={errors[field.name]}
						helperText={field.helperText}
						rows={field.rows}
						inputClass={field.inputClass}
					/>
				{/if}
			</div>
		{/each}
	</div>

	{#if children}
		{@render children()}
	{:else if config.showSubmit !== false}
		<Button type="submit" class="w-full">
			{config.submitText || 'Submit'}
		</Button>
	{/if}
</form>
