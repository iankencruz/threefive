<script lang="ts">
interface Props {
	type?:
		| "text"
		| "email"
		| "password"
		| "number"
		| "tel"
		| "url"
		| "date"
		| "textarea";
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
	["form-input", error && "form-input-error", inputClass]
		.filter(Boolean)
		.join(" "),
);
</script>

<div class={className}>
	{#if label}
		<label for={name} class="form-label">
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
