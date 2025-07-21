
// src/lib/components/Builders/blockDefaults.ts

import type { Block } from '$lib/types';

export function getDefaultProps(type: Block['type']): Block['props'] {
  switch (type) {
    case 'heading':
      return { text: '', level: 1 };
    case 'richtext':
      return { html: '' };
    case 'image':
      return { media_id: '' };
    default:
      throw new Error(`Unsupported block type: ${type}`);
  }
}

