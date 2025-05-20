<script lang="ts">
	import { page } from '$app/state';
	import type { RoleResponse } from '$lib/api';
	import { registerUser } from '$lib/api';
	import ButtonSubmit from '$lib/components/ButtonSubmit.svelte';
	import SelectField from '$lib/components/SelectField.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import { z } from 'zod';
	import type { PageData } from './$types';

	export let data: PageData;

	let email = '';
	let emailErrValidation = 'x'; // disable button on init
	const emailSchema = z.string().email('Invalid email format').min(5, 'Too short');
	let roleAdmin = false;
	let primaryRole: number | null = null;
	let roleErrValidation = 'x';
	let roles: RoleResponse[] = data.roles;

	const roleSchema = z
		.object({
			roleAdmin: z.boolean(),
			primaryRole: z.number().nullable()
		})
		.superRefine((data, ctx) => {
			if (!data.roleAdmin && (!data.primaryRole || data.primaryRole < 1)) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					message: 'Role is required when is admin not checked',
					path: ['primaryRole']
				});
			}
		});

	function validateRole() {
		const result = roleSchema.safeParse({
			roleAdmin: roleAdmin,
			primaryRole: primaryRole
		});

		if (!result.success) {
			return result.error.errors[0]?.message ?? 'Invalid option';
		}
		return '';
	}

	let showModal = false;
	let activationUrl = '';
	let copied = false;
	let loading = false;

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
			<TextInput
				id="email"
				placeholder="Email"
				type="email"
				bind:value={email}
				schema={emailSchema}
				disabled={loading}
				bind:errorValidation={emailErrValidation}
			/>
		</div>

		<div class="flex items-center space-x-2">
			<input
				id="admin"
				type="checkbox"
				bind:checked={roleAdmin}
				on:change={validateRole}
				class="h-4 w-4 rounded border-gray-300 text-blue-600"
				disabled={loading}
			/>
			<label for="admin" class="text-sm text-gray-700">Is Admin?</label>
		</div>

		{#if !roleAdmin}
			<label for="primaryRole" class="block text-sm font-medium text-gray-700">Primary Role</label>
			<SelectField
				id="primaryRole"
				bind:value={primaryRole}
				options={roles.map((r) => ({ id: r.id, label: r.roleName }))}
				disabled={loading}
				validate={validateRole}
				bind:errorValidation={roleErrValidation}
			/>
		{/if}

		<ButtonSubmit
			{loading}
			text="Register"
			textOnLoading="processing..."
			disabled={emailErrValidation !== '' || roleErrValidation !== ''}
		/>
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
