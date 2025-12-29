<!-- frontend/src/lib/components/ui/form/fields/TextareaField.svelte -->
<script lang="ts">
	import Input from '$components/ui/Input.svelte';
	import type { BaseFieldProps } from '../types';

	interface Props extends BaseFieldProps {}

	let { field, value = $bindable(''), error, disabled = false, onchange }: Props = $props();

	function handleInput(e: Event) {
		const target = e.target as HTMLTextAreaElement;
		onchange(target.value);
	}
</script>

<Input
	type="textarea"
	name={field.name}
	label={field.label}
	placeholder={field.placeholder}
	required={field.required}
	{disabled}
	bind:value
	{error}
	helperText={field.helperText}
	rows={field.rows || 4}
	inputClass={field.inputClass}
	oninput={handleInput}
/>
