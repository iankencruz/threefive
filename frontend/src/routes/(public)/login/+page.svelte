<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/store/auth.svelte';
	import { onMount } from 'svelte';
	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	console.log('User - public  login', auth.user);

	// onMount(async () => {
	// 	try {
	// 		const res = await fetch('/api/v1/auth/me');
	// 		if (!res.ok) throw new Error('unauthenticated');
	//
	// 		const json = await res.json();
	// 		auth.setUser(json.user); // hydrate
	// 	} catch {
	// 		goto('/login'); // redirect if not logged in
	// 	} finally {
	// 		loading = false;
	// 	}
	// });
	//
	$effect(() => {
		if (auth.isAuthenticated) {
			goto('/admin');
		}
	});

	async function handleLogin(event: SubmitEvent) {
		event.preventDefault();
		error = '';
		loading = true;

		try {
			const res = await fetch('/api/v1/auth/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email, password })
			});

			const json = await res.json();

			if (!res.ok) {
				error = json.message || 'Login failed';
				return;
			}

			auth.setUser(json.data); // assumes backend returns user data
			goto('/admin');
		} catch (e) {
			console.error(e);
			error = 'Something went wrong. Try again.';
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-full w-full flex-col justify-center py-12 sm:px-6 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-md">
		<img
			class="mx-auto h-10 w-auto"
			src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600"
			alt="Your Company"
		/>
		<h2 class="mt-6 text-center text-2xl/9 font-bold tracking-tight text-gray-900">
			Sign in to your account
		</h2>
		{#if error}
			<p class="mt-2 text-center text-sm/6 text-red-600">{error}</p>
		{/if}
	</div>

	<div class="mt-10 sm:mx-auto sm:w-full sm:max-w-[480px]">
		<div class="bg-white px-6 py-12 shadow-sm sm:rounded-lg sm:px-12">
			<form class="space-y-6" onsubmit={handleLogin}>
				<div>
					<label for="email" class="block text-sm/6 font-medium text-gray-900">Email address</label>
					<div class="mt-2">
						<input
							type="email"
							name="email"
							id="email"
							bind:value={email}
							autocomplete="email"
							required
							class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
						/>
					</div>
				</div>

				<div>
					<label for="password" class="block text-sm/6 font-medium text-gray-900">Password</label>
					<div class="mt-2">
						<input
							type="password"
							name="password"
							id="password"
							bind:value={password}
							autocomplete="current-password"
							required
							class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
						/>
					</div>
				</div>

				<div class="flex items-center justify-between">
					<div class="flex gap-3">
						<div class="flex h-6 shrink-0 items-center">
							<div class="group grid size-4 grid-cols-1">
								<input
									id="remember-me"
									name="remember-me"
									type="checkbox"
									class="col-start-1 row-start-1 appearance-none rounded-sm border border-gray-300 bg-white checked:border-indigo-600 checked:bg-indigo-600 indeterminate:border-indigo-600 indeterminate:bg-indigo-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:border-gray-300 disabled:bg-gray-100 disabled:checked:bg-gray-100 forced-colors:appearance-auto"
								/>
								<svg
									class="pointer-events-none col-start-1 row-start-1 size-3.5 self-center justify-self-center stroke-white group-has-disabled:stroke-gray-950/25"
									viewBox="0 0 14 14"
									fill="none"
								>
									<path
										class="opacity-0 group-has-checked:opacity-100"
										d="M3 8L6 11L11 3.5"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									/>
									<path
										class="opacity-0 group-has-indeterminate:opacity-100"
										d="M3 7H11"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									/>
								</svg>
							</div>
						</div>
						<label for="remember-me" class="block text-sm/6 text-gray-900">Remember me</label>
					</div>

					<div class="text-sm/6">
						<a href="#" class="font-semibold text-indigo-600 hover:text-indigo-500"
							>Forgot password?</a
						>
					</div>
				</div>

				<div>
					<button
						type="submit"
						class="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
						disabled={loading}
					>
						{loading ? 'Signing in...' : 'Sign in'}
					</button>
				</div>
			</form>

			<div>
				<div class="relative mt-10">
					<div class="absolute inset-0 flex items-center" aria-hidden="true">
						<div class="w-full border-t border-gray-200"></div>
					</div>
					<div class="relative flex justify-center text-sm/6 font-medium">
						<span class="bg-white px-6 text-gray-900">Or continue with</span>
					</div>
				</div>

				<!-- Google OAuth2 -->
				<div class="mt-6 grid grid-cols-2 gap-4">
					<a
						href={`/auth/google/login`}
						class="flex w-full items-center justify-center gap-3 rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-xs ring-1 ring-gray-300 ring-inset hover:bg-gray-50 focus-visible:ring-transparent"
					>
						<svg class="h-5 w-5" viewBox="0 0 24 24" aria-hidden="true">
							<path
								d="M12.0003 4.75C13.7703 4.75 15.3553 5.36002 16.6053 6.54998L20.0303 3.125C17.9502 1.19 15.2353 0 12.0003 0C7.31028 0 3.25527 2.69 1.28027 6.60998L5.27028 9.70498C6.21525 6.86002 8.87028 4.75 12.0003 4.75Z"
								fill="#EA4335"
							/>
							<path
								d="M23.49 12.275C23.49 11.49 23.415 10.73 23.3 10H12V14.51H18.47C18.18 15.99 17.34 17.25 16.08 18.1L19.945 21.1C22.2 19.01 23.49 15.92 23.49 12.275Z"
								fill="#4285F4"
							/>
							<path
								d="M5.26498 14.2949C5.02498 13.5699 4.88501 12.7999 4.88501 11.9999C4.88501 11.1999 5.01998 10.4299 5.26498 9.7049L1.275 6.60986C0.46 8.22986 0 10.0599 0 11.9999C0 13.9399 0.46 15.7699 1.28 17.3899L5.26498 14.2949Z"
								fill="#FBBC05"
							/>
							<path
								d="M12.0004 24.0001C15.2404 24.0001 17.9654 22.935 19.9454 21.095L16.0804 18.095C15.0054 18.82 13.6204 19.245 12.0004 19.245C8.8704 19.245 6.21537 17.135 5.2654 14.29L1.27539 17.385C3.25539 21.31 7.3104 24.0001 12.0004 24.0001Z"
								fill="#34A853"
							/>
						</svg>
						<span class="text-sm/6 font-semibold">Google</span>
					</a>

					<a
						href="#"
						class="flex w-full items-center justify-center gap-3 rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-xs ring-1 ring-gray-300 ring-inset hover:bg-gray-50 focus-visible:ring-transparent"
					>
						<svg
							class="size-5 fill-[#24292F]"
							fill="currentColor"
							viewBox="0 0 20 20"
							aria-hidden="true"
						>
							<path
								fill-rule="evenodd"
								d="M10 0C4.477 0 0 4.484 0 10.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0110 4.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.203 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.942.359.31.678.921.678 1.856 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0020 10.017C20 4.484 15.522 0 10 0z"
								clip-rule="evenodd"
							/>
						</svg>
						<span class="text-sm/6 font-semibold">GitHub</span>
					</a>
				</div>
			</div>
		</div>

		<p class="mt-10 text-center text-sm/6 text-gray-500">
			Not a member?
			<a href="/register" class="font-semibold text-indigo-600 hover:text-indigo-500">Register</a>
		</p>
	</div>
</div>
