<script lang="ts">
import DynamicForm, {
	type FormConfig,
} from "$components/ui/DynamicForm.svelte";
import type { ErrorResponse } from "$types/auth";
import { authApi } from "$api/auth";
import { authStore } from "$stores/auth.svelte";
import { goto } from "$app/navigation";

let formData = $state({
	first_name: "",
	last_name: "",
	email: "",
	password: "",
});
let errors = $state<Record<string, string>>({});
let loading = $state(false);
let generalError = $state("");

const formConfig: FormConfig = {
	fields: [
		{
			name: "first_name",
			type: "text",
			label: "First Name",
			placeholder: "John",
			required: true,
			colSpan: 6,
		},
		{
			name: "last_name",
			type: "text",
			label: "Last Name",
			placeholder: "Doe",
			required: true,
			colSpan: 6,
		},
		{
			name: "email",
			type: "email",
			label: "Email",
			placeholder: "john@example.com",
			required: true,
			colSpan: 12,
		},
		{
			name: "password",
			type: "password",
			label: "Password",
			placeholder: "Minimum 8 characters",
			required: true,
			colSpan: 12,
			helperText:
				"Must contain uppercase, lowercase, number, and special character",
		},
	],
	submitText: "Create Account",
	submitVariant: "primary",
	submitFullWidth: true,
};

function handleChange(data: Record<string, any>) {
	formData = {
		first_name: data.first_name || "",
		last_name: data.last_name || "",
		email: data.email || "",
		password: data.password || "",
	};
}

async function handleSubmit(data: Record<string, any>) {
	loading = true;
	errors = {};
	generalError = "";

	try {
		const response = await authApi.register({
			first_name: data.first_name,
			last_name: data.last_name,
			email: data.email,
			password: data.password,
		});

		authStore.setUser(response.user);
		goto("/admin");
	} catch (error) {
		const err = error as ErrorResponse;

		if (err.errors) {
			// Map validation errors to form fields
			err.errors.forEach((e) => {
				errors[e.field] = e.message;
			});
		} else {
			generalError = err.message || "Registration failed";
		}
	} finally {
		loading = false;
	}
}
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50 px-4 py-12">
	<div class="max-w-md w-full space-y-8">
		<div class="text-center">
			<h2 class="text-3xl font-bold text-gray-900">Create your account</h2>
			<p class="mt-2 text-sm text-gray-600">
				Already have an account?
				<a href="/auth/login" class="font-medium text-primary hover:text-primary-dark">
					Sign in
				</a>
			</p>
		</div>

		<div class="bg-white p-8 rounded-lg shadow-md">
			{#if generalError}
				<div class="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
					{generalError}
				</div>
			{/if}

			<DynamicForm
				config={formConfig}
				{formData}
				{errors}
				{loading}
				onchange={handleChange}
				onsubmit={handleSubmit}
			/>
		</div>
	</div>
</div>
