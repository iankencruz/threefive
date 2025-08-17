<script>
	let fname = $state();
	let lname = $state();
	let email = $state();
	let subject = $state();
	let msgbody = $state();
</script>

<div class="mx-auto mt-32 w-full max-w-7xl px-4 sm:px-6 lg:px-8">
	<div class=" flex max-w-3xl flex-col">
		<h2 class="text-3xl font-extrabold text-gray-100">Contact</h2>

		<div class="mt-8 flex w-full flex-col items-center">
			<div class="w-full">
				<form id="contact-form">
					<fieldset>
						<legend class="block text-left text-sm font-medium text-gray-100">
							Name
							<span class="text-xs font-light text-red-400">* Required</span>
						</legend>

						<div class="grid grid-cols-12 items-end gap-x-2">
							<div class="col-span-6">
								<label for="first_name" class="block text-left text-sm font-light text-gray-100">
									First Name
								</label>
								<input
									type="text"
									bind:value={fname}
									name="first_name"
									class="mt-2 block w-full border border-gray-300 px-4 py-2 text-lg text-gray-900 focus:border-gray-900 focus:ring-gray-900"
									required
								/>
							</div>
							<div class="col-span-6">
								<label for="last_name" class="block text-left text-xs font-light text-gray-100"
									>Last Name</label
								>
								<input
									bind:value={lname}
									type="text"
									name="last_name"
									class="mt-2 block w-full border border-gray-300 px-4 py-2 text-lg text-gray-900 focus:border-gray-900 focus:ring-gray-900"
									required
								/>
							</div>
						</div>
						<div></div>
					</fieldset>

					<div class="my-4">
						<label for="email" class="block text-left text-sm font-medium text-gray-100">
							Email
							<span class="text-xs font-light text-red-400">* Required</span>
						</label>
						<input
							type="email"
							bind:value={email}
							name="email"
							class="mt-2 block w-full border border-gray-300 px-4 py-2 text-lg text-gray-900 focus:border-gray-900 focus:ring-gray-900"
							required
						/>
					</div>
					<div class="my-4">
						<label for="subject" class="block text-left text-sm font-medium text-gray-100">
							Subject
							<span class="text-xs font-light text-red-400">* Required</span>
						</label>
						<input
							bind:value={subject}
							type="text"
							name="subject"
							class="mt-2 block w-full border border-gray-300 px-4 py-2 text-lg text-gray-900 focus:border-gray-900 focus:ring-gray-900"
							required
						/>
					</div>
					<div class="my-4">
						<label for="message" class="block text-left text-sm font-medium text-gray-100">
							Message
							<span class="text-xs font-light text-red-400">* Required</span>
						</label>
						<textarea
							bind:value={msgbody}
							name="message"
							rows="5"
							class="mt-2 block w-full border border-gray-300 px-4 py-2 text-lg text-gray-900 focus:border-gray-900 focus:ring-gray-900"
							required
						></textarea>
					</div>
					<!-- 
          Hidden field where we store the v3 token.
          On the server, read it via r.FormValue("g-recaptcha-response"). 
        -->
					<input type="hidden" name="g-recaptcha-response" id="g-recaptcha-response" />

					<!-- 
          Type=\"button\" so it doesn't do a normal submit. 
          We'll call grecaptcha.execute() -> dispatch the \"verified\" event => triggers HTMX 
        -->
					<button
						type="button"
						onclick={() => alert('recaptcha run')}
						class="cursor-pointer px-3.5 py-2.5 text-sm font-semibold text-white shadow-xs outline hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
					>
						Send Message
					</button>
				</form>
			</div>
		</div>
	</div>
</div>
