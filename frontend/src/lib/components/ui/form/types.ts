// frontend/src/lib/components/ui/form/types.ts
import type { Media } from "$api/media";

export type FieldType =
  | "text"
  | "email"
  | "password"
  | "tel"
  | "url"
  | "date"
  | "number"
  | "textarea"
  | "media"
  | "checkbox"
  | "select";

export interface SelectOption {
  value: string;
  label: string;
}

export interface FormField {
  name: string;
  label: string;
  type?: FieldType;
  placeholder?: string;
  required?: boolean;
  value?: string | number | boolean | string[];
  colSpan?: number;
  class?: string;
  helperText?: string;
  rows?: number;
  inputClass?: string;
  multiple?: boolean;
  options?: SelectOption[]; // For select fields
}

export interface FormConfig {
  fields: FormField[];
  submitText?: string;
  showSubmit?: boolean;
  columns?: number;
  submitVariant?: "primary" | "secondary" | "tertiary";
  submitFullWidth?: boolean;
}

// Common props interface for all field components
export interface BaseFieldProps {
  field: FormField;
  value: any;
  error?: string;
  disabled?: boolean;
  onchange: (value: any) => void;
}

// Extended props for media fields
export interface MediaFieldProps extends BaseFieldProps {
  mediaCache: Map<string, Media>;
  onMediaPickerOpen: (fieldName: string) => void;
  onMediaRemove?: (fieldName: string, mediaId?: string) => void;
  onMediaMove?: (fieldName: string, fromIndex: number, toIndex: number) => void;
}
