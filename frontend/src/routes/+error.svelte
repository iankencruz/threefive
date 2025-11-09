<!-- frontend/src/routes/+error.svelte -->
<script lang="ts">
import { page } from "$app/stores";
</script>

<svelte:head>
	<title>{$page.status} - {$page.error?.message || 'Error'}</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 flex items-center justify-center px-4">
	<div class="max-w-md w-full text-center">
		<div class="mb-8">
			{#if $page.status === 404}
				<svg class="w-24 h-24 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>
			{:else}
				<svg class="w-24 h-24 mx-auto text-red-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
				</svg>
			{/if}
		</div>

		<h1 class="text-6xl font-bold text-gray-900 mb-4">
			{$page.status}
		</h1>
		
		<h2 class="text-2xl font-semibold text-gray-700 mb-4">
			{#if $page.status === 404}
				Page Not Found
			{:else if $page.status === 500}
				Internal Server Error
			{:else}
				Something Went Wrong
			{/if}
		</h2>
		
		<p class="text-gray-600 mb-8">
			{$page.error?.message || 'An unexpected error occurred'}
		</p>
		
		<div class="flex flex-col sm:flex-row gap-4 justify-center">
			<a 
				href="/" 
				class="px-6 py-3 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors"
			>
				Go Home
			</a>
			<button 
				onclick={() => window.history.back()}
				class="px-6 py-3 bg-gray-200 text-gray-700 rounded-lg font-medium hover:bg-gray-300 transition-colors"
			>
				Go Back
			</button>
		</div>
	</div>
</div>
