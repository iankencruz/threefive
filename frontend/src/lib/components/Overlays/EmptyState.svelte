<script lang="ts">
	import { FolderPlus } from '@lucide/svelte';
	import type { Snippet } from 'svelte';

	let {
		action,
		title = 'No Project',
		text = 'Get started by creating a new project.',
		icon
	} = $props<{
		action?: () => void;
		title?: string;
		text?: string;
		icon?: Snippet<[any]>;
	}>();
</script>

<div>
	<button
		onclick={action}
		type="button"
		class=" flex w-full flex-col items-center rounded-lg border-2 border-dashed border-gray-300 p-12 text-center hover:border-gray-400 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:outline-hidden"
	>
		{#if icon}
			<div>{@render icon()}</div>
		{:else}
			<FolderPlus />
		{/if}
		<h3 class="mt-2 text-sm font-semibold text-gray-900">{title}</h3>
		<p class="mt-1 text-sm text-gray-500">{text}</p>
	</button>
</div>
