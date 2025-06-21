<script lang="ts">
	import { auth } from '$lib/store/auth.svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { type NavigationItem } from '$lib/components/Navigation/navigation';

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
			console.log('res data: ', res);

			const json = await res.json();
			console.log('json data: ', json);
			auth.setUser(json.user); // ✅ this rehydrates the user
		} catch (err) {
			console.error('Not logged in: Please sign in to access admin', err);
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
	{@render children?.()}
{/if}
