<script lang="ts">
	import type { Snippet } from "svelte";

	interface Props {
		type?: "button" | "submit" | "reset";
		variant?: "primary" | "secondary" | "outline" | "danger";
		size?: "sm" | "md" | "lg";
		fullWidth?: boolean;
		disabled?: boolean;
		loading?: boolean;
		class?: string;
		onclick?: (e: MouseEvent) => void;
		children?: Snippet;
	}

	let {
		type = "button",
		variant = "primary",
		size = "md",
		fullWidth = false,
		disabled = false,
		loading = false,
		class: className = "",
		onclick,
		children,
	}: Props = $props();

	const variantClasses = {
		primary: "btn-primary",
		secondary: "btn-secondary",
		outline: "btn-outline",
		danger: "btn-danger",
	};

	const sizeClasses = {
		sm: "btn-sm",
		md: "",
		lg: "btn-lg",
	};

	const classes = $derived(
		[
			"btn",
			variantClasses[variant],
			sizeClasses[size],
			fullWidth && "btn-full",
			loading && "opacity-50 cursor-wait",
			className,
		]
			.filter(Boolean)
			.join(" "),
	);
</script>

<button
	{type}
	class={classes}
	disabled={disabled || loading}
	onclick={onclick}
>
	{#if loading}
		<span class="inline-block animate-spin mr-2">⏳</span>
	{/if}
	{#if children}
		{@render children()}
{/if}
</button>
