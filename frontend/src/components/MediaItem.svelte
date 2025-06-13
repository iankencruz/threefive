<!-- MediaItem.svelte -->
<script lang="ts">
	import { updateMedia, deleteMedia } from '$lib/api/media';

	interface MediaItem {
		id: string; // Changed from string | number to just string
		url: string;
		title?: string;
		alt_text?: string;
		updated_at: Date;
	}

	interface MediaItemProps {
		item: MediaItem;
		onrefresh?: () => void;
		children?: (context: {
			item: MediaItem;
			title: string;
			altText: string;
			editMode: boolean;
			save: () => Promise<void>;
			remove: () => Promise<void>;
			toggleEdit: () => void;
			updateTitle: (value: string) => void;
			updateAltText: (value: string) => void;
		}) => any;
	}

	let { item, onrefresh, children }: MediaItemProps = $props();
	let title = $state<string>(item.title ?? '');
	let altText = $state<string>(item.alt_text ?? '');
	let editMode = $state<boolean>(false);

	async function save(): Promise<void> {
		await updateMedia(item.id, { title, alt_text: altText });
		editMode = false;
		console.log(item);
		onrefresh?.();
	}

	async function remove(): Promise<void> {
		if (confirm('Delete this media item?')) {
			await deleteMedia(item.id);
			onrefresh?.();
		}
	}

	function toggleEdit(): void {
		editMode = !editMode;
	}

	function updateTitle(value: string): void {
		title = value;
	}

	function updateAltText(value: string): void {
		altText = value;
	}
</script>

{@render children?.({
	item,
	title,
	altText,
	editMode,
	save,
	remove,
	toggleEdit,
	updateTitle,
	updateAltText
})}
