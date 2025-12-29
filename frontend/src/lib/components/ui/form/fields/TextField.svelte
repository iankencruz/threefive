<!-- frontend/src/lib/components/ui/form/fields/TextField.svelte -->
<script lang="ts">
	import Input from '$components/ui/Input.svelte';
	import type { BaseFieldProps } from '../types';

	interface Props extends BaseFieldProps {}

	let { field, value = $bindable(''), error, disabled = false, onchange }: Props = $props();

	function handleInput(e: Event) {
		const target = e.target as HTMLInputElement;
		onchange(target.value);
	}
</script>

<Input
	type={field.type === 'text' ||
	field.type === 'email' ||
	field.type === 'password' ||
	field.type === 'tel' ||
	field.type === 'url' ||
	field.type === 'date' ||
	field.type === 'number'
		? field.type
		: 'text'}
	name={field.name}
	label={field.label}
	placeholder={field.placeholder}
	required={field.required}
	{disabled}
	bind:value
	{error}
	helperText={field.helperText}
	inputClass={field.inputClass}
	oninput={handleInput}
/>
