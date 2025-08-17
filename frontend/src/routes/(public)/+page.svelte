<script lang="ts">
	import { Navigation, Pagination } from 'swiper/modules';

	import { onMount } from 'svelte';

	let { data } = $props();
	console.log('Home data:', data);
</script>

<!--custom css-->
<div class="">
	<section class="h-screen w-full">
		<div class="swiper mySwiper h-full">
			<div class="swiper-wrapper h-full">
				<div
					class="swiper-slide h-full min-h-[800px] w-full bg-[url('https://pagedone.io/asset/uploads/1720172752.png')] bg-cover bg-center bg-no-repeat pb-24 lg:pt-[84px]"
				>
					<section class=" pt-28 pb-0">
						<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
							<div class="inline-flex w-full flex-col items-start justify-start gap-14">
								<div class="inline-flex items-center justify-center gap-[492px]">
									<div class="inline-flex flex-col items-start justify-start gap-10 lg:gap-14">
										<div class="flex flex-col items-start justify-start gap-4">
											<h5 class="text-lg leading-relaxed font-medium text-white">
												Products / Headphone
											</h5>
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
			</div>
		</div>
	</section>
</div>
