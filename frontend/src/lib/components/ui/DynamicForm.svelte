<!-- frontend/src/lib/components/ui/DynamicForm.svelte -->
<script lang="ts">
	import Button from "./Button.svelte";
	import Input from "./Input.svelte";
	import { ImageUp, CheckCircle } from "lucide-svelte";
	import { PUBLIC_API_URL } from "$env/static/public";
	import type { Media } from "$api/media";
	import { untrack } from "svelte";
	import MediaPicker from "./MediaPicker.svelte";

	interface FormField {
		name: string;
		label: string;
		type?:
			| "text"
			| "email"
			| "password"
			| "tel"
			| "url"
			| "date"
			| "number"
			| "textarea"
			| "media";
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
	}

	interface Props {
		config: FormConfig;
		formData?: Record<string, any>;
		onSubmit?: (data: Record<string, any>) => void | Promise<void>;
		onchange?: (data: Record<string, any>) => void;
		errors?: Record<string, string>;
		children?: import("svelte").Snippet;
		asForm?: boolean;
	}

	let {
		config,
		formData: initialFormData = {},
		onSubmit,
		onchange,
		errors = {},
		children,
		asForm = true,
	}: Props = $props();

	// Helper function to get default values for all fields
	function getDefaultFormData(formConfig: FormConfig): Record<string, any> {
		if (!formConfig?.fields) return {};

		return formConfig.fields.reduce(
			(acc, field) => {
				// Media fields default to null, others to empty string
				acc[field.name] = field.value ?? (field.type === "media" ? null : "");
				return acc;
			},
			{} as Record<string, any>,
		);
	}

	// Initialize formData immediately with defaults to prevent undefined binding
	let formData = $state<Record<string, any>>(
		config?.fields
			? initialFormData && Object.keys(initialFormData).length > 0
				? { ...getDefaultFormData(config), ...initialFormData }
				: getDefaultFormData(config)
			: {},
	);

	// Media picker state
	let showMediaPicker = $state(false);
	let currentMediaField = $state<string>("");
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
			console.log("   ❌ No config.fields");
			return;
		}

		config.fields.forEach((field) => {
			if (field.type !== "media") return;

			const mediaId = formData[field.name];

			if (
				mediaId &&
				typeof mediaId === "string" &&
				!selectedMediaCache.has(mediaId)
			) {
				loadMediaInfo(mediaId);
			}
		});
	});

	function getColSpanClass(colSpan?: number): string {
		if (!colSpan) return "col-span-1";
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
		currentMediaField = "";
	}

	function handleMediaSelect(mediaId: string, media: Media) {
		if (currentMediaField) {
			formData[currentMediaField] = mediaId;
			const newCache = new Map(selectedMediaCache);
			newCache.set(mediaId, media);
			selectedMediaCache = newCache;
			currentMediaField = "";
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
			console.log("   ❌ No mediaId");
			return;
		}

		if (selectedMediaCache.has(mediaId)) {
			console.log("   i  Already in cache");
			return;
		}

		try {
			const response = await fetch(
				`${PUBLIC_API_URL}/api/v1/media/${mediaId}`,
				{
					credentials: "include",
				},
			);

			if (response.ok) {
				const media = await response.json();
				const newCache = new Map(selectedMediaCache);
				newCache.set(mediaId, media);
				selectedMediaCache = newCache;
			}
		} catch (err) {
			console.error("   ❌ Failed to load media:", err);
		}
	}

	// Helper to check if field type is media
	function isMediaField(type?: string): boolean {
		return type === "media";
	}

	// Helper to get valid input type (excludes 'media')
	function getInputType(
		type?: string,
	):
		| "text"
		| "email"
		| "password"
		| "tel"
		| "url"
		| "date"
		| "number"
		| "textarea"
		| undefined {
		if (type === "media") return undefined;
		return type as
			| "text"
			| "email"
			| "password"
			| "tel"
			| "url"
			| "date"
			| "number"
			| "textarea"
			| undefined;
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return bytes + " B";
		if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
		return (bytes / (1024 * 1024)).toFixed(1) + " MB";
	}
</script>

{#snippet mediaFieldInput(field: FormField)}
	{@const mediaId = formData[field.name]}
	{@const hasMedia = mediaId && selectedMediaCache.has(mediaId)}
	{@const media = hasMedia ? selectedMediaCache.get(mediaId) : null}
	
	<div class="space-y-2">
		{#if field.label}
			<label class="block text-sm font-medium ">
				{field.label}
				{#if field.required}
					<span class="text-red-500">*</span>
				{/if}
			</label>
		{/if}
		
		

		{#if hasMedia && media}
			<div class="relative group border border-input-border rounded-lg p-4 bg-surface">
				<button
					type="button"
					onclick={() => clearMedia(field.name)}
					class="absolute top-2 right-2 p-1 bg-red-500 text-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity"
					title="Remove"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>

				<div class="flex items-center gap-4">
					<img
						src={media.thumbnail_url || media.url}
						alt={media.original_filename}
						class="w-20 h-20 object-cover rounded"
					/>
					<div class="flex-1 min-w-0">
						<p class="text-sm font-medium  truncate">{media.original_filename}</p>
						<p class="text-xs text-gray-500">
							{formatFileSize(media.size_bytes)}
							{#if media.width && media.height}
								• {media.width}×{media.height}
							{/if}
						</p>
						{#if media.mime_type === 'image/webp'}
							<span class="inline-flex items-center gap-1 mt-1 text-xs text-green-600 font-medium">
								<CheckCircle class="w-3 h-3" />
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
				class="w-full border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-gray-400 transition-colors"
			>
				<ImageUp class="mx-auto h-12 w-12 text-gray-400 mb-2" />
				<p class="text-sm text-gray-600">Click to select or upload media</p>
				<p class="text-xs text-gray-400 mt-1">Images will be converted to WebP • Videos will be optimized</p>
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
				<Button type="submit" class="w-full">
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
<MediaPicker 
	show={showMediaPicker}
	onselect={handleMediaSelect}
	onclose={closeMediaPicker}
/>
