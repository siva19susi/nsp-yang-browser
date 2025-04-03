<script lang="ts">
  import { onMount } from "svelte"

  import { kindView } from "$lib/components/functions"
  import Loading from "$lib/components/Loading.svelte"
	import ErrorNotification from "$lib/components/ErrorNotification.svelte"

  import type { YangTreePayloadPostMessage, YangTreePayloadResponseMessage } from "$lib/workers/structure"

  // DEFAULTS
  let treePayload: any
  let workerComplete = false
  let workerStatus = {status: 404, error: {message: "Unknown Error"}}

  // YANGTREE WORKER
  let yangTreePayloadWorker: Worker | undefined = undefined
  async function loadYangTreePayloadWorker (data: YangTreePayloadPostMessage) {
    const YangTreePayloadWorker = await import('$lib/workers/yangTreePayload.worker?worker')
    yangTreePayloadWorker = new YangTreePayloadWorker.default()
    yangTreePayloadWorker.postMessage(data)
    yangTreePayloadWorker.onmessage = onYangTreePayloadWorkerMessage
  }
  function onYangTreePayloadWorkerMessage(event: MessageEvent<YangTreePayloadResponseMessage>) {
    const response = event.data
    workerStatus.error.message = response.message
    if(event.data.success) {
      treePayload = response.tree
      workerStatus.status = 200
    }
    workerComplete = true
  }

  // ON PAGELOAD
	export let data
  let {kind, basename} = data

  // OTHER BINDING VARIABLES
  onMount(() => loadYangTreePayloadWorker(data))
</script>

<svelte:head>
	<title>Yang Tree Browser {basename} ({kindView(kind)}) Payload</title>
</svelte:head>

{#if !workerComplete}
  <Loading/>
{:else}
  {#if workerStatus.status === 200}
    <pre class="text-sm dark:text-white">{JSON.stringify(treePayload, null, 2)}</pre>
  {:else}
    <ErrorNotification pageError={workerStatus} />
  {/if}
{/if}
