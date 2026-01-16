<!-- frontend/src/routes/(app)/(public)/contact/+page.svelte -->
<script lang="ts">
	import { enhance } from '$app/forms';
	import { PUBLIC_API_URL } from '$env/static/public';
	import DynamicForm from '$components/ui/form/DynamicForm.svelte';
	import type { FormConfig } from '$components/ui/form';
	import { Mail, MapPin, Phone } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let formData = $state({
		name: '',
		email: '',
		subject: '',
		message: ''
	});

	let isSubmitting = $state(false);
	let submitSuccess = $state(false);
	let submitError = $state('');
	let formErrors = $state<Record<string, string>>({});

	const formConfig: FormConfig = {
		fields: [
			{
				name: 'name',
				type: 'text',
				label: 'Name',
				placeholder: 'Your name',
				required: true,
				colSpan: 6
			},
			{
				name: 'email',
				type: 'email',
				label: 'Email',
				placeholder: 'your@email.com',
				required: true,
				colSpan: 6
			},
			{
				name: 'subject',
				type: 'text',
				label: 'Subject',
				placeholder: 'How can we help?',
				colSpan: 12
			},
			{
				name: 'message',
				type: 'textarea',
				label: 'Message',
				placeholder: 'Tell us more about your inquiry...',
				required: true,
				colSpan: 12,
				rows: 6
			}
		],
		submitText: 'Send Message',
		submitVariant: 'primary',
		submitFullWidth: true
	};

	// Helper to capitalise
	function capitalize(s: string) {
		return s.charAt(0).toUpperCase() + s.slice(1);
	}

	async function handleSubmit(data: Record<string, any>) {
		isSubmitting = true;
		submitError = '';
		formErrors = {};
		submitSuccess = false;

		try {
			const response = await fetch(`${PUBLIC_API_URL}/api/v1/contact`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(data)
			});

			const result = await response.json();

			if (!response.ok) {
				if (result.errors && Array.isArray(result.errors)) {
					const newErrors: Record<string, string> = {};

					const errorMessages = result.errors.map((err: { field: string; message: string }) => {
						// Store raw error for the red borders in DynamicForm
						newErrors[err.field] = err.message;
						// Return capitalized string for the summary/toast
						return `${capitalize(err.field)}: ${err.message}`;
					});

					formErrors = newErrors;
					submitError = errorMessages.join('\n');

					errorMessages.forEach((msg: string) => toast.error(msg));
					throw new Error('Validation failed');
				}

				if (response.status === 429) {
					throw new Error('Too many requests. Please try again later.');
				}

				throw new Error(result.message || 'Failed to send message');
			}

			submitSuccess = true;
			toast.success('Your message has been sent successfully!');
			formData = { name: '', email: '', subject: '', message: '' };
		} catch (error) {
			console.log('error:', error);
			// Only set submitError if it hasn't been set by the validation logic above
			if (!submitError) {
				submitError = error instanceof Error ? error.message : 'Failed to send message';
				toast.error(submitError);
			}
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>TFP - Contact Us</title>
	<meta name="description" content="Get in touch with us. We'd love to hear from you." />
</svelte:head>

<div class="min-h-screen bg-background py-16">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<!-- Header -->
		<div class="mb-12 text-center">
			<h1 class="mb-4">Get in Touch</h1>
			<p class="mx-auto max-w-2xl text-lg text-gray-400">
				Have a question or want to work together? We'd love to hear from you.
			</p>
		</div>

		<div class="grid gap-12 lg:grid-cols-3">
			<!-- Contact Information -->
			<div class="lg:col-span-1">
				<h2 class="mb-6 text-xl font-semibold">Contact Information</h2>

				<div class="space-y-6">
					<div class="flex items-start gap-4">
						<Mail class="mt-1 h-5 w-5 text-primary" />
						<div>
							<h3 class="mb-1 font-medium">Email</h3>
							<a href="mailto:hello@example.com" class="text-gray-400 hover:text-primary">
								hello@example.com
							</a>
						</div>
					</div>

					<div class="flex items-start gap-4">
						<Phone class="mt-1 h-5 w-5 text-primary" />
						<div>
							<h3 class="mb-1 font-medium">Phone</h3>
							<a href="tel:+1234567890" class="text-gray-400 hover:text-primary">
								+1 (234) 567-890
							</a>
						</div>
					</div>

					<div class="flex items-start gap-4">
						<MapPin class="mt-1 h-5 w-5 text-primary" />
						<div>
							<h3 class="mb-1 font-medium">Location</h3>
							<p class="text-gray-400">
								123 Main Street<br />
								Suite 100<br />
								City, State 12345
							</p>
						</div>
					</div>
				</div>
			</div>

			<!-- Contact Form -->
			<div class="lg:col-span-2">
				<div class="rounded-lg bg-surface p-8 shadow-lg">
					<DynamicForm config={formConfig} {formData} onSubmit={handleSubmit} asForm={true} />
				</div>
			</div>
		</div>
	</div>
</div>
