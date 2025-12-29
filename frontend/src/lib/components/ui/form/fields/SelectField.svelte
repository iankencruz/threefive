<!-- frontend/src/lib/components/ui/form/fields/SelectField.svelte -->
<script lang="ts">
	import type { BaseFieldProps } from '../types';

	interface Props extends BaseFieldProps {}

	let { field, value = $bindable(''), error, disabled = false, onchange }: Props = $props();

	function handleChange(e: Event) {
		const target = e.target as HTMLSelectElement;
		onchange(target.value);
	}
</script>

<div class={field.class || ''}>
	{#if field.label}
		<label for={field.name} class="mb-1.5 block text-sm font-medium text-foreground/80">
			{field.label}
			{#if field.required}
				<span class="text-red-500">*</span>
			{/if}
		</label>
	{/if}

	<select
		id={field.name}
		name={field.name}
		bind:value
		required={field.required}
		{disabled}
		class="form-input {error ? 'border-danger focus:border-danger focus:ring-danger/50' : ''}"
		onchange={handleChange}
	>
		{#if field.placeholder}
			<option value="" disabled>{field.placeholder}</option>
		{/if}
		{#each field.options || [] as option}
			<option value={option.value}>{option.label}</option>
		{/each}
	</select>

	{#if error}
		<p class="mt-1 text-sm text-danger">{error}</p>
	{:else if field.helperText}
		<p class="mt-1 text-sm text-foreground/60">{field.helperText}</p>
	{/if}
</div>
