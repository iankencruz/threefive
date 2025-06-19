import dayjs from "dayjs";

// src/lib/utils/slugify.ts
export function slugify(str: string): string {
	return str
		.toLowerCase()
		.trim()
		.replace(/[^\w\s-]/g, '')
		.replace(/[\s_-]+/g, '-')
		.replace(/^-+|-+$/g, '');
}


export type FormatStyle = 'full' | 'relative';

export function formatDate(
	dateInput: string | number | Date | undefined | null,
	style: FormatStyle = 'full'
): string {
	if (!dateInput) return '—'; // Or 'loading...', or '' depending on context

	const date = dayjs(dateInput);

	if (!date.isValid()) return 'Invalid date';

	if (style === 'relative') {
		return date.fromNow();
	}

	return date.format('DD MMM YYYY, h:mm:ss A');
}
