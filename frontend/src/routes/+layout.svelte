<script lang="ts">
	import './layout.css';
	import type { Snippet } from 'svelte';
	import {
		LayoutDashboard,
		RefreshCw,
		ChartColumn,
		Folder,
		Users,
		Database,
		FileText,
		PenTool,
		Ellipsis,
		Plus,
		Menu,
		X
	} from '@lucide/svelte';

	let { children }: { children: Snippet } = $props();
	let isMobileMenuOpen = $state(false);

	const navGroups = [
		{
			title: 'Home',
			items: [
				{ name: 'Dashboard', href: '#', icon: LayoutDashboard },
				{ name: 'Lifecycle', href: '#', icon: RefreshCw },
				{ name: 'Analytics', href: '#', icon: ChartColumn },
				{ name: 'Projects', href: '#', icon: Folder },
				{ name: 'Team', href: '#', icon: Users }
			]
		},
		{
			title: 'Documents',
			items: [
				{ name: 'Data Library', href: '#', icon: Database },
				{ name: 'Reports', href: '#', icon: FileText },
				{ name: 'Word Assistant', href: '#', icon: PenTool },
				{ name: 'More', href: '#', icon: Ellipsis }
			]
		}
	];

	const closeMenu = () => (isMobileMenuOpen = false);
</script>

<div class="bg-background text-foreground flex min-h-screen font-sans">
	<div
		onclick={closeMenu}
		onkeydown={(e) => e.key === 'Escape' && closeMenu()}
		role="button"
		tabindex="0"
		class="fixed inset-0 z-30 bg-black/50 backdrop-blur-sm transition-opacity duration-300 lg:hidden
		{isMobileMenuOpen ? 'pointer-events-auto opacity-100' : 'pointer-events-none opacity-0'}"
	></div>

	<aside
		class="border-border bg-sidebar text-sidebar-foreground fixed inset-y-0 left-0 z-40 flex w-64 flex-col border-r opacity-100 transition-transform duration-300 lg:translate-x-0
		{isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full'}"
	>
		<div class="border-border flex h-16 items-center justify-between border-b px-6">
			<div class="flex items-center gap-2.5">
				<div class="bg-primary flex size-6 shrink-0 items-center justify-center rounded-full">
					<div class="bg-primary-foreground size-2 rounded-full"></div>
				</div>
				<span class="text-sm font-bold tracking-tight">Acme Inc.</span>
			</div>
			<button onclick={closeMenu} class="text-muted-foreground hover:text-foreground lg:hidden">
				<X size={20} />
			</button>
		</div>

		<nav class="flex-1 space-y-6 overflow-y-auto p-3">
			{#each navGroups as group}
				<section>
					<h3
						class="text-muted-foreground mb-2 px-3 text-[10px] font-semibold tracking-wider uppercase"
					>
						{group.title}
					</h3>
					<div class="space-y-1">
						{#each group.items as item}
							{@const Icon = item.icon}
							<a
								href={item.href}
								onclick={closeMenu}
								class="rounded-radius hover:bg-sidebar-accent hover:text-sidebar-accent-foreground flex items-center gap-3 px-3 py-2 text-sm font-medium transition-colors"
								class:bg-sidebar-accent={item.name === 'Dashboard'}
								class:text-sidebar-accent-foreground={item.name === 'Dashboard'}
							>
								<Icon size={16} strokeWidth={2} class="opacity-80" />
								{item.name}
							</a>
						{/each}
					</div>
				</section>
			{/each}
		</nav>
	</aside>

	<div class="flex w-full flex-1 flex-col lg:pl-64">
		<header
			class="border-border bg-background/80 sticky top-0 z-10 flex h-16 items-center justify-between border-b px-4 backdrop-blur-md md:px-6"
		>
			<div class="flex items-center gap-3">
				<button
					onclick={() => (isMobileMenuOpen = true)}
					class="text-muted-foreground hover:bg-accent inline-flex items-center justify-center rounded-md p-2 lg:hidden"
				>
					<Menu size={20} />
				</button>
				<h1 class="text-lg font-semibold tracking-tight">Documents</h1>
			</div>

			<div class="flex items-center gap-3">
				<button
					class="bg-primary text-primary-foreground flex items-center gap-2 rounded-md px-3.5 py-1.5 text-xs font-semibold transition-opacity hover:opacity-90"
				>
					<Plus size={14} strokeWidth={3} />
					<span class="xs:inline hidden">Quick Create</span>
				</button>
			</div>
		</header>

		<main class="w-full max-w-7xl p-6 lg:p-10">
			{@render children()}
		</main>
	</div>
</div>
