<script lang="ts">
	import {
		LayoutDashboard,
		Menu,
		X,
		Files,
		Image,
		Tag,
		Cog,
		BookOpen,
		FolderClosed
	} from '@lucide/svelte';

	let isMobileMenuOpen = $state(false);
	const closeMenu = () => (isMobileMenuOpen = false);

	const navGroups = [
		{
			title: 'Home',
			items: [
				{ name: 'Dashboard', href: '/', icon: LayoutDashboard, active: true },
				{ name: 'Pages', href: 'pages', icon: Files, active: true },
				{ name: 'Projects', href: 'projects', icon: FolderClosed, active: true },
				{ name: 'Blogs', href: 'blogs', icon: BookOpen, active: false }
			]
		},
		{
			title: 'Documents',
			items: [
				{ name: 'Media', href: 'media', icon: Image, active: true },
				{ name: 'Tags', href: 'tags', icon: Tag, active: true },
				{ name: 'Settings', href: 'settings', icon: Cog, active: true }
			]
		}
	];
</script>

<div
	onclick={closeMenu}
	onkeydown={(e) => e.key === 'Escape' && closeMenu()}
	role="button"
	tabindex="0"
	class="fixed inset-0 z-30 bg-black/50 backdrop-blur-sm transition-opacity duration-300 lg:hidden
	{isMobileMenuOpen ? 'pointer-events-auto opacity-100' : 'pointer-events-none opacity-0'}"
></div>

<aside
	class="fixed inset-y-0 left-0 z-40 flex w-64 flex-col border-r border-border bg-sidebar text-sidebar-foreground opacity-100 transition-transform duration-300 lg:translate-x-0
	{isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full'}"
>
	<div class="flex h-16 items-center justify-between border-b border-border px-6">
		<div class="flex items-center gap-2.5">
			<div class="flex size-6 shrink-0 items-center justify-center rounded-full bg-primary">
				<div class="size-2 rounded-full bg-primary-foreground"></div>
			</div>
			<span class="text-sm font-bold tracking-tight">ThreeFive.</span>
		</div>
		<button onclick={closeMenu} class="text-muted-foreground hover:text-foreground lg:hidden">
			<X size={20} />
		</button>
	</div>

	<nav class="flex-1 space-y-6 overflow-y-auto p-3">
		{#each navGroups as group}
			<section>
				<h3
					class="mb-2 px-3 text-[10px] font-semibold tracking-wider text-muted-foreground uppercase"
				>
					{group.title}
				</h3>
				<div class="space-y-1">
					{#each group.items as item}
						{@const Icon = item.icon}
						<a
							href={item.href}
							onclick={closeMenu}
							class="rounded-radius flex items-center gap-3 px-3 py-2 text-sm font-medium transition-colors hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
							class:text-muted-foreground={item.active == false}
							class:pointer-events-none={item.active == false}
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

<button
	onclick={() => (isMobileMenuOpen = true)}
	class="fixed top-3 left-4 z-20 inline-flex items-center justify-center rounded-md p-2 text-muted-foreground hover:bg-accent lg:hidden"
>
	<Menu size={20} />
</button>
