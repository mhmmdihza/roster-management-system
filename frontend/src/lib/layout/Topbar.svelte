<script lang="ts">
	import { goto } from '$app/navigation';
	import { logout } from '$lib/api';
	export let data;

	async function handleLogout() {
		await logout();
		goto('/login');
	}
</script>

<div class="flex items-center justify-end gap-4 bg-white p-4 shadow">
	<span class="text-gray-700">{data.user.employee_name}</span>
	<button on:click={handleLogout} class="text-red-500 hover:underline">Logout</button>
</div>
