<script lang="ts">
	export let role: string;
	let links =
		role === 'admin'
			? [{ name: 'Register New User', href: '/admin/register' }]
			: [{ name: 'shifts', href: '/employee/shifts' }];
</script>

<aside class="h-full w-64 bg-gray-800 p-4 text-white">
	<h2 class="mb-4 text-xl font-bold capitalize">{role} Panel</h2>
	<nav class="flex flex-col gap-2">
		{#each links as link}
			<a href={link.href} class="rounded px-2 py-1 hover:bg-gray-700">{link.name}</a>
		{/each}
	</nav>
</aside>
