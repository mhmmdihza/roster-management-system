<script lang="ts">
	export let loading: boolean;
	export let disabled: boolean;
	export let text: string;
	export let textOnLoading: string;
</script>

<button
	class="flex w-full items-center justify-center gap-2 rounded bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-blue-400"
	type="submit"
	disabled={loading || disabled}
>
	{#if loading}
		<!-- Spinner -->
		<svg class="h-5 w-5 animate-spin text-white" fill="none" viewBox="0 0 24 24">
			<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
			></circle>
			<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
		</svg>
		<span>{textOnLoading}</span>
	{:else}
		{text}
	{/if}
</button>
