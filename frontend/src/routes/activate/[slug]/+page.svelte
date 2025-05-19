<script lang="ts">
	import { page } from '$app/state';
	import { activateAccount } from '$lib/api';

	export let data: { id: string };

	let name = '';
	let password = '';
	let showModal = false;
	let message = '';
	let error = '';

	async function handleSubmit() {
		error = '';
		message = '';

		try {
			const res = await activateAccount({
				id: data.id,
				name,
				password
			});

			if (!res.ok) {
				error = 'Activation failed, make sure registration link are valid';
				return;
			}

			const resData = await res.json();
			message = resData.message || 'Account activated successfully';
			showModal = true;
		} catch (err) {
			error = 'Unexpected error: ' + err;
		}
	}
	function closeAndRedirect() {
		showModal = false;
		window.location.href = page.url.origin;
	}
</script>

<section class="mx-auto mt-10 max-w-xl rounded-2xl bg-white p-6 shadow-md">
	<h1 class="mb-6 text-2xl font-bold text-gray-800">Activate Account</h1>

	<p class="mb-4 text-gray-600">Activating ID: <code>{data.id}</code></p>

	<form on:submit|preventDefault={handleSubmit} class="space-y-4">
		<div>
			<label for="name" class="block text-sm font-medium text-gray-700">Name</label>
			<input
				id="name"
				type="text"
				bind:value={name}
				required
				class="mt-1 w-full rounded-md border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
			/>
		</div>

		<div>
			<label for="password" class="block text-sm font-medium text-gray-700">Password</label>
			<input
				id="password"
				type="password"
				bind:value={password}
				required
				minlength="6"
				class="mt-1 w-full rounded-md border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
			/>
		</div>

		{#if error}
			<p class="text-sm text-red-600">{error}</p>
		{/if}

		<button
			type="submit"
			class="w-full rounded-md bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700"
		>
			Activate Account
		</button>
	</form>
</section>

{#if showModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
		<div class="w-full max-w-md rounded-lg bg-white p-6 text-center shadow-lg">
			<h2 class="mb-4 text-lg font-semibold text-gray-800">Success!</h2>
			<p class="mb-4 text-gray-700">{message}</p>
			<button
				on:click={closeAndRedirect}
				class="mx-auto mt-4 block text-sm text-gray-500 hover:underline"
			>
				Close
			</button>
		</div>
	</div>
{/if}
