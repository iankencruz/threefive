<!-- frontend/src/routes/+error.svelte -->
<script lang="ts">
	import { page } from '$app/state';
</script>

<svelte:head>
	<title>{page.status} - {page.error?.message || 'Error'}</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-gray-50 px-4">
	<div class="w-full max-w-md text-center">
		<div class="mb-8">
			{#if page.status === 404}
				<svg
					class="mx-auto mb-4 h-24 w-24 text-gray-400"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
					/>
				</svg>
			{:else}
				<svg
					class="mx-auto mb-4 h-24 w-24 text-red-400"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
					/>
				</svg>
			{/if}
		</div>

		<h1 class="mb-4 text-6xl font-bold text-gray-900">
			{page.status}
		</h1>

		{console.log(page.status)}
		<h2 class="mb-4 text-2xl font-semibold text-gray-700">
			{#if page.status === 404}
				Page Not Found
			{:else if page.status === 401}
				Not Authorized
			{:else if page.status === 500}
				Internal Server Error
			{:else}
				Something Went Wrong
			{/if}
		</h2>

		<p class="mb-8 text-gray-600">
			{page.error?.message || 'An unexpected error occurred'}
		</p>

		<div class="flex flex-col justify-center gap-4 sm:flex-row">
			<a
				href="/"
				class="rounded-lg bg-primary px-6 py-3 font-medium text-white transition-colors hover:bg-primary"
			>
				Go Home
			</a>
			<button
				onclick={() => window.history.back()}
				class="rounded-lg bg-gray-200 px-6 py-3 font-medium text-gray-700 transition-colors hover:bg-gray-300"
			>
				Go Back
			</button>
		</div>
	</div>
</div>
