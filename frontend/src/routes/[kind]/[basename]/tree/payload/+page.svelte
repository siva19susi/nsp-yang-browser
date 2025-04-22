<script lang="ts">
  import { onMount } from "svelte"

  import Navbar from "$lib/components/Navbar.svelte"
  import Loading from "$lib/components/Loading.svelte"
	import ErrorNotification from "$lib/components/ErrorNotification.svelte"

  import type { YangTreePayloadPostMessage, YangTreePayloadResponseMessage } from "./structure"
	import Footer from "$lib/components/Footer.svelte";

  // DEFAULTS
  let treePayload: any
  let workerComplete = false
  let workerStatus = {status: 404, error: {message: "Unknown Error"}}

  // YANGTREE WORKER
  let yangTreePayloadWorker: Worker | undefined = undefined
  async function loadYangTreePayloadWorker (data: YangTreePayloadPostMessage) {
    const YangTreePayloadWorker = await import('./yangTreePayload.worker?worker')
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
  let {kind, basename, urlPath, nspIp} = data

  // OTHER BINDING VARIABLES
  onMount(() => loadYangTreePayloadWorker(data))
</script>

<svelte:head>
	<title>NSP YANG Browser | Payload - {basename} ({kind})</title>
</svelte:head>

{#if !workerComplete}
  <Loading/>
{:else}
  {#if workerStatus.status === 200}
    <Navbar {kind} {basename} {nspIp}/>
    <div class="font-nunito">
      <div class="px-6 py-4 text-sm dark:text-white pt-[85px]">
        <p class="pb-1 font-semibold text-black dark:text-white">Target Path:</p>
        <pre>{urlPath}</pre>
        <p class="pt-4 pb-1 font-semibold text-black dark:text-white">Sample Payload:</p>
        <pre>{JSON.stringify(treePayload, null, 2)}</pre>
      </div>
      <Footer home={false}/>
    </div>
  {:else}
    <ErrorNotification pageError={workerStatus} />
  {/if}
{/if}
