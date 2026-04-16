<script lang="ts">
	import type { NavItem } from '$lib/types';
	import { FilesIcon, FolderOpenIcon, ImageIcon, LayoutDashboard, TagIcon } from '@lucide/svelte';

	let { path }: { path: string } = $props();

	const adminNavItems: NavItem[] = [
		{
			title: 'Dashboard',
			href: '/admin/dashboard',
			icon: LayoutDashboard
		},
		{
			title: 'Pages',
			href: '/admin/pages',
			icon: FilesIcon
		},
		{
			title: 'Projects',
			href: '/admin/projects',
			icon: FolderOpenIcon
		},
		{
			title: 'Media',
			href: '/admin/media',
			icon: ImageIcon
		},
		{
			title: 'Tags',
			href: '/admin/tags',
			icon: TagIcon
		}
	];
</script>

<ul role="list" class="-mx-2 space-y-1">
	{#each adminNavItems as item}
		<li>
			<a
				href={item.href}
				class="group flex items-center gap-x-3 rounded-md p-2 text-sm/6 font-semibold transition-colors"
				class:text-primary={path === item.href}
				class:bg-muted={path === item.href}
				class:text-foreground={path !== item.href}
				class:hover:bg-muted={path !== item.href}
				class:hover:text-primary={path !== item.href}
				class:dark:text-gray-400={path !== item.href}
				class:dark:hover:bg-white={path !== item.href}
				class:dark:hover:text-white={path !== item.href}
			>
				{#if item.icon}
					{@const Icon = item.icon}
					<Icon size={18} weight="regular" class="shrink-0" />
				{/if}
				{item.title}
			</a>
		</li>
	{/each}
</ul>
