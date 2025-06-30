<script lang="ts">
	import '$src/app.css';
	import { auth } from '$lib/store/auth.svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import {
		sidebarNavigation,
		userMenuItems,
		type NavigationItem
	} from '$lib/components/Navigation/navigation';
	import { PanelLeftClose, PanelLeftOpen } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';

	let { children } = $props();

	let loading = $state(true);

	let isSidebarOpen = $state(false);
	let menuOpen = $state(false);
	let collapsed = $state(false); // controls full vs compact sidebar

	function isActive(path?: string): boolean {
		return page.url.pathname === path || page.url.pathname.startsWith(path + '/');
	}

	onMount(async () => {
		try {
			const res = await fetch('/api/v1/auth/me');
			if (!res.ok) throw new Error('Unauthenticated');

			const json = await res.json();
			auth.setUser(json.data); // ✅ this rehydrates the user
		} catch (err) {
			console.error('Not logged in: Please sign in to access admin', err);
			toast.error('You must be logged in to access this page');
			goto('/login'); // ✅ ensure you redirect if not logged in
		} finally {
			loading = false;
		}
	});

	async function handleLogout() {
		await auth.logout();
		goto('/login');
	}
</script>

<!-- NavItem Snippet -->
{#snippet NavItem(NavItems: NavigationItem[])}
	{#each NavItems as item}
		<a
			href={item.href}
			class={`group  flex items-center rounded-md text-sm font-semibold text-gray-700
			${!collapsed && isActive(item.href) ? 'bg-gray-50 text-indigo-600' : ''}
			${collapsed && isActive(item.href) ? 'text-indigo-600' : ''}
			${!collapsed ? 'gap-x-4 px-2 py-1  hover:bg-gray-50' : 'justify-center px-2 py-1 hover:bg-transparent'}
			hover:text-indigo-600`}
		>
			<!-- icon container (always visible) -->
			<div class="flex h-10 w-10 items-center justify-center">
				<item.icon class="h-5 w-5 shrink-0" />
			</div>

			<!-- label (conditionally visible) -->
			<span class={collapsed ? 'sr-only' : 'truncate'}>{item.label}</span>
		</a>
	{/each}
{/snippet}

<!-- UserProfile Navigation loop snippet -->
{#snippet UserItems(NavItems: NavigationItem[])}
	{#each NavItems as item}
		<button
			onclick={handleLogout}
			class="flex w-full items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
		>
			<item.icon class="h-4 w-4" />
			<span>{item.label}</span>
		</button>
	{/each}
{/snippet}

{#if loading}
	<p class="mt-10 text-center text-gray-500">Loading...</p>
{:else}
	<div>
		<!-- Mobile sidebar -->
		{#if isSidebarOpen}
			<div class="fixed inset-0 z-50 flex lg:hidden" role="dialog" aria-modal="true">
				<button
					aria-label="Sidebar Open"
					class="fixed inset-0 bg-gray-900/80"
					onclick={() => (isSidebarOpen = false)}
				></button>
				<div class="relative mr-16 flex w-full max-w-xs flex-1">
					<div class="absolute top-0 left-full flex w-16 justify-center pt-5">
						<button class="-m-2.5 p-2.5" onclick={() => (isSidebarOpen = false)}>
							<span class="sr-only">Close sidebar</span>
							<svg
								class="size-6 text-white"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="1.5"
							>
								<path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
							</svg>
						</button>
					</div>
					<div class="flex grow flex-col gap-y-5 overflow-y-auto bg-white px-6 pb-2">
						<div class="flex h-16 shrink-0 items-center">
							<img
								class="h-8 w-auto"
								src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600"
								alt="Logo"
							/>
						</div>
						<nav class="flex flex-1 flex-col">
							<ul class="flex flex-1 flex-col gap-y-2">
								<li>
									<ul class="-mx-2 space-y-1">
										{@render NavItem(sidebarNavigation)}
									</ul>
								</li>
							</ul>
						</nav>
					</div>
				</div>
			</div>
		{/if}
		<!-- Desktop sidebar -->
		<div
			class={`hidden lg:fixed lg:inset-y-0 lg:z-50 lg:flex lg:flex-col ${
				collapsed ? 'lg:w-20' : 'lg:w-72'
			}`}
		>
			<div
				class="flex grow flex-col gap-y-5 overflow-y-auto border-r border-gray-200 bg-white px-6"
			>
				<div class="flex h-16 shrink-0 items-center">
					<img
						class="h-8 w-auto"
						src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600"
						alt="Logo"
					/>
				</div>

				<div class="flex">
					<button
						title={collapsed ? 'Open Sidebar' : 'Colllapse Sidebar'}
						onclick={() => (collapsed = !collapsed)}
						aria-label="Toggle sidebar"
						class={`flex items-center justify-center rounded-md text-gray-500 hover:bg-gray-100
						${collapsed ? 'mx-auto h-10 w-10' : '  h-10 w-10  '}`}
					>
						{#if collapsed}
							<PanelLeftOpen class="h-6 w-6" />
						{:else}
							<PanelLeftClose class="h-6 w-6" />
						{/if}
					</button>
				</div>
				<nav class="flex flex-1 flex-col">
					<ul class="flex flex-1 flex-col">
						<li>
							<ul class="-mx-2">
								{@render NavItem(sidebarNavigation)}
							</ul>
						</li>
					</ul>
				</nav>
				<!-- User menu at the bottom -->
				<div class="mt-auto pb-6">
					<div class="relative -mx-2 hover:bg-gray-50">
						<button
							class="flex w-full items-center gap-3 rounded-lg p-2"
							onclick={() => (menuOpen = !menuOpen)}
						>
							<img
								class="h-8 w-8 rounded-full"
								src={`https://ui-avatars.com/api/?name=${auth.user?.first_name}+${auth.user?.last_name}`}
								alt="User avatar"
							/>
							<span class="truncate text-sm font-medium text-gray-900"
								>{auth.user?.first_name} {auth.user?.last_name}</span
							>
						</button>

						{#if menuOpen}
							<div
								class="absolute bottom-14 left-0 z-20 w-full origin-top-left rounded-md bg-white shadow ring-1 ring-black/5"
							>
								{@render UserItems(userMenuItems)}
							</div>
						{/if}
					</div>
				</div>
			</div>
		</div>
		<!-- Mobile top bar -->
		<div
			class="sticky top-0 z-40 flex items-center gap-x-6 bg-white px-4 py-4 shadow sm:px-6 lg:hidden"
		>
			<button class="-m-2.5 p-2.5 text-gray-700" onclick={() => (isSidebarOpen = true)}>
				<span class="sr-only">Open sidebar</span>
				<svg
					class="size-6"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="1.5"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
					/>
				</svg>
			</button>
			<div class="flex-1 text-sm font-semibold text-gray-900">Dashboard</div>
			<div class="relative -mx-2 hover:bg-gray-50">
				<button
					class="flex w-full items-center gap-3 rounded-lg p-2"
					onclick={() => (menuOpen = !menuOpen)}
				>
					<img
						class="h-8 w-8 rounded-full"
						src={`https://ui-avatars.com/api/?name=${auth.user?.first_name}+${auth.user?.last_name}`}
						alt="User avatar"
					/>
					<span class="truncate text-sm font-medium text-gray-900"
						>{auth.user?.first_name} {auth.user?.last_name}</span
					>
				</button>

				{#if menuOpen}
					<div
						class="absolute left-0 z-20 w-full origin-top-left rounded-md bg-white shadow ring-1 ring-black/5"
					>
						{@render UserItems(userMenuItems)}
					</div>
				{/if}
			</div>
		</div>
		<main
			class={`min-h-screen  overflow-x-hidden px-4 py-10 sm:px-4  ${collapsed ? 'lg:ml-20' : 'lg:ml-72'}`}
		>
			{@render children()}
		</main>
	</div>{/if}
