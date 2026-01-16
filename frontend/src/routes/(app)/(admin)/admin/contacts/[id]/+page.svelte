<!-- frontend/src/routes/(app)/(admin)/admin/contacts/[id]/+page.svelte -->
<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { enhance } from '$app/forms';
	import type { PageData, ActionData } from './$types';
	import { ArrowLeft, Mail, Calendar, User, MessageSquare, Trash2 } from 'lucide-svelte';

	let { data, form }: { data: PageData; form: ActionData } = $props();

	let isUpdating = $state(false);
	let isDeleting = $state(false);
	let showDeleteConfirm = $state(false);

	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
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

	// Auto-mark as read when viewing
	$effect(() => {
		if (data.contact.status === 'new') {
			const formData = new FormData();
			formData.append('status', 'read');

			fetch(`/api/v1/admin/contacts/${data.contact.id}/status`, {
				method: 'PATCH',
				credentials: 'include',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ status: 'read' })
			}).then(() => {
				invalidateAll();
			});
		}
	});
</script>

<svelte:head>
	<title>Contact: {data.contact.name}</title>
</svelte:head>

<div class="mx-auto max-w-4xl">
	<!-- Header -->
	<div class="mb-8">
		<button
			onclick={() => goto('/admin/contacts')}
			class="mb-4 flex items-center gap-2 text-sm text-gray-400 transition-colors hover:text-gray-200"
		>
			<ArrowLeft size={16} />
			Back to contacts
		</button>

		<div class="flex items-start justify-between">
			<div>
				<h1 class="mb-2">Contact Submission</h1>
				<p class="text-sm text-gray-400">
					Received {formatDate(data.contact.created_at)}
				</p>
			</div>

			<div class="flex items-center gap-3">
				<!-- Status Badge -->
				<span
					class="inline-flex items-center rounded-full px-3 py-1 text-sm font-medium capitalize {getStatusColor(
						data.contact.status
					)}"
				>
					{data.contact.status}
				</span>

				<!-- Delete Button -->
				<button
					type="button"
					onclick={() => (showDeleteConfirm = true)}
					class="rounded-lg border border-red-600 px-4 py-2 text-sm font-medium text-red-600 transition-colors hover:bg-red-600 hover:text-white"
					disabled={isDeleting}
				>
					<Trash2 size={16} class="inline-block" />
				</button>
			</div>
		</div>
	</div>

	<!-- Contact Information Card -->
	<div class="mb-6 overflow-hidden rounded-lg bg-surface shadow">
		<div class="border-b border-gray-700 px-6 py-4">
			<h2 class="text-lg font-semibold">Contact Information</h2>
		</div>

		<div class="grid gap-6 p-6 md:grid-cols-2">
			<!-- Name -->
			<div class="flex items-start gap-3">
				<User size={20} class="mt-0.5 text-gray-400" />
				<div>
					<p class="text-sm text-gray-400">Name</p>
					<p class="font-medium">{data.contact.name}</p>
				</div>
			</div>

			<!-- Email -->
			<div class="flex items-start gap-3">
				<Mail size={20} class="mt-0.5 text-gray-400" />
				<div>
					<p class="text-sm text-gray-400">Email</p>
					<a href="mailto:{data.contact.email}" class="font-medium text-primary hover:underline">
						{data.contact.email}
					</a>
				</div>
			</div>

			<!-- Subject (if exists) -->
			{#if data.contact.subject}
				<div class="flex items-start gap-3 md:col-span-2">
					<MessageSquare size={20} class="mt-0.5 text-gray-400" />
					<div>
						<p class="text-sm text-gray-400">Subject</p>
						<p class="font-medium">{data.contact.subject}</p>
					</div>
				</div>
			{/if}

			<!-- Date -->
			<div class="flex items-start gap-3 md:col-span-2">
				<Calendar size={20} class="mt-0.5 text-gray-400" />
				<div>
					<p class="text-sm text-gray-400">Submitted</p>
					<p class="font-medium">{formatDate(data.contact.created_at)}</p>
				</div>
			</div>
		</div>
	</div>

	<!-- Message Card -->
	<div class="mb-6 overflow-hidden rounded-lg bg-surface shadow">
		<div class="border-b border-gray-700 px-6 py-4">
			<h2 class="text-lg font-semibold">Message</h2>
		</div>
		<div class="p-6">
			<p class="whitespace-pre-wrap text-gray-300">{data.contact.message}</p>
		</div>
	</div>

	<!-- Actions Card -->
	<div class="overflow-hidden rounded-lg bg-surface shadow">
		<div class="border-b border-gray-700 px-6 py-4">
			<h2 class="text-lg font-semibold">Actions</h2>
		</div>

		<div class="p-6">
			<form
				method="POST"
				action="?/updateStatus"
				use:enhance={() => {
					isUpdating = true;
					return async ({ update }) => {
						await update();
						isUpdating = false;
					};
				}}
			>
				<label class="mb-2 block text-sm font-medium">Change Status</label>
				<div class="flex gap-3">
					<select
						name="status"
						class="form-input flex-1"
						onchange={(e) => e.currentTarget.form?.requestSubmit()}
						disabled={isUpdating}
					>
						<option value="new" selected={data.contact.status === 'new'}>New</option>
						<option value="read" selected={data.contact.status === 'read'}>Read</option>
						<option value="archived" selected={data.contact.status === 'archived'}>Archived</option>
					</select>
				</div>

				{#if form?.success}
					<p class="mt-2 text-sm text-green-600">Status updated successfully!</p>
				{/if}
				{#if form?.error}
					<p class="mt-2 text-sm text-red-600">{form.error}</p>
				{/if}
			</form>

			<!-- Quick Reply Button -->
			<div class="mt-6">
				<a
					href="mailto:{data.contact.email}?subject=Re: {data.contact.subject || 'Your message'}"
					class="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 font-medium text-white transition-colors hover:bg-primary/90"
				>
					<Mail size={18} />
					Reply via Email
				</a>
			</div>
		</div>
	</div>
</div>

<!-- Delete Confirmation Modal -->
{#if showDeleteConfirm}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		onclick={() => (showDeleteConfirm = false)}
	>
		<div
			class="w-full max-w-md rounded-lg bg-surface p-6 shadow-xl"
			onclick={(e) => e.stopPropagation()}
		>
			<h3 class="mb-4 text-lg font-semibold">Delete Contact Submission</h3>
			<p class="mb-6 text-gray-400">
				Are you sure you want to delete this contact submission? This action cannot be undone.
			</p>

			<div class="flex justify-end gap-3">
				<button
					type="button"
					onclick={() => (showDeleteConfirm = false)}
					class="rounded-lg border border-gray-600 px-4 py-2 font-medium transition-colors hover:bg-gray-700"
				>
					Cancel
				</button>

				<form
					method="POST"
					action="?/delete"
					use:enhance={() => {
						isDeleting = true;
						return async ({ result, update }) => {
							await update();
							isDeleting = false;
							if (result.type === 'success' && result.data?.deleted) {
								goto('/admin/contacts');
							}
						};
					}}
				>
					<button
						type="submit"
						disabled={isDeleting}
						class="rounded-lg bg-red-600 px-4 py-2 font-medium text-white transition-colors hover:bg-red-700 disabled:cursor-not-allowed disabled:opacity-50"
					>
						{isDeleting ? 'Deleting...' : 'Delete'}
					</button>
				</form>
			</div>
		</div>
	</div>
{/if}
