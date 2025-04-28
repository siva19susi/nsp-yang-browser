<script lang="ts">
  import { onMount } from "svelte"

  import Navbar from "$lib/components/Navbar.svelte"

  import type { IntentQueryPostMessage, IntentQueryResponseMessage } from "./structure"
	import Footer from "$lib/components/Footer.svelte";

  // DEFAULTS
  let isSubmitting = false
  let trimUrlListPath = false
  let nspResponse: any = {}

  let workerStatus = {
    complete: false, 
    success: false,
    progress: 0,
    error: {
      message: "Unknown Error"
    }
  }

  // INVENTORY FIND WORKER
  let inventoryFindWorker: Worker | undefined = undefined
  async function loadInventoryFindWorker (data: IntentQueryPostMessage) {
    workerStatus.progress = 30
    const InventoryFindWorker = await import('./intentQuery.worker?worker')
    inventoryFindWorker = new InventoryFindWorker.default()
    inventoryFindWorker.postMessage(data)
    inventoryFindWorker.onmessage = onInventoryFindWorkerMessage
  }
  function onInventoryFindWorkerMessage(event: MessageEvent<IntentQueryResponseMessage>) {
    const response = event.data

    if (response.type === "progress") {
      workerStatus.progress = response.value
    }
    if (response.type === "complete") {
      workerStatus.progress = 100
      workerStatus.error.message = response.message
      workerStatus.complete = true
      if(response.success) {
        nspResponse = response.output
      }
    }
  }

  // ON PAGELOAD
	export let data
  let {kind, basename, urlPath, nspIp} = data
  let findPayload = {
    "xpath-filter": urlPath,
    "include-meta": false,
    fields: "",
    depth: 2,
    limit: 1,
    offset: 1
  }

  async function inventoryFind(event: SubmitEvent) {
    isSubmitting = false
    workerStatus.complete = false
    const formData = new FormData(event.currentTarget as HTMLFormElement)
		const findPayload = {
      url: formData.get("xpath-filter") as string,
      target: formData.get("target") as string,
      "intent-key": basename.split("_")[0]
    }

    loadInventoryFindWorker(findPayload)
  }

  
  function trimBeforeFirstBracketSegment (path: string) {
    // Split path into segments
    const segments = path.split("/")

    // Find the index of the first segment that contains [=]
    const index = segments.findIndex(seg => /\[.*=.*\]/.test(seg))

    // If such a segment exists and there's at least one segment before it
    if (index > 0) {
      return segments.slice(0, index).join("/")
    }

    // If no filter or it's the first segment, return as-is
    return path
  }
</script>

<svelte:head>
	<title>NSP YANG Browser | Query - {basename} ({kind})</title>
</svelte:head>


<Navbar {kind} {basename} {nspIp}/>
<div class="px-6 py-4 font-nunito container mx-auto pt-[85px]">
  <p class="text-gray-800 dark:text-gray-300 pb-1">Connected NSP intent explorer:</p>
  <form class="pt-4" method="POST" action="?/inventoryFind" on:submit|preventDefault={inventoryFind}>
    <div class="space-y-4">
      <div>
        <label for="xpath-filter" class="block uppercase text-gray-800 dark:text-gray-200 text-xs mb-2">Intent Path*</label>
        <input id="xpath-filter" name="xpath-filter" type="text" required value="{trimUrlListPath ? trimBeforeFirstBracketSegment(urlPath) : urlPath}" class="font-fira px-3 py-2 rounded-lg w-full text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
        {#if urlPath.includes("[")}
          <div class="px-1 py-1.5 flex items-center">
            <input id="default-checkbox" type="checkbox" class="w-3 h-3" bind:checked={trimUrlListPath}>
            <label for="default-checkbox" class="ms-2 text-xs text-nowrap text-gray-900 dark:text-gray-300 cursor-pointer">Trim the path if the first list key is unknown</label>
          </div>
        {/if}
      </div>
      <div>
        <label for="target" class="block uppercase text-gray-800 dark:text-gray-200 text-xs mb-2">Intent Target</label>
        <input id="target" name="target" type="text" value="" class="font-fira px-3 py-2 rounded-lg w-full text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
      </div>
    </div>
    <div class="flex items-center justify-end pt-4">
      <button type="submit" class="px-4 py-2 rounded-lg text-sm text-white bg-green-600 hover:bg-green-700 {isSubmitting ? 'animate-pulse' : ''}" disabled={isSubmitting}>{isSubmitting ? 'Submitting...' : 'Submit'}</button>
    </div>
  </form>
  <div class="text-sm dark:text-white">
    {#if workerStatus.progress > 0}
      {#if !workerStatus.complete}
        <div class="px-1 py-2">
          <div class="w-full bg-gray-200 rounded-full h-1 dark:bg-gray-700">
            <div class="bg-blue-600 h-1 rounded-full" style="width: {workerStatus.progress}%"></div>
          </div>
          <p class="pt-2 text-black dark:text-white">Quering NSP...</p>
        </div>
      {:else}
        <p class="pt-4 pb-1 font-semibold text-black dark:text-white">NSP Response:</p>
        {#if Object.keys(nspResponse).length === 0}
          <p class="py-2 text-red-600">{workerStatus.error.message}</p>
        {:else}
          <div class="overflow-x-auto scroll-light dark:scroll-dark">
            <pre>{JSON.stringify(nspResponse, null, 2)}</pre>
          </div>
        {/if}
      {/if}
    {/if}
  </div>
  <Footer home={false}/>
</div>
