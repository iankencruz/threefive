
<!-- frontend/src/routes/admin/+layout.svelte -->
<script lang="ts">
	import { page } from "$app/state";
	import {
		BookOpenText,
		Files,
		House,
		Image,
		ScrollText,
		Settings,
	} from "lucide-svelte";

	const { children } = $props();

	const isActive = (path: string, exact = false) =>
		exact
			? page.url.pathname === path
			: page.url.pathname === path || page.url.pathname.startsWith(path + "/");

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
			<a href="/admin/dashboard" class={linkClasses("/admin/dashboard")}>
				<House size={18} />
				Dashboard
			</a>

			<a href="/admin/pages" class={linkClasses("/admin/pages")}>
				<BookOpenText size={18} />
				Pages
			</a>
			<a href="/admin/projects" class={linkClasses("/admin/projects")}>
				<Files size={18} />
				Projects
			</a>
			<a href="/admin/blogs" class={linkClasses("/admin/blogs")}>
				<ScrollText size={18}/>
				Blogs
			</a>

			<a href="/admin/media" class={linkClasses("/admin/media")}>
				<Image size={18} />
				Media
			</a>


			<a href="/admin/settings" class={linkClasses("/admin/settings")}>
				<Settings size={18}/>
				Settings
			</a>
		</nav>
	</aside>

	<main class="flex-1 p-8 overflow-y-auto">
		{@render children?.()}
	</main>
</div>

