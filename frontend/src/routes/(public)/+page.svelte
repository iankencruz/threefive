<script>
	import { onMount } from 'svelte';

	let currentSlide = $state(0);
	let sliderContainer;
	let isTransitioning = $state(false);

	let { data } = $props();
	console.log('Home data:', data);
	console.log('Galleries data:', data.Galleries);

	// Find the specific gallery by slug
	// @ts-ignore
	const heroGallery = data.Galleries?.find((item) => item.gallery.slug === 'home-hero-gallery');

	// Create slides from the gallery media or fallback to default slides
	// @ts-ignore
	const slides = heroGallery?.gallery.media?.map((mediaItem) => ({
		bg: mediaItem.url || mediaItem.file_path, // adjust property name based on your media structure
		category: 'Products / Headphone',
		title: mediaItem.alt || mediaItem.title || 'AirPods Max',
		description:
			mediaItem.description ||
			'You can listen to music, make phone calls, use Siri, and more with your AirPods Max.'
	})) || [
		// Fallback slides if no gallery found
		{
			bg: 'https://pagedone.io/asset/uploads/1720172752.png',
			category: 'Products / Headphone',
			title: 'AirPods Max',
			description:
				'You can listen to music, make phone calls, use Siri, and more with your AirPods Max. The AirPods Max feature a stainless steel frame with a breathable knit mesh canopy and memory foam ear cushions for comfort.'
		}
	];

	// @ts-ignore
	let navOpen = $state(false);
	// @ts-ignore
	let dropdownOpen = $state(false);
	// @ts-ignore
	let megamenuOpen = $state(false);

	function nextSlide() {
		if (isTransitioning) return;
		isTransitioning = true;
		currentSlide = (currentSlide + 1) % slides.length;
		setTimeout(() => (isTransitioning = false), 500);
	}

	function prevSlide() {
		if (isTransitioning) return;
		isTransitioning = true;
		currentSlide = currentSlide === 0 ? slides.length - 1 : currentSlide - 1;
		setTimeout(() => (isTransitioning = false), 500);
	}

	/**
	 * @param {number} index
	 */
	function goToSlide(index) {
		if (isTransitioning) return;
		isTransitioning = true;
		currentSlide = index;
		setTimeout(() => (isTransitioning = false), 500);
	}

	onMount(() => {
		// Auto-play functionality (optional)
		const interval = setInterval(nextSlide, 5000);
		return () => clearInterval(interval);
	});
</script>

<div class="relative">
	<!-- Navigation -->

	<!-- Slider Section -->
	<section class="relative h-screen w-full overflow-hidden">
		<div class="relative h-full w-full" bind:this={sliderContainer}>
			{#each slides as slide, index}
				<div
					class="absolute inset-0 h-full min-h-[800px] w-full bg-cover bg-center bg-no-repeat pb-24 transition-transform duration-500 ease-in-out lg:pt-[84px]"
					style="background-image: url('{slide.bg}'); transform: translateX({(index -
						currentSlide) *
						100}%)"
				>
					<section class="pt-8 pt-[120px] pb-14">
						<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
							<div class="inline-flex w-full flex-col items-start justify-start gap-14">
								<div class="inline-flex items-center justify-center gap-[492px]">
									<div class="inline-flex flex-col items-start justify-start gap-10 lg:gap-14">
										<div class="flex flex-col items-start justify-start gap-4">
											<div
												class="prose prose-headings:text-white prose-h1:text-7xl prose-code:text-white flex flex-col items-start justify-start gap-2 text-white"
											>
												{@html data.Content}
											</div>
										</div>
										<button
											class="flex items-center justify-center rounded-xl bg-white px-5 py-2.5 shadow-[0px_1px_2px_0px_rgba(16,_24,_40,_0.05)] transition-all duration-700 ease-in-out hover:bg-gray-200"
										>
											<span class="px-2 py-px text-base leading-relaxed font-semibold text-gray-900"
												>Buy Now</span
											>
											<svg
												xmlns="http://www.w3.org/2000/svg"
												width="20"
												height="20"
												viewBox="0 0 20 20"
												fill="none"
											>
												<path
													d="M4.5845 4.99988L9.5847 10.0001L4.58154 15.0032M10.4178 4.99988L15.418 10.0001L10.4149 15.0032"
													stroke="#111827"
													stroke-width="1.6"
													stroke-linecap="round"
													stroke-linejoin="round"
												/>
											</svg>
										</button>
									</div>
								</div>
							</div>
						</div>
					</section>
				</div>
			{/each}
		</div>

		<!-- Navigation Controls -->
		<div
			class="absolute right-10 bottom-4 z-10 mx-auto flex max-w-[320px] -translate-x-1/2 transform items-center justify-between lg:bottom-24 xl:bottom-28 2xl:bottom-24"
		>
			<button
				onclick={prevSlide}
				class="group relative z-50 flex items-center justify-center p-2 transition-all duration-700 ease-in-out"
				disabled={isTransitioning}
			>
				<svg
					class="text-white transition-all duration-700 ease-in-out group-hover:text-gray-400"
					width="20"
					height="20"
					viewBox="0 0 20 20"
					fill="none"
				>
					<path
						d="M7.55112 16.0613L2.5 9.99994M2.5 9.99994L7.55112 3.9386M2.5 9.99994L17.4999 9.99994"
						stroke="currentColor"
						stroke-width="1.6"
						stroke-linecap="round"
						stroke-linejoin="round"
					/>
				</svg>
			</button>

			<!-- Progress Bar -->
			<div class="relative mx-14">
				<div class="h-0.5 w-48 rounded bg-gray-300">
					<div
						class="h-full rounded bg-indigo-400 transition-all duration-500 ease-in-out"
						style="width: {((currentSlide + 1) / slides.length) * 100}%"
					></div>
				</div>

				<!-- Slide Indicators (optional) -->
				<div class="absolute -bottom-8 left-1/2 z-10 flex -translate-x-1/2 transform space-x-2">
					{#each slides as _, index}
						<button
							onclick={() => goToSlide(index)}
							class="h-2 w-2 rounded-full transition-all duration-300 {index === currentSlide
								? 'bg-white'
								: 'bg-white/30'}"
							disabled={isTransitioning}
						></button>
					{/each}
				</div>
			</div>

			<button
				onclick={nextSlide}
				class="group relative z-50 flex items-center justify-center p-2 transition-all duration-700 ease-in-out"
				disabled={isTransitioning}
			>
				<svg
					class="text-white transition-all duration-700 ease-in-out group-hover:text-gray-400"
					width="20"
					height="20"
					viewBox="0 0 20 20"
					fill="none"
				>
					<path
						d="M12.4488 5L17.4999 10.0511M17.4999 10.0511L12.4488 15.1022M17.4999 10.0511L2.5 10.0511"
						stroke="currentColor"
						stroke-width="1.6"
						stroke-linecap="round"
						stroke-linejoin="round"
					/>
				</svg>
			</button>
		</div>

		<!-- Slide Counter -->
		<div class="absolute right-32 z-10 transform lg:top-[40%]">
			<div class="font-manrope text-3xl leading-[46px] font-normal text-gray-400">
				<span class="text-5xl leading-[62px] font-semibold text-white"
					>{(currentSlide + 1).toString().padStart(2, '0')}</span
				>
				/ {slides.length.toString().padStart(2, '0')}
			</div>
		</div>
	</section>
</div>

<style>
	.font-manrope {
		font-family: 'Manrope', sans-serif;
	}

	@media (max-width: 1003px) {
		.mx-14 {
			margin-left: 2rem;
			margin-right: 2rem;
		}
	}

	/* Custom scrollbar styling for progress bar */
	.w-48 {
		width: 12rem;
	}

	/* Animation for fade in effect */
	.animate-fade {
		animation: fadeIn 0.3s ease-in-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
