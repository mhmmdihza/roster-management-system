<script lang="ts">
	import { page } from '$app/state';
	import type { RoleResponse } from '$lib/api';
	import { listRoles, registerUser } from '$lib/api';
	import { onMount } from 'svelte';

	let email = '';
	let roleAdmin = false;
	let primaryRole: number | null = null;

	let roles: RoleResponse[] = [];
	let showModal = false;
	let activationUrl = '';
	let copied = false;

	onMount(async () => {
		try {
			const allRoles = await listRoles();
			roles = allRoles.filter((role) => role.id !== 0);
		} catch (err) {
			console.error('Failed to load roles:', err);
		}
	});

	async function handleSubmit() {
		try {
			const res = await registerUser(email, primaryRole, roleAdmin);
			if (!res.ok) {
				const error = await res.json();
				alert(`Registration failed: ${error.error}`);
				return;
			}
			const data = await res.json();
			const currentUrl = page.url.origin;
			activationUrl = `${currentUrl}/activate/${data.id}`;
			showModal = true;
		} catch (err) {
			alert('Unexpected error: ' + err);
		}
	}

	function copyToClipboard() {
		navigator.clipboard.writeText(activationUrl).then(() => {
			copied = true;
			setTimeout(() => (copied = false), 2000);
		});
	}
</script>

<section class="mx-auto mt-10 max-w-xl rounded-2xl bg-white p-6 shadow-md">
	<h1 class="mb-6 text-2xl font-bold text-gray-800">Register New User</h1>

	<form on:submit|preventDefault={handleSubmit} class="space-y-4">
		<div>
			<label for="email" class="block text-sm font-medium text-gray-700">Email</label>
			<input
				id="email"
				type="email"
				bind:value={email}
				required
				class="mt-1 w-full rounded-md border border-gray-300 px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
			/>
		</div>

		<div class="flex items-center space-x-2">
			<input
				id="admin"
				type="checkbox"
				bind:checked={roleAdmin}
				class="h-4 w-4 rounded border-gray-300 text-blue-600"
			/>
			<label for="admin" class="text-sm text-gray-700">Is Admin?</label>
		</div>

		{#if !roleAdmin}
			<div>
				<label for="primaryRole" class="block text-sm font-medium text-gray-700">Primary Role</label
				>
				<select
					id="primaryRole"
					bind:value={primaryRole}
					class="mt-1 w-full rounded-md border border-gray-300 bg-white px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
				>
					<option value="" disabled selected>Select a role</option>
					{#each roles as role}
						<option value={role.id}>{role.roleName}</option>
					{/each}
				</select>
			</div>
		{/if}

		<button
			type="submit"
			class="w-full rounded-md bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700"
		>
			Register
		</button>
	</form>
</section>

{#if showModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
		<div class="w-full max-w-md rounded-lg bg-white p-6 text-center shadow-lg">
			<h2 class="mb-4 text-lg font-semibold text-gray-800">User Created Successfully!</h2>
			<p class="mb-4 break-all text-gray-700">{activationUrl}</p>
			<button
				on:click={copyToClipboard}
				class="rounded bg-green-600 px-4 py-2 text-white transition hover:bg-green-700"
			>
				{#if copied}Copied!{/if}
				{#if !copied}Copy Link{/if}
			</button>
			<button
				on:click={() => (showModal = false)}
				class="mx-auto mt-4 block text-sm text-gray-500 hover:underline"
			>
				Close
			</button>
		</div>
	</div>
{/if}
