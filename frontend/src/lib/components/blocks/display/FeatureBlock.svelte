<script lang="ts">
	import type { FeatureBlockData } from '../forms/FeatureBlockForm.svelte';

	interface Props {
		data: FeatureBlockData;
	}

	let { data }: Props = $props();
</script>

<section class=" bg-background py-12">
	<div class="mx-auto mt-12 w-full max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
		<div class="grid grid-cols-1 gap-8 md:grid-cols-8">
			<h2
				class="col-span-2 border-r border-gray-200/5 text-2xl/tight font-semibold sm:text-xl/tight"
			>
				{data?.title}
			</h2>

			<div class="col-span-6 mb-8 ml-32 space-y-6">
				<h3 class="mt-4 w-full text-8xl">{data?.heading}</h3>
				<h4 class="mt-8 text-4xl">{data?.subheading}</h4>
				<p class="text w-full max-w-2xl py-6 text-lg font-semibold text-gray-400">
					{data?.description}
				</p>
				<a href={`/projects`} class="inline-block">
					<button
						class="rounded-lg bg-primary px-6 py-3 text-lg font-semibold text-white transition-all hover:bg-primary/90"
					>
						Discover All Projects
					</button>
				</a>
			</div>
		</div>
		<div class="mt-8 grid grid-cols-1">
			<div class=" rounded-lg border border-gray-200 dark:border-gray-700">
				<a href="#" class="text-lg font-medium text-primary hover:underline">
					<img
						src={data?.media?.url}
						alt={data?.heading}
						class="h-full max-h-[32rem] w-full rounded-lg object-cover"
					/>
				</a>
			</div>
		</div>
	</div>
</section>
