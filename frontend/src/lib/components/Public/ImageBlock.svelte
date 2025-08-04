<script lang="ts">
	let { block }: { block: { props: any } } = $props();
	const { url, alt_text = '', fit = '', align = '' } = block.props;

	const fitClass = fit === 'cover' ? 'object-cover' : fit === 'contain' ? 'object-contain' : '';

	const alignClass =
		align === 'left' ? 'ml-0 mr-auto' : align === 'right' ? 'ml-auto mr-0' : 'mx-auto';

	console.log('imageblock:', block);
</script>

<div class={`my-6 w-full ${alignClass}`}>
	<img src={url} alt={alt_text} class={`max-w-full rounded-md shadow ${fitClass}`} loading="lazy" />
</div>
