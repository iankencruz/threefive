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
		type?: 'text' | 'email' | 'password' | 'tel' | 'url' | 'date' | 'number' | 'textarea' | 'media';
		placeholder?: string;
		required?: boolean;
		value?: string | number;
		colSpan?: number;
		class?: string;
		helperText?: string;
		rows?: number;
		inputClass?: string;
	}

	export interface FormConfig {
		fields: FormField[];
		submitText?: string;
		showSubmit?: boolean;
		columns?: number;
		submitVariant?: 'primary' | 'secondary' | 'tertiary';
		submitFullWidth?: boolean;
	}

	interface Props {
		config: FormConfig;
		formData?: Record<string, any>;
		onSubmit?: (data: Record<string, any>) => void | Promise<void>;
		onchange?: (data: Record<string, any>) => void;
		errors?: Record<string, string>;
		children?: import('svelte').Snippet;
		asForm?: boolean;
	}

	let {
		config,
		formData: initialFormData = {},
		onSubmit,
		onchange,
		errors = {},
		children,
		asForm = true
	}: Props = $props();

	// Helper function to get default values for all fields
	function getDefaultFormData(formConfig: FormConfig): Record<string, any> {
		if (!formConfig?.fields) return {};

		return formConfig.fields.reduce(
			(acc, field) => {
				// Media fields default to null, others to empty string
				acc[field.name] = field.value ?? (field.type === 'media' ? null : '');
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
	let selectedMediaCache = $state<Map<string, Media>>(new Map());

	// Notify parent of changes (with untrack to prevent infinite loops)
	$effect(() => {
		if (onchange && Object.keys(formData).length > 0) {
			untrack(() => onchange(formData));
		}
	});

	// Load media info for fields with existing values
	$effect.pre(() => {
		if (!config?.fields) {
			console.log('   ❌ No config.fields');
			return;
		}

		config.fields.forEach((field) => {
			if (field.type !== 'media') return;

			const mediaId = formData[field.name];

			if (mediaId && typeof mediaId === 'string' && !selectedMediaCache.has(mediaId)) {
				loadMediaInfo(mediaId);
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
			formData[currentMediaField] = mediaId;
			const newCache = new Map(selectedMediaCache);
			newCache.set(mediaId, media);
			selectedMediaCache = newCache;
			currentMediaField = '';
		}
		closeMediaPicker();
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

	// Helper to check if field type is media
	function isMediaField(type?: string): boolean {
		return type === 'media';
	}

	// Helper to get valid input type (excludes 'media')
	function getInputType(
		type?: string
	): 'text' | 'email' | 'password' | 'tel' | 'url' | 'date' | 'number' | 'textarea' | undefined {
		if (type === 'media') return undefined;
		return type as
			| 'text'
			| 'email'
			| 'password'
			| 'tel'
			| 'url'
			| 'date'
			| 'number'
			| 'textarea'
			| undefined;
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

{#if config?.fields}
	{#if asForm}
		<!-- Render as a form when asForm is true -->
		<form onsubmit={handleSubmit} class="space-y-6">
			<div class="grid grid-cols-{config.columns || 1} gap-4">
				{#each config.fields as field}
					<div class="{getColSpanClass(field.colSpan)} {field.class || ''}">
						{#if isMediaField(field.type)}
							{@render mediaFieldInput(field)}
						{:else}
							<Input
								type={getInputType(field.type) || 'text'}
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
						{/if}
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
			<div class="grid grid-cols-{config.columns || 1} gap-4">
				{#each config.fields as field}
					<div class="{getColSpanClass(field.colSpan)} {field.class || ''}">
						{#if isMediaField(field.type)}
							{@render mediaFieldInput(field)}
						{:else}
							<Input
								type={getInputType(field.type) || 'text'}
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
						{/if}
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
