<script lang="ts">
	import { Tipex } from '@friendofsvelte/tipex';
	import '@friendofsvelte/tipex/styles/index.css';
	// let body = `<p>The initial html content.</p>`;

	let body = $state<string>(
		'<p>Type your content here. You can use <strong>bold</strong>, <em>italic</em>, and <u>underline</u>.</p>'
	);

	let editor = $state();
</script>

<!-- Shorthands to disable: !floating !focal -->
<Tipex {body} floating focal class="mt-2 mb-0 h-[20vh] border border-neutral-200" />
