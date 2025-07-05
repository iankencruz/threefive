<script lang="ts">
	let { onselect }: { onselect: (type: string) => void } = $props();

	const blockTypes: { label: string; value: string }[] = [
		{ label: 'Heading', value: 'heading' },
		{ label: 'Image', value: 'image' },
		{ label: 'Gallery', value: 'gallery' },
		{ label: 'Rich Text', value: 'richtext' }
	];

	function handleChange(e: Event) {
		const target = e.target as HTMLSelectElement;
		const selected = target.value;
		if (selected) {
			onselect(selected);
			target.selectedIndex = 0; // Reset
		}
	}
</script>

<select
	onchange={handleChange}
	class="w-full max-w-32 rounded border bg-white px-2 py-1 text-sm shadow-sm"
>
	<option disabled selected>+ Add Block</option>
	{#each blockTypes as block}
		<option value={block.value}>{block.label}</option>
	{/each}
</select>
