<!-- frontend/src/lib/components/ui/form/DynamicForm.svelte -->
<script lang="ts">
	import { untrack } from 'svelte';
	import { PUBLIC_API_URL } from '$env/static/public';
	import type { Media } from '$api/media';
	import Button from '$components/ui/Button.svelte';
	import MediaPicker from '$components/ui/MediaPicker.svelte';

	import TextField from './fields/TextField.svelte';
	import TextareaField from './fields/TextareaField.svelte';
	import SelectField from './fields/SelectField.svelte';
	import CheckboxField from './fields/CheckboxField.svelte';
	import MediaField from './fields/MediaField.svelte';
	import MultiMediaField from './fields/MultiMediaField.svelte';

	import type { FormConfig, FormField } from './types';

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

	// Initialize form data with defaults
	let formData = $state<Record<string, any>>(getDefaultFormData());

	// Media picker state
	let showMediaPicker = $state(false);
	let currentMediaField = $state<string>('');
	let selectedMediaCache = $state<Map<string, Media>>(new Map(Object.entries(initialMediaCache)));

	// Helper: Get default values for all fields
	function getDefaultFormData(): Record<string, any> {
		if (!config?.fields) return {};

		const defaults = config.fields.reduce(
			(acc, field) => {
				if (field.type === 'media') {
					acc[field.name] = field.multiple ? [] : null;
				} else if (field.type === 'checkbox') {
					acc[field.name] = field.value ?? false;
				} else {
					acc[field.name] = field.value ?? '';
				}
				return acc;
			},
			{} as Record<string, any>
		);

		return { ...defaults, ...initialFormData };
	}

	// Helper: Get grid column span class
	function getColSpanClass(colSpan?: number): string {
		if (!colSpan) return 'col-span-12';
		return `col-span-${colSpan}`;
	}

	// Notify parent of changes
	$effect(() => {
		if (onchange && Object.keys(formData).length > 0) {
			untrack(() => onchange(formData));
		}
	});

	// Load media info for existing values
	$effect.pre(() => {
		if (!config?.fields) return;

		config.fields.forEach((field) => {
			if (field.type !== 'media') return;

			if (field.multiple) {
				const mediaIds = formData[field.name] || [];
				mediaIds.forEach((id: string) => {
					if (id && !selectedMediaCache.has(id)) {
						loadMediaInfo(id);
					}
				});
			} else {
				const mediaId = formData[field.name];
				if (mediaId && !selectedMediaCache.has(mediaId)) {
					loadMediaInfo(mediaId);
				}
			}
		});
	});

	// Load media info from API
	async function loadMediaInfo(mediaId: string) {
		if (!mediaId || selectedMediaCache.has(mediaId)) return;

		try {
			const response = await fetch(`${PUBLIC_API_URL}/api/v1/media/${mediaId}`, {
				credentials: 'include'
			});

			if (response.ok) {
				const media = await response.json();
				selectedMediaCache = new Map(selectedMediaCache).set(mediaId, media);
			}
		} catch (err) {
			console.error('Failed to load media:', err);
		}
	}

	// Handle form submission
	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (onSubmit) {
			await onSubmit(formData);
		}
	}

	// Handle field value changes
	function handleFieldChange(fieldName: string, value: any) {
		formData[fieldName] = value;
	}

	// Media picker handlers
	function openMediaPicker(fieldName: string) {
		currentMediaField = fieldName;
		showMediaPicker = true;
	}

	function closeMediaPicker() {
		showMediaPicker = false;
		currentMediaField = '';
	}

	function handleMediaSelect(mediaId: string, media: Media) {
		if (!currentMediaField) return;

		const field = config.fields.find((f) => f.name === currentMediaField);

		if (field?.multiple) {
			const current = formData[currentMediaField] || [];
			if (!current.includes(mediaId)) {
				formData[currentMediaField] = [...current, mediaId];
			}
		} else {
			formData[currentMediaField] = mediaId;
		}

		selectedMediaCache = new Map(selectedMediaCache).set(mediaId, media);
		closeMediaPicker();
	}

	function removeMedia(fieldName: string, mediaId?: string) {
		const field = config.fields.find((f) => f.name === fieldName);

		if (field?.multiple && mediaId) {
			formData[fieldName] = (formData[fieldName] || []).filter((id: string) => id !== mediaId);
		} else {
			formData[fieldName] = null;
		}
	}

	function moveMedia(fieldName: string, fromIndex: number, toIndex: number) {
		const mediaIds = [...(formData[fieldName] || [])];
		const [removed] = mediaIds.splice(fromIndex, 1);
		mediaIds.splice(toIndex, 0, removed);
		formData[fieldName] = mediaIds;
	}
</script>

{#if config?.fields}
	{#if asForm}
		<form onsubmit={handleSubmit} class="space-y-6">
			<div class="grid grid-cols-12 gap-4">
				{#each config.fields as field}
					<div class="{getColSpanClass(field.colSpan || 12)} {field.class || ''}">
						{#if field.type === 'textarea'}
							<TextareaField
								{field}
								value={formData[field.name]}
								error={errors[field.name]}
								onchange={(value) => handleFieldChange(field.name, value)}
							/>
						{:else if field.type === 'select'}
							<SelectField
								{field}
								value={formData[field.name]}
								error={errors[field.name]}
								onchange={(value) => handleFieldChange(field.name, value)}
							/>
						{:else if field.type === 'checkbox'}
							<CheckboxField
								{field}
								value={formData[field.name]}
								error={errors[field.name]}
								onchange={(value) => handleFieldChange(field.name, value)}
							/>
						{:else if field.type === 'media' && field.multiple}
							<MultiMediaField
								{field}
								value={formData[field.name]}
								error={errors[field.name]}
								mediaCache={selectedMediaCache}
								onMediaPickerOpen={openMediaPicker}
								onMediaRemove={removeMedia}
								onMediaMove={moveMedia}
								onchange={(value) => handleFieldChange(field.name, value)}
							/>
						{:else if field.type === 'media'}
							<MediaField
								{field}
								value={formData[field.name]}
								error={errors[field.name]}
								mediaCache={selectedMediaCache}
								onMediaPickerOpen={openMediaPicker}
								onMediaRemove={removeMedia}
								onchange={(value) => handleFieldChange(field.name, value)}
							/>
						{:else}
							<TextField
								{field}
								value={formData[field.name]}
								error={errors[field.name]}
								onchange={(value) => handleFieldChange(field.name, value)}
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
		<div class="space-y-6">
			<div class="grid grid-cols-12 gap-4">
				{#each config.fields as field}
					<div class="{getColSpanClass(field.colSpan || 12)} {field.class || ''}">
						{#if field.type === 'textarea'}
							<TextareaField
								{field}
								value={formData[field.name]}
								error={errors[field.name]}
								onchange={(value) => handleFieldChange(field.name, value)}
							/>
						{:else if field.type === 'select'}
							<SelectField
								{field}
								value={formData[field.name]}
								error={errors[field.name]}
								onchange={(value) => handleFieldChange(field.name, value)}
							/>
						{:else if field.type === 'checkbox'}
							<CheckboxField
								{field}
								value={formData[field.name]}
								error={errors[field.name]}
								onchange={(value) => handleFieldChange(field.name, value)}
							/>
						{:else if field.type === 'media' && field.multiple}
							<MultiMediaField
								{field}
								value={formData[field.name]}
								error={errors[field.name]}
								mediaCache={selectedMediaCache}
								onMediaPickerOpen={openMediaPicker}
								onMediaRemove={removeMedia}
								onMediaMove={moveMedia}
								onchange={(value) => handleFieldChange(field.name, value)}
							/>
						{:else if field.type === 'media'}
							<MediaField
								{field}
								value={formData[field.name]}
								error={errors[field.name]}
								mediaCache={selectedMediaCache}
								onMediaPickerOpen={openMediaPicker}
								onMediaRemove={removeMedia}
								onchange={(value) => handleFieldChange(field.name, value)}
							/>
						{:else}
							<TextField
								{field}
								value={formData[field.name]}
								error={errors[field.name]}
								onchange={(value) => handleFieldChange(field.name, value)}
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

<MediaPicker show={showMediaPicker} onselect={handleMediaSelect} onclose={closeMediaPicker} />
