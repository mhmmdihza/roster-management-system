<script lang="ts">
	import { userStore } from '$lib/auth';
	import Sidebar from '$lib/layout/Sidebar.svelte';
	import Topbar from '$lib/layout/Topbar.svelte';
	let { data, children } = $props();

	import '../../app.css';

	if (data.user) {
		userStore.set(data.user);
	}
</script>

<div class="flex h-screen">
	<Sidebar role={data.user.role} />

	<div class="flex flex-1 flex-col">
		<Topbar {data} />
		<main class="p-4">
			{@render children()}
		</main>
	</div>
</div>
