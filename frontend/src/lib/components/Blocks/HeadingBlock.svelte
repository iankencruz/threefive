<script lang="ts">
	export interface HeadingBlockProps {
		title: string;
		description: string;
	}

	let {
		props,
		onupdate
	}: {
		props: HeadingBlockProps;
		onupdate: (props: HeadingBlockProps) => void;
	} = $props();

	function handleChange(e: Event) {
		const target = e.target as HTMLInputElement | HTMLTextAreaElement;
		onupdate({ ...props, [target.name]: target.value });
	}
</script>

<div class="space-y-2">
	<input
		name="title"
		type="text"
		placeholder="Heading title"
		class="w-full rounded border px-3 py-2"
		value={props.title}
		oninput={handleChange}
	/>

	<textarea
		name="description"
		rows="2"
		placeholder="Optional subtext"
		class="w-full rounded border px-3 py-2"
		oninput={handleChange}>{props.description}</textarea
	>
</div>
