<!-- frontend/src/routes/(app)/(admin)/admin/contacts/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import type { PageData } from './$types';
	import { Mail, Eye, MessageSquare } from 'lucide-svelte';

	let { data }: { data: PageData } = $props();

	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	};

	const getStatusColor = (status: string) => {
		switch (status) {
			case 'new':
				return 'bg-blue-100 text-blue-800';
			case 'read':
				return 'bg-green-100 text-green-800';
			case 'archived':
				return 'bg-gray-100 text-gray-800';
			default:
				return 'bg-gray-100 text-gray-800';
		}
	};

	const currentStatus = $derived(page.url.searchParams.get('status') || 'all');

	const updateStatusFilter = (newStatus: string) => {
		const params = new URLSearchParams(page.url.searchParams);
		if (newStatus === 'all') {
			params.delete('status');
		} else {
			params.set('status', newStatus);
		}
		params.set('page', '1');
		goto(`/admin/contacts?${params.toString()}`);
	};
</script>

<svelte:head>
	<title>Admin: Contacts</title>
</svelte:head>

{#snippet pagination(data: PageData)}
	{#if data.pagination && data.pagination.total_pages > 1}
		<div class="mt-8 flex items-center justify-center gap-4">
			<button
				class="rounded-lg border border-gray-300 bg-surface px-4 py-2 font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
				disabled={data.pagination.page === 1}
				onclick={() => {
					const params = new URLSearchParams(page.url.searchParams);
					params.set('page', (data.pagination.page - 1).toString());
					goto(`/admin/contacts?${params.toString()}`);
				}}
			>
				Previous
			</button>
			<span class="text-sm text-gray-600">
				Page {data.pagination.page} of {data.pagination.total_pages}
			</span>
			<button
				class="rounded-lg border border-gray-300 bg-surface px-4 py-2 font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
				disabled={data.pagination.page === data.pagination.total_pages}
				onclick={() => {
					const params = new URLSearchParams(page.url.searchParams);
					params.set('page', (data.pagination.page + 1).toString());
					goto(`/admin/contacts?${params.toString()}`);
				}}
			>
				Next
			</button>
		</div>
	{/if}
{/snippet}

<div class="mx-auto max-w-7xl">
	<div class="mb-8 flex items-center justify-between">
		<div>
			<h1 class="mb-2">Contact Submissions</h1>
			<p class="text-sm text-gray-400">
				{data.pagination.total_count} total submission{data.pagination.total_count !== 1 ? 's' : ''}
			</p>
		</div>
	</div>

	<!-- Status Filter Tabs -->
	<div class="mb-6 border-b border-gray-700">
		<nav class="-mb-px flex gap-8">
			<button
				type="button"
				onclick={() => updateStatusFilter('all')}
				class="border-b-2 px-1 py-4 text-sm font-medium transition-colors {currentStatus === 'all'
					? 'border-primary text-primary'
					: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-200'}"
			>
				All
			</button>
			<button
				type="button"
				onclick={() => updateStatusFilter('new')}
				class="border-b-2 px-1 py-4 text-sm font-medium transition-colors {currentStatus === 'new'
					? 'border-primary text-primary'
					: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-200'}"
			>
				New
			</button>
			<button
				type="button"
				onclick={() => updateStatusFilter('read')}
				class="border-b-2 px-1 py-4 text-sm font-medium transition-colors {currentStatus === 'read'
					? 'border-primary text-primary'
					: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-200'}"
			>
				Read
			</button>
			<button
				type="button"
				onclick={() => updateStatusFilter('archived')}
				class="border-b-2 px-1 py-4 text-sm font-medium transition-colors {currentStatus ===
				'archived'
					? 'border-primary text-primary'
					: 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-200'}"
			>
				Archived
			</button>
		</nav>
	</div>

	{#if !data.contacts || data.contacts.length === 0}
		<!-- Empty State -->
		<div class="rounded-lg border-2 border-dashed border-gray-700 bg-surface py-20 text-center">
			<MessageSquare class="mx-auto mb-4 h-12 w-12 text-gray-400" />
			<h3 class="mb-2 text-lg font-medium text-gray-200">No contact submissions</h3>
			<p class="text-gray-400">
				{currentStatus === 'all'
					? 'Contact submissions will appear here.'
					: `No ${currentStatus} submissions.`}
			</p>
		</div>
	{:else}
		<!-- Contacts Table -->
		<div class="overflow-hidden rounded-lg bg-surface shadow">
			<table class="min-w-full divide-y divide-gray-700">
				<thead class="bg-surface">
					<tr>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase"> Name </th>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
							Email
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
							Subject
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
							Status
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase"> Date </th>
						<th class="px-6 py-3 text-left text-xs font-medium tracking-wider uppercase">
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-700 bg-surface">
					{#each data.contacts as contact}
						<tr class="hover:bg-white/5">
							<td class="px-6 py-4">
								<div class="flex items-center gap-2">
									{#if contact.status === 'new'}
										<span class="inline-block h-2 w-2 rounded-full bg-blue-500" title="Unread"
										></span>
									{/if}
									<span class="font-medium">{contact.name}</span>
								</div>
							</td>
							<td class="px-6 py-4">
								<a
									href="mailto:{contact.email}"
									class="flex items-center gap-2 text-sm text-gray-300 hover:text-primary"
								>
									<Mail size={16} />
									{contact.email}
								</a>
							</td>
							<td class="px-6 py-4 text-sm text-gray-300">
								{contact.subject || '—'}
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span
									class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium capitalize {getStatusColor(
										contact.status
									)}"
								>
									{contact.status}
								</span>
							</td>
							<td class="px-6 py-4 text-sm whitespace-nowrap text-gray-400">
								{formatDate(contact.created_at)}
							</td>
							<td class="px-6 py-4 text-sm whitespace-nowrap">
								<div class="flex items-center gap-2">
									<button
										onclick={() => goto(`/admin/contacts/${contact.id}`)}
										class="rounded p-2 text-gray-400 transition-colors hover:bg-gray-700 hover:text-white"
										title="View details"
									>
										<Eye size={18} />
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		{@render pagination(data)}
	{/if}
</div>
