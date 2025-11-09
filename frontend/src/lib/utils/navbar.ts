// src/lib/utils/navbar.ts

/**
 * Determines navbar variant based on page blocks
 * Returns 'transparent' if first block is a hero, 'opaque' otherwise
 */
export function getNavbarVariant(blocks: any[]): "transparent" | "opaque" {
	if (!blocks || blocks.length === 0) {
		return "opaque";
	}

	// Check if first block is a hero block
	const firstBlock = blocks[0];
	return firstBlock.type === "hero" ? "transparent" : "opaque";
}
