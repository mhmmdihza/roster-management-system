<script lang="ts">
	import { goto } from '$app/navigation';
	import { login } from '$lib/api';

	let username = '';
	let password = '';
	let error = '';
	let loading = false;

	async function handleLogin() {
		error = '';
		loading = true;
		try {
			const res = await login(username, password);
			const data = await res.json();

			if (res.ok) {
				goto('/');
			} else {
				error = data?.error || data?.message || 'Login failed.';
			}
		} catch (e) {
			error = 'Login failed. Please try again later.';
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-100 px-4">
	<form
		on:submit|preventDefault={handleLogin}
		class="w-full max-w-sm rounded-2xl bg-white p-8 shadow-lg transition-all"
	>
		<h2 class="mb-6 text-center text-3xl font-bold text-gray-800">Role Management</h2>

		{#if error}
			<div class="mb-4 rounded bg-red-100 p-2 text-sm text-red-600">
				{error}
			</div>
		{/if}

		<input
			class="mb-4 w-full rounded border border-gray-300 p-3 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-400"
			placeholder="Username"
			bind:value={username}
			disabled={loading}
		/>

		<input
			class="mb-6 w-full rounded border border-gray-300 p-3 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-400"
			type="password"
			placeholder="Password"
			bind:value={password}
			disabled={loading}
		/>

		<button
			class="flex w-full items-center justify-center gap-2 rounded bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-blue-400"
			type="submit"
			disabled={loading}
		>
			{#if loading}
				<!-- Spinner -->
				<svg class="h-5 w-5 animate-spin text-white" fill="none" viewBox="0 0 24 24">
					<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
					></circle>
					<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
					></path>
				</svg>
				<span>Logging in...</span>
			{:else}
				Login
			{/if}
		</button>
	</form>
</div>
