<script lang="ts">
	import DynamicForm, { type FormConfig } from '$components/ui/DynamicForm.svelte';
	import type { ErrorResponse } from '$types/auth';
	import { authApi } from '$api/auth';
	import { authStore } from '$stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';

	let formData = $state({
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
				name: 'first_name',
				type: 'text',
				label: 'First Name',
				placeholder: 'John',
				required: true,
				colSpan: 6
			},
			{
				name: 'last_name',
				type: 'text',
				label: 'Last Name',
				placeholder: 'Doe',
				required: true,
				colSpan: 6
			},
			{
				name: 'email',
				type: 'email',
				label: 'Email',
				placeholder: 'john@example.com',
				required: true,
				colSpan: 12
			},
			{
				name: 'password',
				type: 'password',
				label: 'Password',
				placeholder: 'Minimum 8 characters',
				required: true,
				colSpan: 12,
				helperText: 'Must contain uppercase, lowercase, number, and special character'
			}
		],
		submitText: 'Create Account',
		submitVariant: 'primary',
		submitFullWidth: true
	};

	function handleChange(data: Record<string, any>) {
		formData = {
			first_name: data.first_name || '',
			last_name: data.last_name || '',
			email: data.email || '',
			password: data.password || ''
		};
	}

	async function handleSubmit(data: Record<string, any>) {
		loading = true;
		errors = {};
		generalError = '';

		try {
			const response = await authApi.register({
				first_name: data.first_name,
				last_name: data.last_name,
				email: data.email,
				password: data.password
			});

			authStore.setUser(response.user);
			toast.success('Register success');
			goto('/admin');
		} catch (error) {
			const err = error as ErrorResponse;

			if (err.errors) {
				// Map validation errors to form fields
				err.errors.forEach((e) => {
					errors[e.field] = e.message;
					toast.error(e.message);
				});
			} else {
				generalError = err.message || 'Registration failed';
				toast.error(generalError);
			}
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-background px-4 py-12">
	<div class="w-full max-w-md space-y-8">
		<div class="text-center">
			<h2 class="">Create your account</h2>
			<p class="mt-2">
				Already have an account?
				<a href="/auth/login" class="hover:text-primary-dark font-medium text-primary"> Sign in </a>
			</p>
		</div>

		<div class="rounded-lg bg-surface p-8 shadow-md">
			{#if generalError}
				<div class="mb-4 rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700">
					{generalError}
				</div>
			{/if}

			<DynamicForm
				config={formConfig}
				{formData}
				{errors}
				onchange={handleChange}
				onSubmit={handleSubmit}
			/>
		</div>
	</div>
</div>
