<script lang="ts">
	import { writable } from 'svelte/store';
	import { EyeOff, Eye } from '@lucide/svelte';

	const form = writable({
		first_name: '',
		last_name: '',
		email: '',
		password: ''
	});

	let error = $state('');
	let success = $state('');

	let fieldErrors = $state<{ [key: string]: string }>({});

	let showPassword = $state(false);

	function toLabel(key: string) {
		return key
			.replace(/_/g, ' ') // e.g. first_name => first name
			.replace(/\b\w/g, (c) => c.toUpperCase()); // capitalise each word
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();

		error = '';
		success = '';
		fieldErrors = {};

		const res = await fetch('/api/v1/admin/register', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include', // ✅ IMPORTANT
			body: JSON.stringify($form)
		});

		const data = await res.json();
		console.log('✅ JSON from backend:', data);

		if (!res.ok) {
			const errors = data.data || {};

			if (typeof errors === 'object' && Object.keys(errors).length > 0) {
				fieldErrors = errors;
				error = Object.values(errors).flat().join(', ') || data.message || 'Validation failed';
			} else {
				error = data.message || 'Something went wrong';
			}

			console.log('❗ Error:', errors);
			return;
		}
		success = 'Registered successfully! Redirecting...';
		setTimeout(() => (location.href = '/'), 1000);
	}

	$effect(() => {
		console.log('📦 fieldErrors changed:', fieldErrors);
	});
</script>

<svelte:head>
	<title>Register | ThreeFive</title>
	<meta name="description" content="Login to Sabiflow" />
</svelte:head>
<!--
  This example requires updating your template:

  ```
  <html class="h-full bg-gray-50">
  <body class="h-full">
  ```
-->
<div class="flex min-h-full flex-col justify-center py-12 sm:px-6 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-md">
		<img
			class="mx-auto h-10 w-auto"
			src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600"
			alt="Your Company"
		/>
		<h2 class="mt-6 text-center text-2xl/9 font-bold tracking-tight text-gray-900">Register</h2>
	</div>

	<div class="mt-10 sm:mx-auto sm:w-full sm:max-w-[480px]">
		<div class="bg-white px-6 py-12 shadow-sm sm:rounded-lg sm:px-12">
			<form class="space-y-4" action="#" method="POST" onsubmit={handleSubmit}>
				<div class="grid grid-cols-2 gap-x-4">
					<div>
						<label for="firstname" class="block text-sm/6 font-medium text-gray-900"
							>First name
						</label>
						<div class="mt-2">
							<input
								bind:value={$form.first_name}
								type="text"
								id="firstname"
								name="firstname"
								placeholder="First name"
								required
								class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
							/>
							{#if fieldErrors.first_name}
								<p class="mt-1 text-sm text-red-600">{fieldErrors.first_name}</p>
							{/if}
						</div>
					</div>
					<div>
						<label for="lastname" class="block text-sm/6 font-medium text-gray-900"
							>Last name
						</label>
						<div class="mt-2">
							<input
								bind:value={$form.last_name}
								type="text"
								id="lastname"
								name="lastname"
								placeholder="Last name"
								required
								class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
							/>
							{#if fieldErrors.last_name}
								<p class="mt-1 text-sm text-red-600">{fieldErrors.last_name}</p>
							{/if}
						</div>
					</div>
				</div>

				<div>
					<label for="email" class="block text-sm/6 font-medium text-gray-900">Email</label>
					<div class="mt-2">
						<input
							bind:value={$form.email}
							placeholder="Email"
							type="email"
							name="email"
							id="email"
							autocomplete="email"
							required
							class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
						/>
						{#if fieldErrors.email}
							<p class="mt-1 text-sm text-red-600">{fieldErrors.email}</p>
						{/if}
					</div>
				</div>
				<div>
					<label for="password" class="block text-sm/6 font-medium text-gray-900">Password</label>
					<div class=" mt-2">
						<div class="relative">
							<input
								bind:value={$form.password}
								type={showPassword ? 'text' : 'password'}
								name="password"
								id="password"
								placeholder="Password"
								minlength="8"
								autocomplete="current-password"
								required
								class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
							/>
							{#if $form.password != ''}
								<button
									type="button"
									onclick={() => (showPassword = !showPassword)}
									class="absolute inset-y-0 right-0 flex items-center px-3"
								>
									{#if showPassword}
										<span class="sr-only">Hide password</span>
										<EyeOff class="size-5" />
									{:else}
										<span class="sr-only">Show password</span>
										<Eye class="size-5" />
									{/if}
								</button>
							{/if}
						</div>
						{#if fieldErrors.password}
							<p class="mt-1 text-sm text-red-600">{fieldErrors.password}</p>
						{/if}
					</div>
				</div>

				<div>
					{#if Object.keys(fieldErrors).length > 0}
						<ul class="mb-4 rounded-md border border-red-200 bg-red-50 p-4 text-sm text-red-800">
							{#each Object.entries(fieldErrors) as [key, message]}
								<li>• {toLabel(key)}: {message}</li>
							{/each}
						</ul>
					{/if}

					{#if success}
						<p class="text-green-600">{success}</p>
					{/if}
				</div>

				<div>
					<button
						type="submit"
						class="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
						>Register</button
					>
				</div>
			</form>
		</div>

		<p class="mt-10 text-center text-sm/6 text-gray-500">
			Already a member?
			<a href="/admin/login" class="font-semibold text-indigo-600 hover:text-indigo-500">Login</a>
		</p>
	</div>
</div>
