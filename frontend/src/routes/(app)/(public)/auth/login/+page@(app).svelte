<script lang="ts">
	import DynamicForm, { type FormConfig } from '$components/ui/DynamicForm.svelte';
	import type { ErrorResponse } from '$types/auth';
	import { authApi } from '$api/auth';
	import { authStore } from '$stores/auth.svelte';
	import { goto } from '$app/navigation';

	interface FormData {
		first_name: string;
		last_name: string;
		email: string;
		password: string;
	}

	let formData: FormData = $state({
		first_name: '',
		last_name: '',
		email: '',
		password: ''
	});
	let errors = $state<Record<string, string>>({});
	let loading = $state(false);
	let generalError = $state('');

	const formConfig: FormConfig = {
		fields: [
			{
				name: 'email',
				type: 'email',
				label: 'Email',
				placeholder: 'Enter your email',
				required: true,
				colSpan: 12
			},
			{
				name: 'password',
				type: 'password',
				label: 'Password',
				placeholder: 'Enter your password',
				required: true,
				colSpan: 12
			}
		],
		submitText: 'Sign In',
		submitVariant: 'primary',
		submitFullWidth: true
	};

	async function handleSubmit(data: Record<string, any>) {
		loading = true;
		errors = {};
		generalError = '';

		try {
			const response = await authApi.login({
				email: data.email,
				password: data.password
			});

			authStore.setUser(response.user);
			goto('/admin');
		} catch (error) {
			const err = error as ErrorResponse;

			if (err.errors) {
				// Map validation errors to form fields
				err.errors.forEach((e) => {
					errors[e.field] = e.message;
				});
			} else {
				generalError = err.message || 'Login failed';
			}
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-background px-4">
	<div class="w-full max-w-md space-y-8">
		<div class="text-center">
			<h2 class="text-foreground">Sign in to your account</h2>
			<p class="mt-2">
				Or
				<a href="/auth/register" class="font-medium hover:underline"> create a new account </a>
			</p>
		</div>

		<div class="rounded-lg bg-surface p-8 shadow-md">
			{#if generalError}
				<div class="mb-4 rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700">
					{generalError}
				</div>
			{/if}

			<DynamicForm config={formConfig} {formData} {errors} onSubmit={handleSubmit} asForm={true} />

			<div class="mt-4 text-center">
				<a
					href="/auth/forgot-password"
					class="hover:text-primary-dark text-sm text-primary hover:underline"
				>
					Forgot your password?
				</a>
			</div>
		</div>
	</div>
</div>
