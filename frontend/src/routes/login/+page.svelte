<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { login } from '$lib/api';
	import ButtonSubmit from '$lib/components/ButtonSubmit.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import { z } from 'zod';

	let username = '';
	let usernameErrValidation = 'x'; // this non empty value will disable the button at init
	const usernameSchema = z.string().email('Invalid email format').min(5, 'Too short');
	let password = '';
	let passwordErrValidation = 'x';
	const passwordSchema = z.string().min(6, 'Too short');

	let error = '';
	let loading = false;

	async function handleLogin() {
		error = '';
		loading = true;
		try {
			const res = await login(username, password);
			const data = await res.json();

			if (res.ok) {
				await invalidateAll();
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

		<TextInput
			id="username"
			placeholder="Username"
			type="email"
			bind:value={username}
			schema={usernameSchema}
			disabled={loading}
			bind:errorValidation={usernameErrValidation}
		/>

		<TextInput
			id="password"
			placeholder="Password"
			type="password"
			bind:value={password}
			schema={passwordSchema}
			disabled={loading}
			bind:errorValidation={passwordErrValidation}
		/>
		<ButtonSubmit
			{loading}
			text="Login"
			textOnLoading="Logging in..."
			disabled={usernameErrValidation !== '' || passwordErrValidation !== ''}
		/>
	</form>
</div>
