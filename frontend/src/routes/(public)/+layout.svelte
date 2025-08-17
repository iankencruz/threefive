<script lang="ts">
	import { navigationBar, type NavigationItem } from '$lib/components/Navigation/navigation';
	import '../../app.css';

	let { children, data } = $props();

	console.log('User from layout:', data.user);
</script>

<!-- NavItem Snippet -->
{#snippet NavItem(NavItems: NavigationItem[])}
	<div
		class="mt-5 hidden w-full rounded-lg bg-zinc-700 p-4 backdrop-blur-3xl max-lg:py-4 lg:mt-0 lg:flex lg:bg-transparent lg:pl-10 lg:backdrop-blur-none"
		id="navbar"
	>
		<ul class="mt-4 flex flex-col gap-6 max-lg:mb-4 lg:mt-0 lg:flex-row lg:items-center">
			{#each NavItems as item}
				<li class="px-1.5">
					<a
						href={item.href}
						class="block text-base font-medium text-white transition-all duration-500 hover:text-gray-400 lg:px-1.5 lg:text-left"
						>{item.label}</a
					>
				</li>
			{/each}
		</ul>
	</div>
{/snippet}

<main class="flex h-screen w-full flex-col justify-start bg-black text-white">
	<div>
		<nav class="fixed z-10 w-full bg-transparent py-5 transition-all duration-500">
			<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
				<div class="flex w-full flex-col lg:flex-row">
					<div class="flex shrink-0 justify-between lg:flex-row">
						<a href="/" class="flex w-full grow items-center text-2xl font-bold text-white">
							<div class="">ThreeFive Projects</div>
						</a>
						<button
							data-collapse-toggle="navbar"
							type="button"
							class="ml-3 inline-flex items-center rounded-lg p-2 text-sm text-white hover:bg-gray-100 hover:text-gray-900 focus:ring-2 focus:ring-gray-200 focus:outline-none lg:hidden"
							aria-controls="navbar-default"
							aria-expanded="false"
						>
							<span class="sr-only">Open main menu</span>
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
						</button>
					</div>
					<div
						class="mt-5 hidden w-full justify-end rounded-lg bg-zinc-700 p-4 backdrop-blur-3xl max-lg:py-4 lg:mt-0 lg:flex lg:bg-transparent lg:pl-10 lg:backdrop-blur-none"
						id="navbar"
					>
						<ul class="mt-4 flex flex-col gap-6 max-lg:mb-4 lg:mt-0 lg:flex-row lg:items-center">
							{@render NavItem(navigationBar)}
						</ul>
					</div>
				</div>
			</div>
		</nav>
	</div>
	{@render children?.()}
</main>
