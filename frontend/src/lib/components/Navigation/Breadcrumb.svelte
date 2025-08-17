<script lang="ts">
	import { page } from '$app/stores';
	import { ChevronRight } from '@lucide/svelte';

	let { lastLabel }: { lastLabel?: string } = $props();

	let breadcrumbs = $state<{ label: string; href: string }[]>([]);
	let uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

	let shouldHide = $state(false);

	$effect(() => {
		let pathname = $page.url.pathname;

		if (pathname === '/admin') {
			breadcrumbs = []; // force clear
			shouldHide = true;
			return;
		}

		let segments = $page.url.pathname.split('/').filter(Boolean); // e.g. ['admin', 'projects', 'uuid']
		let path = '';
		let parts: { label: string; href: string }[] = [];

		segments.forEach((segment, index) => {
			const isAdmin = index === 0 && segment === 'admin';
			const isLast = index === segments.length - 1;
			const isUUID = uuidRegex.test(segment);

			path += `/${segment}`;

			parts.push({
				label: isAdmin
					? 'Dashboard'
					: isLast && isUUID && lastLabel
						? lastLabel
						: segment.charAt(0).toUpperCase() + segment.slice(1),
				href: isAdmin ? '/admin' : path
			});
		});

		breadcrumbs = parts;
	});
</script>

{#if !shouldHide}
	<nav class="mb-4 px-4 text-sm text-gray-500">
		<ul class="flex items-center gap-1">
			{#each breadcrumbs as crumb, i (i)}
				{#if i > 0}
					<li class="text-gray-400">
						<ChevronRight class="h-4 w-4" />
					</li>
				{/if}
				<li>
					<a
						href={crumb.href}
						class="block max-w-64 truncate hover:underline {i === breadcrumbs.length - 1
							? 'font-medium text-gray-800'
							: 'text-gray-500'}"
					>
						{crumb.label}
					</a>
				</li>
			{/each}
		</ul>
	</nav>
{/if}
