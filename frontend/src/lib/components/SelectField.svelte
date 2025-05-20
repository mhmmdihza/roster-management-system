<script lang="ts">
	export let id: string;
	export let value: string | number | null;
	export let options: { id: string | number; label: string }[] = [];
	export let required: boolean = false;
	export let errorValidation: string = '';
	export let disabled: boolean | undefined | null;
	export let validate: (() => string) | undefined;

	let touched = false;
</script>

<div class="mb-4">
	<select
		{id}
		bind:value
		{required}
		{disabled}
		class="mt-1 w-full rounded-md border px-4 py-2 focus:outline-none focus:ring-2
			{errorValidation && touched
			? 'border-red-500 focus:ring-red-500'
			: 'border-gray-300 focus:ring-blue-500'}"
		on:change={() => {
			touched = true;
			if (validate) {
				errorValidation = validate();
			}
		}}
	>
		{#each options as option}
			<option value={option.id}>{option.label}</option>
		{/each}
	</select>

	{#if touched && errorValidation}
		<p class="mt-1 text-sm text-red-600">{errorValidation}</p>
	{/if}
</div>
