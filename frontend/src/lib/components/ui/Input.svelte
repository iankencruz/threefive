<script lang="ts">
interface Props {
	type?: "text" | "email" | "password" | "number" | "tel" | "url" | "date" | "textarea";
	name: string;
	label?: string;
	placeholder?: string;
	value?: string | number;
	error?: string;
	helperText?: string;
	required?: boolean;
	disabled?: boolean;
	readonly?: boolean;
	rows?: number; // For textarea
	class?: string;
	inputClass?: string;
	oninput?: (e: Event) => void;
	onblur?: (e: FocusEvent) => void;
	onfocus?: (e: FocusEvent) => void;
}

let {
	type = "text",
	name,
	label,
	placeholder,
	value = $bindable(""),
	error,
	helperText,
	required = false,
	disabled = false,
	readonly = false,
	rows = 4,
	class: className = "",
	inputClass = "",
	oninput,
	onblur,
	onfocus,
}: Props = $props();

const inputClasses = $derived(
	[
		"w-full px-4 py-2.5 rounded-sm border bg-input-bg border-input-border text-foreground placeholder:text-input-placeholder focus:outline-none focus:ring-2 focus:ring-input-focus-ring focus:border-input-focus-border transition-colors duration-200 appearance-none [background-clip:padding-box] disabled:bg-input-disabled-bg disabled:text-input-disabled-text disabled:cursor-not-allowed",
		error && "form-input-error",
		inputClass,
	]
		.filter(Boolean)
		.join(" "),
);
</script>

<div class={className}>
	{#if label}
		<label for={name} class="block text-sm font-medium text-foreground/80 mb-1.5">
			{label}
			{#if required}
				<span class="text-danger">*</span>
			{/if}
		</label>
	{/if}

	{#if type === 'textarea'}
		<textarea
			id={name}
			{name}
			{placeholder}
			{required}
			{disabled}
			{readonly}
			{rows}
			class={inputClasses}
			bind:value
			oninput={oninput}
			onblur={onblur}
			onfocus={onfocus}
		></textarea>
	{:else}
		<input
			id={name}
			{type}
			{name}
			{placeholder}
			{required}
			{disabled}
			{readonly}
			class={inputClasses}
			bind:value
			oninput={oninput}
			onblur={onblur}
			onfocus={onfocus}
		/>
	{/if}

	{#if error}
		<p class="form-error">{error}</p>
	{:else if helperText}
		<p class="form-helper">{helperText}</p>
	{/if}
</div>
