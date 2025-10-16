<!-- frontend/src/lib/components/ui/DynamicForm.svelte -->
<script lang="ts">
import { ImageUp, CheckIcon } from "lucide-svelte";
import { PUBLIC_API_URL } from "$env/static/public";
import type { Media } from "$api/media";
import MediaPicker from "$components/ui/MediaPicker.svelte";
import Input from "$components/ui/Input.svelte";
import Button from "$components/ui/Button.svelte";

export interface HeroBlockData {
	title: string;
	subtitle?: string;
	image_id?: string; // ✅ Changed from optional to match backend
	cta_text?: string;
	cta_url?: string;
}

interface FormField {
	name: string;
	label: string;
	type?: "text" | "email" | "password" | "tel" | "url" | "date" | "number" | "textarea" | "media";
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
	formData: initialFormData,
	onSubmit,
	onchange,
	errors = {},
	children,
	asForm = true,
}: Props = $props();

// Helper function to get default values for all fields
function getDefaultFormData(): Record<string, any> {
	return config.fields.reduce(
		(acc, field) => {
			acc[field.name] = field.value ?? (field.type === "media" ? "" : "");
			return acc;
		},
		{} as Record<string, any>,
	);
}

// Initialize formData - merge initialFormData with defaults to ensure all fields exist
let formData = $state<Record<string, any>>(
	initialFormData && Object.keys(initialFormData).length > 0
		? { ...getDefaultFormData(), ...initialFormData }
		: getDefaultFormData(),
);

// Media picker state
let showMediaPicker = $state(false);
let currentMediaField = $state<string>("");
let selectedMediaCache = $state<Map<string, Media>>(new Map());

// Notify parent of changes
$effect(() => {
	if (onchange) {
		onchange(formData);
	}
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

function handleMediaSelect(mediaId: string, media: Media) {
	if (currentMediaField) {
		formData[currentMediaField] = mediaId;
		selectedMediaCache.set(mediaId, media);
		currentMediaField = "";
	}
}

function clearMedia(fieldName: string) {
	formData[fieldName] = "";
	selectedMediaCache.delete(formData[fieldName]);
}

// Load selected media info when field has a value
async function loadMediaInfo(mediaId: string) {
	if (!mediaId || selectedMediaCache.has(mediaId)) return;

	try {
		const response = await fetch(`${PUBLIC_API_URL}/api/v1/media/${mediaId}`, {
			credentials: "include",
		});
		if (response.ok) {
			const media = await response.json();
			selectedMediaCache.set(mediaId, media);
		}
	} catch (err) {
		console.error("Failed to load media:", err);
	}
}

// Helper to check if field type is media
function isMediaField(type?: string): boolean {
	return type === "media";
}

// Helper to get valid input type (excludes 'media')
function getInputType(
	type?: string,
): "text" | "email" | "password" | "tel" | "url" | "date" | "number" | "textarea" | undefined {
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

// Load media info for fields with existing values
$effect(() => {
	config.fields.forEach((field) => {
		if (field.type === "media" && formData[field.name]) {
			loadMediaInfo(formData[field.name]);
		}
	});
});
</script>

{#snippet mediaFieldInput(field: FormField)}
	<div class="space-y-2">
		{#if field.label}
			<label class="block text-sm font-medium text-gray-700">
				{field.label}
				{#if field.required}
					<span class="text-red-500">*</span>
				{/if}
			</label>
		{/if}

		{#if formData[field.name] && selectedMediaCache.has(formData[field.name])}
			{@const media = selectedMediaCache.get(formData[field.name])}
			{#if media}
				<div class="relative group border-2 border-gray-200 rounded-lg p-4 bg-gray-50">
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
							<p class="text-sm font-medium text-gray-900 truncate">{media.original_filename}</p>
							<p class="text-xs text-gray-500">
								{formatFileSize(media.size_bytes)}
								{#if media.width && media.height}
									• {media.width}×{media.height}
								{/if}
							</p>
							{#if media.mime_type === 'image/webp'}
								<span class="inline-flex items-center gap-1 mt-1 text-xs text-green-600 font-medium">
									<CheckIcon class="w-3 h-3" />
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
			{/if}
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

<!-- Media Picker Modal -->
<MediaPicker
	show={showMediaPicker}
	onselect={handleMediaSelect}
/>
