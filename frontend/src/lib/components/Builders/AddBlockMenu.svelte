<script lang="ts">
	import { Plus } from '@lucide/svelte';
	import { onDestroy } from 'svelte';

	const blockTypes = [
		{ label: 'Heading', value: 'heading' },
		{ label: 'Image', value: 'image' },
		{ label: 'Rich Text', value: 'richtext' }
	];

	let { onselect }: { onselect: (type: string) => void } = $props();

	let open = $state(false);

	$effect(() => {
		function handleClick(event: MouseEvent) {
			const menu = document.getElementById('block-menu-dropdown');
			if (open && menu && !menu.contains(event.target as Node)) {
				open = false;
			}
		}

		// Use capture phase to ensure this runs before internal clicks
		window.addEventListener('click', handleClick, true);

		onDestroy(() => {
			window.removeEventListener('click', handleClick, true);
		});
	});
	function handleSelect(value: string) {
		onselect(value);
		open = false;
	}

	function toggleMenu() {
		open = !open;
	}
</script>

<!-- Divider + Button Wrapper -->
<div class="relative my-1 w-full">
	<!-- Horizontal line -->
	<div class="absolute inset-0 flex items-center" aria-hidden="true">
		<div class="w-full border-t border-gray-300"></div>
	</div>

	<!-- Centered Button -->
	<div class="relative flex justify-center">
		<button
			onclick={toggleMenu}
			type="button"
			class="inline-flex cursor-pointer items-center gap-x-1.5 rounded-full bg-white p-1 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-gray-300 hover:bg-gray-100 hover:text-gray-700 focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 focus:outline-none"
		>
			<Plus size={16} class="text-gray-400" />
		</button>

		<!-- Dropdown Menu -->
		{#if open}
			<div
				id="block-menu-dropdown"
				class="ring-opacity-5 absolute bottom-full z-20 mb-2 w-44 origin-bottom-left rounded-md bg-white shadow-lg ring-1 ring-black"
			>
				<ul class="py-1 text-sm text-gray-700">
					{#each blockTypes as type}
						<li>
							<button
								type="button"
								class="w-full px-4 py-2 text-left hover:bg-gray-100"
								onclick={() => handleSelect(type.value)}
							>
								{type.label}
							</button>
						</li>
					{/each}
				</ul>
			</div>
		{/if}
	</div>
</div>
