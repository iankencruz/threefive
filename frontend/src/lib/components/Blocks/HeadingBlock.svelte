<script lang="ts">
	import type { Block } from '$lib/types';

	let {
		block,
		onupdate
	}: {
		block: Block;
		onupdate: (updated: Block) => void;
	} = $props();

	let localProps: Record<string, string> = $state({
		title: '',
		description: ''
	});

	// Safely hydrate props, even if props were nested like { props: { title, description } }
	$effect(() => {
		const incoming = block.props?.props ?? block.props ?? {};
		localProps.title = incoming.title ?? '';
		localProps.description = incoming.description ?? '';
	});

	function handleChange(e: Event) {
		const target = e.target as HTMLInputElement | HTMLTextAreaElement;

		localProps = {
			...localProps,
			[target.name]: target.value
		};

		// This avoids re-nesting props again
		onupdate({
			...block,
			props: { ...localProps }
		});
	}
</script>

<div class="space-y-2">
	<label for="title" class="block text-xs font-light text-gray-700">Title</label>
	<input
		name="title"
		type="text"
		placeholder="Heading title"
		class="w-full rounded border px-3 py-2"
		value={localProps.title}
		oninput={handleChange}
	/>

	<label for="description" class="block text-xs font-light text-gray-700">Description</label>
	<textarea
		name="description"
		rows="2"
		placeholder="Optional subtext"
		class="w-full rounded border px-3 py-2"
		oninput={handleChange}>{localProps.description}</textarea
	>
</div>
