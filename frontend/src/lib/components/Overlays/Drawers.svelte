<!-- src/lib/components/Drawer.svelte -->
<script lang="ts">
	let { open = false, title = '', description = '', onclose, onsubmit, children } = $props();
</script>

{#if open}
	<div class="fixed inset-0 z-50" role="dialog" aria-modal="true">
		<!-- Backdrop -->
		<div
			class="fixed inset-0 bg-black/30 backdrop-blur-sm"
			role="dialog"
			tabindex="0"
			aria-modal="true"
			aria-labelledby="dialog_drawer"
			aria-describedby="dialog_desc"
			onkeydown={(e) => {
				if (e.key === 'Escape') onclose?.();
			}}
			onclick={() => onclose?.()}
		></div>

		<!-- Drawer panel -->
		<div class="fixed inset-y-0 right-0 flex max-w-full pl-10 sm:pl-16">
			<div class="pointer-events-auto w-screen max-w-2xl">
				<form
					class="flex h-full flex-col bg-white shadow-xl"
					onsubmit={(e) => (e.preventDefault(), onsubmit?.(e))}
				>
					<!-- Header -->
					<div class="flex items-start justify-between bg-gray-50 px-4 py-6 sm:px-6">
						<div>
							<h2 class="text-base font-semibold text-gray-900">{title}</h2>
							{#if description}
								<p class="mt-1 text-sm text-gray-500">{description}</p>
							{/if}
						</div>
						<button
							aria-label="Close"
							type="button"
							onclick={() => onclose?.()}
							class="text-gray-400 hover:text-gray-600"
						>
							<svg
								class="h-6 w-6"
								fill="none"
								stroke="currentColor"
								stroke-width="1.5"
								viewBox="0 0 24 24"
							>
								<path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
							</svg>
						</button>
					</div>

					<!-- Content slot -->
					<div class="flex-1 space-y-6 overflow-y-auto px-4 py-6 sm:px-6">
						{@render children()}
					</div>

					<!-- Actions -->
					<div class="flex justify-end gap-2 border-t px-4 py-4 sm:px-6">
						<button
							type="button"
							onclick={() => onclose?.()}
							class="rounded-md bg-white px-4 py-2 text-sm font-semibold text-gray-900 ring-1 ring-gray-300 hover:bg-gray-50"
							>Cancel</button
						>
						<button
							type="submit"
							class="rounded-md bg-indigo-600 px-4 py-2 text-sm font-semibold text-white hover:bg-indigo-500"
							>Save</button
						>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}
