<script lang="ts">
	import { z } from 'zod';

	export let id: string;
	export let placeholder: string;
	export let type: string = 'text';
	export let value: string;
	export let required: boolean = false;
	export let errorValidation: string = '';
	export let disabled: boolean | undefined | null;

	export let schema: z.ZodType<any> | null = null;

	let touched = false;

	// validate on input change or when field is touched
	function validate() {
		if (!schema) {
			errorValidation = '';
			return;
		}

		const result = schema.safeParse(value);

		if (!result.success) {
			errorValidation = result.error.errors[0]?.message ?? 'Invalid input';
		} else {
			errorValidation = '';
		}
	}
</script>

<div class="mb-4">
	<input
		{placeholder}
		on:input={(e) => {
			touched = true;
			value = (e.target as HTMLInputElement).value;
			validate();
		}}
		{disabled}
		{id}
		{type}
		bind:value
		{required}
		class="mt-1 w-full rounded-md border px-4 py-2 focus:outline-none focus:ring-2
			{errorValidation && touched
			? 'border-red-500 focus:ring-red-500'
			: 'border-gray-300 focus:ring-blue-500'}"
	/>
	{#if touched && errorValidation}
		<p class="mt-1 text-sm text-red-600">{errorValidation}</p>
	{/if}
</div>
