<script lang="ts">
	// 1. Import the SvelteKit page store to get the current URL
	import { page } from '$app/state';
	import type { Link } from '$types/pages';
	import { Projects } from '$types/projects';
	import { Blogs } from '$types/blogs';

	let { variant } = $props<'standard' | 'ghost'>();

	// 1. Separate the main navigation links (excluding 'Contact Us')
	let PrimaryNavLinks: Link[] = [
		{ id: 1, title: 'Home', href: '/' },
		{ id: 2, title: 'About Us', href: '/about' }
		// We removed Contact Us here
	];

	const ContactLink: Link = { id: 3, title: 'Contact', href: '/contact' };

	let NavLinks: Link[] = [...PrimaryNavLinks, ContactLink];

	// State Declarations
	let navbarOpen = $state(false);
	let featuresOpen = $state(false);
	let blogsOpen = $state(false);

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
			<a href="/" class="flex w-max items-center"> Threefive Project </a>

			<div class="hidden w-full lg:flex lg:pl-11" id="navbar-desktop">
				<ul class="flex gap-6 lg:mt-0 lg:ml-auto lg:flex-row lg:items-center lg:justify-center">
					{#each PrimaryNavLinks as link (link.id)}
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

					{#if Projects.length > 0}
						<li class="relative">
							<button
								popovertarget="desktop-menu-features"
								onclick={() => (featuresOpen = !featuresOpen)}
								class="dropdown-toggle mr-auto flex items-center justify-between text-center text-base font-medium text-white transition-all duration-500 hover:text-gray-400 lg:m-0 lg:mx-3 lg:mb-0 lg:text-left"
								aria-expanded={featuresOpen}
							>
								Projects
								<svg
									class="ml-1.5 h-2 w-3"
									width="10"
									height="6"
									viewBox="0 0 10 6"
									fill="none"
									xmlns="http://www.w3.org/2000/svg"
								>
									<path
										d="M1 1L3.58579 3.58579C4.25245 4.25245 4.58579 4.58579 5 4.58579C5.41421 4.58579 5.74755 4.25245 6.41421 3.58579L9 1"
										stroke="currentColor"
										stroke-width="1.6"
										stroke-linecap="round"
										stroke-linejoin="round"
									></path>
								</svg>
							</button>

							<el-popover
								id="desktop-menu-features"
								popover="auto"
								class:hidden={!featuresOpen}
								onclose={() => (featuresOpen = false)}
								class="absolute top-16 w-full overflow-visible bg-white transition transition-discrete backdrop:bg-transparent open:block data-closed:-translate-y-1 data-closed:opacity-0 data-enter:duration-200 data-enter:ease-out data-leave:duration-150 data-leave:ease-in"
							>
								<!-- Presentational element used to render the bottom shadow, if we put the shadow on the actual panel it pokes out the top, so we use this shorter element to hide the top of the shadow -->
								<div
									aria-hidden="true"
									class="absolute inset-0 top-1/2 bg-white shadow-lg ring-1 ring-gray-900/5"
								></div>
								<div class="relative bg-white">
									<div
										class="mx-auto grid max-w-7xl grid-cols-4 gap-x-4 px-6 py-10 lg:px-8 xl:gap-x-8"
									>
										{#each Projects.slice(0, 4) as project}
											<div class="group relative rounded-lg p-6 text-sm/6 hover:bg-gray-50">
												<a
													href={`/projects/${project.slug}`}
													class=" block font-semibold text-gray-900"
												>
													{project.title}
													<span class="absolute inset-0"></span>
												</a>
												<p class="mt-1 text-gray-600">
													{project.description || 'Explore our project'}
												</p>
											</div>
										{/each}
									</div>
									<div class="bg-gray-50">
										<div class="mx-auto max-w-7xl px-6 lg:px-8">
											<div class="grid divide-x divide-gray-900/5 border-x border-gray-900/5">
												<a
													href={`/projects`}
													class="flex items-center justify-center gap-x-2.5 p-3 text-sm/6 font-semibold text-gray-900 hover:bg-gray-100"
												>
													<svg
														viewBox="0 0 20 20"
														fill="currentColor"
														data-slot="icon"
														aria-hidden="true"
														class="size-5 flex-none text-gray-400"
													>
														<path
															d="M2.5 3A1.5 1.5 0 0 0 1 4.5v4A1.5 1.5 0 0 0 2.5 10h6A1.5 1.5 0 0 0 10 8.5v-4A1.5 1.5 0 0 0 8.5 3h-6Zm11 2A1.5 1.5 0 0 0 12 6.5v7a1.5 1.5 0 0 0 1.5 1.5h4a1.5 1.5 0 0 0 1.5-1.5v-7A1.5 1.5 0 0 0 17.5 5h-4Zm-10 7A1.5 1.5 0 0 0 2 13.5v2A1.5 1.5 0 0 0 3.5 17h6a1.5 1.5 0 0 0 1.5-1.5v-2A1.5 1.5 0 0 0 9.5 12h-6Z"
															clip-rule="evenodd"
															fill-rule="evenodd"
														/>
													</svg>
													View all projects
												</a>
											</div>
										</div>
									</div>
								</div>
							</el-popover>
						</li>
					{/if}

					{#if Blogs.length > 0}
						<li class="relative">
							<button
								popovertarget="desktop-menu-blogs"
								onclick={() => (blogsOpen = !blogsOpen)}
								class="dropdown-toggle mr-auto flex items-center justify-between text-center text-base font-medium text-white transition-all duration-500 hover:text-gray-400 lg:m-0 lg:mx-3 lg:mb-0 lg:text-left"
								aria-expanded={blogsOpen}
							>
								Blogs
								<svg
									class="ml-1.5 h-2 w-3"
									width="10"
									height="6"
									viewBox="0 0 10 6"
									fill="none"
									xmlns="http://www.w3.org/2000/svg"
								>
									<path
										d="M1 1L3.58579 3.58579C4.25245 4.25245 4.58579 4.58579 5 4.58579C5.41421 4.58579 5.74755 4.25245 6.41421 3.58579L9 1"
										stroke="currentColor"
										stroke-width="1.6"
										stroke-linecap="round"
										stroke-linejoin="round"
									></path>
								</svg>
							</button>

							<el-popover
								id="desktop-menu-blogs"
								popover="auto"
								class:hidden={!featuresOpen}
								onclose={() => (featuresOpen = false)}
								class="absolute top-16 w-full overflow-visible bg-white transition transition-discrete backdrop:bg-transparent open:block data-closed:-translate-y-1 data-closed:opacity-0 data-enter:duration-200 data-enter:ease-out data-leave:duration-150 data-leave:ease-in"
							>
								<!-- Presentational element used to render the bottom shadow, if we put the shadow on the actual panel it pokes out the top, so we use this shorter element to hide the top of the shadow -->
								<div
									aria-hidden="true"
									class="absolute inset-0 top-1/2 bg-white shadow-lg ring-1 ring-gray-900/5"
								></div>
								<div class="relative bg-white">
									<div
										class="max-w-8xl mx-auto grid grid-cols-4 gap-x-4 px-6 py-10 lg:px-8 xl:gap-x-8"
									>
										{#each Blogs.slice(0, 4) as blog}
											<div class="group relative rounded-lg p-6 text-sm/6 hover:bg-gray-50">
												<a href={`/blogs/${blog.slug}`} class=" block font-semibold text-gray-900">
													{blog.title}
													<span class="absolute inset-0"></span>
												</a>
												<p class="mt-1 line-clamp-2 w-full text-ellipsis text-gray-600">
													{blog.excerpt || 'Read our blog'}
												</p>
											</div>
										{/each}
									</div>
									<div class="bg-gray-50">
										<div class="mx-auto max-w-7xl px-6 lg:px-8">
											<div class="grid divide-x divide-gray-900/5 border-x border-gray-900/5">
												<a
													href={`/projects`}
													class="flex items-center justify-center gap-x-2.5 p-3 text-sm/6 font-semibold text-gray-900 hover:bg-gray-100"
												>
													<svg
														viewBox="0 0 20 20"
														fill="currentColor"
														data-slot="icon"
														aria-hidden="true"
														class="size-5 flex-none text-gray-400"
													>
														<path
															d="M2.5 3A1.5 1.5 0 0 0 1 4.5v4A1.5 1.5 0 0 0 2.5 10h6A1.5 1.5 0 0 0 10 8.5v-4A1.5 1.5 0 0 0 8.5 3h-6Zm11 2A1.5 1.5 0 0 0 12 6.5v7a1.5 1.5 0 0 0 1.5 1.5h4a1.5 1.5 0 0 0 1.5-1.5v-7A1.5 1.5 0 0 0 17.5 5h-4Zm-10 7A1.5 1.5 0 0 0 2 13.5v2A1.5 1.5 0 0 0 3.5 17h6a1.5 1.5 0 0 0 1.5-1.5v-2A1.5 1.5 0 0 0 9.5 12h-6Z"
															clip-rule="evenodd"
															fill-rule="evenodd"
														/>
													</svg>
													View all blogs
												</a>
											</div>
										</div>
									</div>
								</div>
							</el-popover>
						</li>
					{/if}

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
