<!-- frontend/src/lib/components/ui/form/fields/CheckboxField.svelte -->
<script lang="ts">
	import type { BaseFieldProps } from '../types';

	interface Props extends BaseFieldProps {
		value: boolean;
	}

	let { field, value = $bindable(false), error, disabled = false, onchange }: Props = $props();

	function handleChange(e: Event) {
		const target = e.target as HTMLInputElement;
		onchange(target.checked);
	}
</script>

<div class="my-2 flex w-full items-start gap-3">
	<div class="flex-1">
		{#if field.label}
			<label for={field.name} class="text-sm font-medium">
				{field.label}
				{#if field.required}
					<span class="text-red-500">*</span>
				{/if}
			</label>
		{/if}
		{#if field.helperText}
			<p class="mt-1 text-sm text-gray-500">{field.helperText}</p>
		{/if}
		{#if error}
			<p class="mt-1 text-sm text-red-600">{error}</p>
		{/if}
	</div>
	<div
		class="group relative inline-flex w-11 shrink-0 rounded-full bg-gray-200 p-0.5 inset-ring inset-ring-gray-900/5 outline-offset-2 outline-primary transition-colors duration-200 ease-in-out has-checked:bg-primary has-focus-visible:outline-2"
	>
		<span
			class="relative size-5 rounded-full bg-white shadow-xs ring-1 ring-gray-900/5 transition-transform duration-200 ease-in-out group-has-checked:translate-x-5"
		>
			<span
				aria-hidden="true"
				class="absolute inset-0 flex size-full items-center justify-center opacity-100 transition-opacity duration-200 ease-in group-has-checked:opacity-0 group-has-checked:duration-100 group-has-checked:ease-out"
			>
				<svg viewBox="0 0 12 12" fill="none" class="size-3 text-gray-400">
					<path
						d="M4 8l2-2m0 0l2-2M6 6L4 4m2 2l2 2"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					/>
				</svg>
			</span>
			<span
				aria-hidden="true"
				class="absolute inset-0 flex size-full items-center justify-center opacity-0 transition-opacity duration-100 ease-out group-has-checked:opacity-100 group-has-checked:duration-200 group-has-checked:ease-in"
			>
				<svg viewBox="0 0 12 12" fill="currentColor" class="size-3 text-primary">
					<path
						d="M3.707 5.293a1 1 0 00-1.414 1.414l1.414-1.414zM5 8l-.707.707a1 1 0 001.414 0L5 8zm4.707-3.293a1 1 0 00-1.414-1.414l1.414 1.414zm-7.414 2l2 2 1.414-1.414-2-2-1.414 1.414zm3.414 2l4-4-1.414-1.414-4 4 1.414 1.414z"
					/>
				</svg>
			</span>
		</span>
		<input
			type="checkbox"
			id={field.name}
			name={field.name}
			bind:checked={value}
			{disabled}
			aria-label={field.label}
			class="absolute inset-0 size-full appearance-none rounded-lg focus:outline-hidden"
			onchange={handleChange}
		/>
	</div>
</div>
