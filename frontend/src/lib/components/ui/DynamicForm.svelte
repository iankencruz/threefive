<!-- frontend/src/lib/components/ui/DynamicForm.svelte -->
<script lang="ts">
	import Button from './Button.svelte';
	import Input from './Input.svelte';
	import { ImageUp, CircleCheck } from 'lucide-svelte';
	import { PUBLIC_API_URL } from '$env/static/public';
	import type { Media } from '$api/media';
	import { untrack } from 'svelte';
	import MediaPicker from './MediaPicker.svelte';

	interface FormField {
		name: string;
		label: string;
		type?:
			| 'text'
			| 'email'
			| 'password'
			| 'tel'
			| 'url'
			| 'date'
			| 'number'
			| 'textarea'
			| 'media'
			| 'checkbox'
			| 'select';
		placeholder?: string;
		required?: boolean;
		value?: string | number | boolean;
		colSpan?: number;
		class?: string;
		helperText?: string;
		rows?: number;
		inputClass?: string;
		multiple?: boolean;
	}

	export interface FormConfig {
		fields: FormField[];
		submitText?: string;
		showSubmit?: boolean;
		columns?: number;
		submitVariant?: 'primary' | 'secondary' | 'tertiary';
		submitFullWidth?: boolean;
		initialMediaCache?: Record<string, Media>;
	}

	interface Props {
		config: FormConfig;
		formData?: Record<string, any>;
		onSubmit?: (data: Record<string, any>) => void | Promise<void>;
		onchange?: (data: Record<string, any>) => void;
		errors?: Record<string, string>;
		children?: import('svelte').Snippet;
		asForm?: boolean;
		initialMediaCache?: Record<string, Media>;
	}

	let {
		config,
		formData: initialFormData = {},
		onSubmit,
		onchange,
		errors = {},
		children,
		asForm = true,
		initialMediaCache = {}
	}: Props = $props();

	// Helper function to get default values for all fields
	function getDefaultFormData(formConfig: FormConfig): Record<string, any> {
		if (!formConfig?.fields) return {};

		return formConfig.fields.reduce(
			(acc, field) => {
				// Different defaults based on field type
				if (field.type === 'media') {
					// Support both single and multiple media
					acc[field.name] = field.multiple ? [] : (field.value ?? null);
				} else if (field.type === 'checkbox') {
					acc[field.name] = field.value ?? false;
				} else {
					acc[field.name] = field.value ?? '';
				}
				return acc;
			},
			{} as Record<string, any>
		);
	}

	// Initialize formData immediately with defaults to prevent undefined binding
	let formData = $state<Record<string, any>>(
		config?.fields
			? initialFormData && Object.keys(initialFormData).length > 0
				? { ...getDefaultFormData(config), ...initialFormData }
				: getDefaultFormData(config)
			: {}
	);

	// Media picker state
	let showMediaPicker = $state(false);
	let currentMediaField = $state<string>('');
	let selectedMediaCache = $state<Map<string, Media>>(new Map(Object.entries(initialMediaCache)));

	// Notify parent of changes (with untrack to prevent infinite loops)
	$effect(() => {
		if (onchange && Object.keys(formData).length > 0) {
			untrack(() => onchange(formData));
		}
	});

	// Load media info for fields with existing values
	$effect.pre(() => {
		if (!config?.fields) {
			return;
		}

		config.fields.forEach((field) => {
			if (field.type !== 'media') return;

			if (field.multiple) {
				// Handle array of media IDs
				const mediaIds = formData[field.name] || [];
				mediaIds.forEach((mediaId: string) => {
					if (mediaId && typeof mediaId === 'string' && !selectedMediaCache.has(mediaId)) {
						loadMediaInfo(mediaId);
					}
				});
			} else {
				// Handle single media ID
				const mediaId = formData[field.name];
				if (mediaId && typeof mediaId === 'string' && !selectedMediaCache.has(mediaId)) {
					loadMediaInfo(mediaId);
				}
			}
		});
	});

	function getColSpanClass(colSpan?: number): string {
		if (!colSpan) return 'col-span-1';
		return `col-span-${colSpan}`;
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (onSubmit) {
			await onSubmit(formData);
		}
	}

	function openMediaPicker(fieldName: string) {
		currentMediaField = fieldName;
		showMediaPicker = true;
	}

	function closeMediaPicker() {
		showMediaPicker = false;
		currentMediaField = '';
	}

	function handleMediaSelect(mediaId: string, media: Media) {
		if (currentMediaField) {
			const field = config.fields.find((f) => f.name === currentMediaField);

			if (field?.multiple) {
				// Multi-select: add to array
				const currentMedia = formData[currentMediaField] || [];
				const exists = currentMedia.some((id: string) => id === mediaId);

				if (!exists) {
					formData[currentMediaField] = [...currentMedia, mediaId];
				}
			} else {
				// Single select: replace
				formData[currentMediaField] = mediaId;
			}

			const newCache = new Map(selectedMediaCache);
			newCache.set(mediaId, media);
			selectedMediaCache = newCache;
			currentMediaField = '';
		}
		closeMediaPicker();
	}

	function removeMediaFromArray(fieldName: string, mediaId: string) {
		const mediaIds = formData[fieldName] || [];
		formData[fieldName] = mediaIds.filter((id: string) => id !== mediaId);
	}

	function moveMediaInArray(fieldName: string, fromIndex: number, toIndex: number) {
		const mediaIds = [...(formData[fieldName] || [])];
		const [removed] = mediaIds.splice(fromIndex, 1);
		mediaIds.splice(toIndex, 0, removed);
		formData[fieldName] = mediaIds;
	}

	function clearMedia(fieldName: string) {
		const mediaId = formData[fieldName];
		formData[fieldName] = null;
		if (mediaId) {
			const newCache = new Map(selectedMediaCache);
			newCache.delete(mediaId);
			selectedMediaCache = newCache;
		}
	}

	// Load selected media info when field has a value
	async function loadMediaInfo(mediaId: string) {
		if (!mediaId) {
			console.log('   ❌ No mediaId');
			return;
		}

		if (selectedMediaCache.has(mediaId)) {
			console.log('   i  Already in cache');
			return;
		}

		try {
			const response = await fetch(`${PUBLIC_API_URL}/api/v1/media/${mediaId}`, {
				credentials: 'include'
			});

			if (response.ok) {
				const media = await response.json();
				const newCache = new Map(selectedMediaCache);
				newCache.set(mediaId, media);
				selectedMediaCache = newCache;
			}
		} catch (err) {
			console.error('   ❌ Failed to load media:', err);
		}
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return bytes + ' B';
		if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
		return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
	}
</script>

{#snippet mediaFieldInput(field: FormField)}
	{@const mediaId = formData[field.name]}
	{@const hasMedia = mediaId && selectedMediaCache.has(mediaId)}
	{@const media = hasMedia ? selectedMediaCache.get(mediaId) : null}

	<div class="space-y-2">
		{#if field.label}
			<label for={field.name} class="block text-sm font-medium">
				{field.label}
				{#if field.required}
					<span class="text-red-500">*</span>
				{/if}
			</label>
		{/if}

		{#if hasMedia && media}
			<div class="group relative rounded-lg border border-input-border bg-surface p-4">
				<button
					type="button"
					onclick={() => clearMedia(field.name)}
					class="absolute top-2 right-2 rounded-full bg-red-500 p-1 text-white opacity-0 transition-opacity group-hover:opacity-100"
					title="Remove"
				>
					<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						/>
					</svg>
				</button>

				<div class="flex items-center gap-4">
					<img
						src={media.thumbnail_url || media.url}
						alt={media.original_filename}
						class="h-20 w-20 rounded object-cover"
					/>
					<div class="min-w-0 flex-1">
						<p class="truncate text-sm font-medium">{media.original_filename}</p>
						<p class="text-xs text-gray-500">
							{formatFileSize(media.size_bytes)}
							{#if media.width && media.height}
								• {media.width}×{media.height}
							{/if}
						</p>
						{#if media.mime_type === 'image/webp'}
							<span class="mt-1 inline-flex items-center gap-1 text-xs font-medium text-green-600">
								<CircleCheck class="h-3 w-3" />
								WebP Optimized
							</span>
						{/if}
					</div>
					<button
						type="button"
						onclick={() => openMediaPicker(field.name)}
						class="px-3 py-1.5 text-sm font-medium text-blue-600 hover:text-blue-700"
					>
						Change
					</button>
				</div>
			</div>
		{:else}
			<button
				type="button"
				onclick={() => openMediaPicker(field.name)}
				class="w-full rounded-lg border-2 border-dashed border-gray-300 p-6 text-center transition-colors hover:border-gray-400"
			>
				<ImageUp class="mx-auto mb-2 h-12 w-12 text-gray-400" />
				<p class="text-sm text-gray-600">Click to select or upload media</p>
				<p class="mt-1 text-xs text-gray-400">
					Images will be converted to WebP • Videos will be optimized
				</p>
			</button>
		{/if}

		{#if errors[field.name]}
			<p class="text-sm text-red-600">{errors[field.name]}</p>
		{/if}
	</div>
{/snippet}

{#snippet checkboxFieldInput(field: FormField)}
	<div class="flex items-center gap-3">
		<div
			class="group relative inline-flex w-11 shrink-0 rounded-full bg-gray-200 p-0.5 inset-ring inset-ring-gray-900/5 outline-offset-2 outline-primary transition-colors duration-200 ease-in-out has-checked:bg-primary has-focus-visible:outline-2"
		>
			<span
				class="relative size-5 rounded-full bg-white shadow-xs ring-1 ring-gray-900/5 transition-transform duration-200 ease-in-out group-has-checked:translate-x-5"
			>
				<span
					aria-hidden="true"
					class="absolute inset-0 flex size-full items-center justify-center opacity-100 transition-opacity duration-200 ease-in group-has-checked:opacity-0 group-has-checked:duration-100 group-has-checked:ease-out"
				>
					<svg viewBox="0 0 12 12" fill="none" class="size-3 text-gray-400">
						<path
							d="M4 8l2-2m0 0l2-2M6 6L4 4m2 2l2 2"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
						/>
					</svg>
				</span>
				<span
					aria-hidden="true"
					class="absolute inset-0 flex size-full items-center justify-center opacity-0 transition-opacity duration-100 ease-out group-has-checked:opacity-100 group-has-checked:duration-200 group-has-checked:ease-in"
				>
					<svg viewBox="0 0 12 12" fill="currentColor" class="size-3 text-primary">
						<path
							d="M3.707 5.293a1 1 0 00-1.414 1.414l1.414-1.414zM5 8l-.707.707a1 1 0 001.414 0L5 8zm4.707-3.293a1 1 0 00-1.414-1.414l1.414 1.414zm-7.414 2l2 2 1.414-1.414-2-2-1.414 1.414zm3.414 2l4-4-1.414-1.414-4 4 1.414 1.414z"
						/>
					</svg>
				</span>
			</span>
			<input
				type="checkbox"
				id={field.name}
				name={field.name}
				bind:checked={formData[field.name]}
				aria-label={field.label}
				class="absolute inset-0 size-full appearance-none rounded-lg focus:outline-hidden"
			/>
		</div>
		{#if field.label}
			<label for={field.name} class="text-sm font-medium">
				{field.label}
				{#if field.required}
					<span class="text-red-500">*</span>
				{/if}
			</label>
		{/if}
	</div>
	{#if field.helperText}
		<p class="mt-1 ml-7 text-sm text-gray-500">{field.helperText}</p>
	{/if}
	{#if errors[field.name]}
		<p class="mt-1 ml-7 text-sm text-red-600">{errors[field.name]}</p>
	{/if}
{/snippet}

{#snippet standardFieldInput(field: FormField)}
	<Input
		type={field.type === 'checkbox' || field.type === 'media'
			? 'text'
			: (field.type as
					| 'text'
					| 'email'
					| 'password'
					| 'tel'
					| 'url'
					| 'date'
					| 'number'
					| 'textarea'
					| undefined)}
		name={field.name}
		label={field.label}
		placeholder={field.placeholder}
		required={field.required}
		bind:value={formData[field.name]}
		error={errors[field.name]}
		helperText={field.helperText}
		rows={field.rows}
		inputClass={field.inputClass}
	/>
{/snippet}

{#snippet multiMediaFieldInput(field: FormField)}
	{@const mediaIds = formData[field.name] || []}
	{@const mediaList = mediaIds.map((id: string) => selectedMediaCache.get(id)).filter(Boolean)}

	<div class="space-y-2">
		{#if field.label}
			<label class="block text-sm font-medium">
				{field.label}
				{#if field.required}
					<span class="text-red-500">*</span>
				{/if}
			</label>
		{/if}

		{#if mediaList.length > 0}
			<div class="space-y-2">
				{#each mediaList as media, index}
					{@const mediaId = mediaIds[index]}
					<div class="flex items-center gap-3 rounded-lg border border-gray-300 bg-gray-50 p-3">
						<img
							src={media.thumbnail_url || media.url}
							alt={media.original_filename}
							class="h-16 w-16 rounded object-cover"
						/>
						<div class="flex-1">
							<p class="text-sm font-medium text-gray-900">{media.original_filename}</p>
							<p class="text-xs text-gray-500">{media.mime_type}</p>
						</div>
						<div class="flex items-center gap-1">
							<button
								type="button"
								onclick={() => moveMediaInArray(field.name, index, index - 1)}
								disabled={index === 0}
								class="rounded p-1 text-gray-600 hover:bg-gray-200 disabled:opacity-30"
								title="Move up"
							>
								<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M5 15l7-7 7 7"
									/>
								</svg>
							</button>
							<button
								type="button"
								onclick={() => moveMediaInArray(field.name, index, index + 1)}
								disabled={index === mediaList.length - 1}
								class="rounded p-1 text-gray-600 hover:bg-gray-200 disabled:opacity-30"
								title="Move down"
							>
								<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M19 9l-7 7-7-7"
									/>
								</svg>
							</button>
							<button
								type="button"
								onclick={() => removeMediaFromArray(field.name, mediaId)}
								class="rounded p-1 text-red-600 hover:bg-red-100"
								title="Remove"
							>
								<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M6 18L18 6M6 6l12 12"
									/>
								</svg>
							</button>
						</div>
					</div>
				{/each}
			</div>

			<button
				type="button"
				onclick={() => openMediaPicker(field.name)}
				class="mt-2 flex w-full items-center justify-center gap-2 rounded-lg border-2 border-dashed border-gray-300 p-4 text-gray-600 transition-colors hover:border-gray-400 hover:text-gray-700"
			>
				<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 4v16m8-8H4"
					/>
				</svg>
				<span class="text-sm font-medium">Add More Images</span>
			</button>
		{:else}
			<button
				type="button"
				onclick={() => openMediaPicker(field.name)}
				class="w-full rounded-lg border-2 border-dashed border-gray-300 p-6 text-center transition-colors hover:border-gray-400"
			>
				<ImageUp class="mx-auto mb-2 h-12 w-12 text-gray-400" />
				<p class="text-sm text-gray-600">Click to select media</p>
				<p class="mt-1 text-xs text-gray-400">
					{field.required ? 'At least one file is required' : 'Optional'}
				</p>
			</button>
		{/if}

		{#if errors[field.name]}
			<p class="text-sm text-red-600">{errors[field.name]}</p>
		{/if}
	</div>
{/snippet}

{#snippet renderField(field: FormField)}
	{#if field.type === 'media'}
		{#if field.multiple}
			{@render multiMediaFieldInput(field)}
		{:else}
			{@render mediaFieldInput(field)}
		{/if}
	{:else if field.type === 'checkbox'}
		{@render checkboxFieldInput(field)}
	{:else}
		{@render standardFieldInput(field)}
	{/if}
{/snippet}

{#if config?.fields}
	{#if asForm}
		<!-- Render as a form when asForm is true -->
		<form onsubmit={handleSubmit} class="space-y-6">
			<div class="grid grid-cols-12 gap-4">
				{#each config.fields as field}
					<div class="{getColSpanClass(field.colSpan || 12)} {field.class || ''}">
						{@render renderField(field)}
					</div>
				{/each}
			</div>

			{#if children}
				{@render children()}
			{:else if config.showSubmit !== false}
				<Button type="submit" class="w-full py-2.5">
					{config.submitText || 'Submit'}
				</Button>
			{/if}
		</form>
	{:else}
		<!-- Render just the fields when asForm is false -->
		<div class="space-y-6">
			<div class="grid grid-cols-12 gap-4">
				{#each config.fields as field}
					<div class="{getColSpanClass(field.colSpan || 12)} {field.class || ''}">
						{@render renderField(field)}
					</div>
				{/each}
			</div>

			{#if children}
				{@render children()}
			{/if}
		</div>
	{/if}
{/if}

<!-- Media Picker Modal -->
<MediaPicker show={showMediaPicker} onselect={handleMediaSelect} onclose={closeMediaPicker} />
