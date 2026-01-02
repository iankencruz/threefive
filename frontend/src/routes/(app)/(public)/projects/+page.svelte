<script lang="ts">
	// You can add any necessary script logic here

	let { data } = $props();
	console.log('Projects List Page Data:', data);
</script>

<div class="grid grid-cols-3 gap-2">
	<div>
		<img
			class="rounded-base h-auto w-full"
			src="https://flowbite.s3.amazonaws.com/docs/gallery/square/image-1.jpg"
			alt=""
		/>
	</div>
	<div>
		<img
			class="rounded-base h-auto w-full"
			src="https://flowbite.s3.amazonaws.com/docs/gallery/square/image-2.jpg"
			alt=""
		/>
	</div>
	<div>
		<img
			class="rounded-base h-auto w-full"
			src="https://flowbite.s3.amazonaws.com/docs/gallery/square/image-3.jpg"
			alt=""
		/>
	</div>
	<div>
		<img
			class="rounded-base h-auto w-full"
			src="https://flowbite.s3.amazonaws.com/docs/gallery/square/image-4.jpg"
			alt=""
		/>
	</div>
	<div>
		<img
			class="rounded-base h-auto w-full"
			src="https://flowbite.s3.amazonaws.com/docs/gallery/square/image-5.jpg"
			alt=""
		/>
	</div>

	<div>
		<img
			class="rounded-base h-auto w-full"
			src="https://flowbite.s3.amazonaws.com/docs/gallery/square/image-6.jpg"
			alt=""
		/>
	</div>
</div>
