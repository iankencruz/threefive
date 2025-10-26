
<!-- frontend/src/routes/admin/+layout.svelte -->
<script lang="ts">
	import { page } from "$app/state";

	const { children } = $props();

	const isActive = (path: string) =>
		page.url.pathname === path || page.url.pathname.startsWith(path + "/");

	const baseLink =
		"flex items-center gap-3 px-4 py-3 rounded-lg font-medium no-underline transition-colors";
	const linkClasses = (path: string) =>
		[
			baseLink,
			"text-foreground/60",
			isActive(path)
				? "bg-primary/10 text-primary"
				: "hover:bg-background/40 hover:text-foreground",
		].join(" ");
</script>

<div class="flex min-h-screen bg-background text-foreground">
	<aside class="hidden md:flex w-[260px] flex-col bg-surface border-r border-surface/60">
		<div class="px-6 py-8 border-b border-surface/60">
			<h2 class="m-0 text-xl font-bold text-foreground">Admin Panel</h2>
		</div>

		<nav class="flex flex-col gap-1 px-4 py-6">
			<a href="/admin" class={linkClasses("/admin")}>
				<svg class="icon h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
				</svg>
				Dashboard
			</a>

			<a href="/admin/pages" class={linkClasses("/admin/pages")}>
				<svg class="icon h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
				</svg>
				Pages
			</a>

			<a href="/admin/media" class={linkClasses("/admin/media")}>
				<svg class="icon h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
				</svg>
				Media
			</a>

			<a href="/admin/users" class={linkClasses("/admin/users")}>
				<svg class="icon h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
				</svg>
				Users
			</a>

			<a href="/admin/settings" class={linkClasses("/admin/settings")}>
				<svg class="icon h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
				</svg>
				Settings
			</a>
		</nav>
	</aside>

	<main class="flex-1 p-8 overflow-y-auto">
		{@render children?.()}
	</main>
</div>

