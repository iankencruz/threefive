<script lang="ts">
	// 1. Import the SvelteKit page store to get the current URL
	import { page } from '$app/state';
	import type { Link } from '$types/pages';
	import { type Project } from '$types/projects';
	import { type Blog } from '$types/blogs';
	import { LayoutPanelLeft } from 'lucide-svelte';

	let { variant } = $props<'standard' | 'ghost'>();

	let NavLinks: Link[] = [
		{ id: 1, title: 'Home', href: '/' },
		{ id: 2, title: 'About Us', href: '/about' },
		{ id: 3, title: 'Projects', href: '/projects' },
		{ id: 4, title: 'Blogs', href: '/blogs' }
		// We removed Contact Us here
	];

	const ContactLink: Link = { id: 3, title: 'Contact', href: '/contact' };

	// State Declarations
	let navbarOpen = $state(false);

	// Helper to check if a link is active based on the current SvelteKit URL
	// Checks for exact match OR if the current URL starts with the link's href (useful for path matching, e.g., /projects)
	function checkIsActive(href: string): boolean {
		const currentPath = page.url.pathname;

		if (href === '/') {
			// Home link is only active on the exact root path
			return currentPath === '/';
		}

		// For non-root links, check if the current path starts with the link's href.
		// E.g., if href is /about, and currentPath is /about/team, it matches.
		return currentPath.startsWith(href);
	}

	// Function to close the main navbar when a link is clicked
	function closeNavbar() {
		if (navbarOpen) {
			navbarOpen = false;
		}
	}
</script>

<nav
	class="z-20 w-full transition-all duration-500"
	class:absolute={variant === 'ghost'}
	class:bg-background={variant !== 'ghost'}
>
	<div class="max-w-8xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="flex w-full items-center justify-between py-4">
			{#if variant !== 'ghost'}
				<a href="/" class="w-full">
					<div class="flex flex-row gap-2 leading-none font-bold text-primary sm:hidden">
						<span>三</span>
						<span>五</span>
					</div>
					<div class="hidden w-full text-lg font-bold text-foreground sm:block">
						Threefive Project
					</div>
				</a>
			{/if}

			<div class="hidden w-full pt-4 lg:flex lg:pl-11" id="navbar-desktop">
				<ul class="flex gap-2 lg:mt-0 lg:ml-auto lg:flex-row lg:items-center lg:justify-center">
					{#each NavLinks as link (link.id)}
						<li>
							<a
								href={link.href}
								class="nav-link block text-base font-medium transition-all duration-500 lg:mx-3"
								class:text-primary={checkIsActive(link.href)}
								class:text-white={!checkIsActive(link.href)}
								class:hover:text-gray-400={!checkIsActive(link.href)}
							>
								{link.title}
							</a>
						</li>
					{/each}

					<li>
						<a
							href={ContactLink.href}
							class="nav-link block text-base font-medium transition-all duration-500 lg:mx-3"
							class:text-primary={checkIsActive(ContactLink.href)}
							class:text-white={!checkIsActive(ContactLink.href)}
							class:hover:text-gray-400={!checkIsActive(ContactLink.href)}
						>
							{ContactLink.title}
						</a>
					</li>
				</ul>
			</div>

			<div class="flex items-center justify-end gap-5 lg:hidden">
				<button
					onclick={() => (navbarOpen = !navbarOpen)}
					type="button"
					class="inline-flex items-center rounded-lg p-2 text-sm text-white hover:bg-gray-100 hover:text-gray-900 focus:ring-2 focus:ring-gray-200 focus:outline-none"
					aria-controls="mobile-drawer"
					aria-expanded={navbarOpen}
				>
					<span class="sr-only">Open main menu</span>
					{#if navbarOpen}
						<svg
							class="h-6 w-6"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							xmlns="http://www.w3.org/2000/svg"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M6 18L18 6M6 6l12 12"
							></path>
						</svg>
					{:else}
						<svg
							class="h-6 w-6"
							aria-hidden="true"
							fill="currentColor"
							viewBox="0 0 20 20"
							xmlns="http://www.w3.org/2000/svg"
						>
							<path
								fill-rule="evenodd"
								d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
								clip-rule="evenodd"
							></path>
						</svg>
					{/if}
				</button>
			</div>
		</div>
	</div>
</nav>

<div
	id="mobile-drawer"
	class="fixed inset-y-0 left-0 z-40 h-full w-full bg-white shadow-xl transition-transform duration-300 sm:w-84 lg:hidden"
	class:translate-x-0={navbarOpen}
	class:translate-x-[-100%]={!navbarOpen}
>
	<div class="flex items-center justify-between border-b border-gray-100 px-4 py-4">
		<a href="/" class="flex w-max items-center">
			<span class="text-lg font-bold text-gray-900">Threefive Project</span>
		</a>

		<button
			onclick={() => (navbarOpen = false)}
			type="button"
			class="rounded-md p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 focus:ring-2 focus:ring-indigo-500 focus:outline-none"
			aria-label="Close menu"
		>
			<svg
				class="h-6 w-6"
				fill="none"
				stroke="currentColor"
				viewBox="0 0 24 24"
				xmlns="http://www.w3.org/2000/svg"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M6 18L18 6M6 6l12 12"
				></path>
			</svg>
		</button>
	</div>

	<div class="overflow-y-auto py-2">
		<ul class="space-y-0">
			{#each NavLinks as link (link.id)}
				<li>
					<a
						onclick={closeNavbar}
						href={link.href}
						class="block border-l-4 py-4 pr-4 pl-3 text-base font-medium transition duration-150 ease-in-out"
						class:border-indigo-700={checkIsActive(link.href)}
						class:bg-indigo-50={checkIsActive(link.href)}
						class:text-indigo-700={checkIsActive(link.href)}
						class:border-transparent={!checkIsActive(link.href)}
						class:text-gray-700={!checkIsActive(link.href)}
						class:hover:bg-gray-50={!checkIsActive(link.href)}
						class:hover:text-gray-900={!checkIsActive(link.href)}
					>
						{link.title}
					</a>
				</li>
			{/each}
		</ul>

		<div class="mt-4 border-t border-gray-100 pt-2">
			<h6 class="mb-1 px-3 text-sm font-medium text-gray-500">Projects</h6>
			<ul class="space-y-0">
				{#each Projects as project}
					<li>
						<a
							onclick={closeNavbar}
							href={`/projects/${project.slug}`}
							class="block border-l-4 py-4 pr-4 pl-3 text-base font-medium transition duration-150 ease-in-out"
							class:border-indigo-700={checkIsActive(project.slug)}
							class:bg-indigo-50={checkIsActive(project.slug)}
							class:text-indigo-700={checkIsActive(project.slug)}
							class:border-transparent={!checkIsActive(project.slug)}
							class:text-gray-700={!checkIsActive(project.slug)}
							class:hover:bg-gray-50={!checkIsActive(project.slug)}
							class:hover:text-gray-900={!checkIsActive(project.slug)}
						>
							{project.title}
						</a>
					</li>
				{/each}
				<li>
					<a
						onclick={closeNavbar}
						href="/projects"
						class="block border-l-4 border-transparent py-4 pr-4 pl-3 text-base font-medium text-gray-700 hover:bg-gray-50 hover:text-gray-900"
					>
						More...
					</a>
				</li>
			</ul>
		</div>
		<div class="mt-4 border-t border-gray-100 pt-2">
			<h6 class="mb-1 px-3 text-sm font-medium text-gray-500">Socials</h6>
		</div>
	</div>
</div>

<div
	class="fixed inset-0 z-30 bg-black/50 transition-opacity duration-300 lg:hidden"
	class:opacity-100={navbarOpen}
	class:opacity-0={!navbarOpen}
	class:pointer-events-none={!navbarOpen}
></div>
