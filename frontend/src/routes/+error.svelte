<script lang="ts">
	import { page } from '$app/stores'
	import ErrorNotification from '$lib/components/ErrorNotification.svelte';
</script>

<ErrorNotification pageError={$page} />