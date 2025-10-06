<script lang="ts">
import Input from "$components/ui/Input.svelte";
import Button from "$components/ui/Button.svelte";
import type { Snippet } from "svelte";

export interface FormField {
	name: string;
	type?:
		| "text"
		| "email"
		| "password"
		| "number"
		| "tel"
		| "url"
		| "date"
		| "textarea";
	label?: string;
	placeholder?: string;
	required?: boolean;
	colSpan?: number; // Grid column span (1-12)
	class?: string; // Additional classes for field wrapper
	inputClass?: string; // Additional classes for input element
	helperText?: string;
	rows?: number; // For textarea
}

export interface FormConfig {
	fields: FormField[];
	submitText?: string;
	submitVariant?: "primary" | "secondary" | "outline" | "danger";
	submitFullWidth?: boolean;
	showSubmit?: boolean;
}

interface Props {
	config: FormConfig;
	formData?: Record<string, any>;
	errors?: Record<string, string>;
	loading?: boolean;
	class?: string;
	onchange?: (data: Record<string, any>) => void;
	onsubmit: (data: Record<string, any>) => void;
	children?: Snippet; // For custom form footer/actions
}

let {
	config,
	formData = $bindable({}),
	errors = {},
	loading = false,
	class: className = "",
	onchange, // Add this
	onsubmit,
	children,
}: Props = $props();

// Initialize SYNCHRONOUSLY before render
config.fields.forEach((field) => {
	if (!(field.name in formData)) {
		formData[field.name] = "";
	}
});

$effect(() => {
	onchange?.(formData);
});

function handleSubmit(e: SubmitEvent) {
	e.preventDefault();
	onsubmit(formData);
}

function getColSpanClass(colSpan?: number): string {
	const spanMap: Record<number, string> = {
		1: "col-span-1",
		2: "col-span-2",
		3: "col-span-3",
		4: "col-span-4",
		5: "col-span-5",
		6: "col-span-6",
		7: "col-span-7",
		8: "col-span-8",
		9: "col-span-9",
		10: "col-span-10",
		11: "col-span-11",
		12: "col-span-12",
	};

	return spanMap[colSpan || 12] || "col-span-12";
}
</script>

<form onsubmit={handleSubmit} class={className}>
	<div class="grid grid-cols-12 gap-4">
		{#each config.fields as field}
			<div class="{getColSpanClass(field.colSpan)} {field.class || ''}">
				<Input
					type={field.type || 'text'}
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
			</div>
		{/each}
	</div>

	{#if children}
		<div class="mt-6">
			{@render children()}
		</div>
	{:else if config.showSubmit !== false}
		<div class="mt-6">
			<Button
				type="submit"
				variant={config.submitVariant || 'primary'}
				fullWidth={config.submitFullWidth}
				{loading}
			>
				{config.submitText || 'Submit'}
			</Button>
		</div>
	{/if}
</form>
