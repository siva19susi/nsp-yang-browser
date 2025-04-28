<script lang="ts">
  import { onMount } from "svelte"

  import Navbar from "$lib/components/Navbar.svelte"

  import type { InventoryFindPostMessage, InventoryFindResponseMessage } from "./structure"
	import Footer from "$lib/components/Footer.svelte";

  // DEFAULTS
  let isSubmitting = false
  let nspResponse: any = {}

  let workerStatus = {
    complete: false, 
    success: false,
    progress: 30,
    error: {
      message: "Unknown Error"
    }
  }

  // INVENTORY FIND WORKER
  let inventoryFindWorker: Worker | undefined = undefined
  async function loadInventoryFindWorker (data: InventoryFindPostMessage) {
    const InventoryFindWorker = await import('./inventoryFind.worker?worker')
    inventoryFindWorker = new InventoryFindWorker.default()
    inventoryFindWorker.postMessage(data)
    inventoryFindWorker.onmessage = onInventoryFindWorkerMessage
  }
  function onInventoryFindWorkerMessage(event: MessageEvent<InventoryFindResponseMessage>) {
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

  onMount(() => loadInventoryFindWorker({kind, nsp: findPayload}))

  async function inventoryFind(event: SubmitEvent) {
    isSubmitting = false
    workerStatus.complete = false
    const formData = new FormData(event.currentTarget as HTMLFormElement)
		let findPayload: any = {}
    formData.forEach(function(value, key) {
      if(key === "include-meta") {
        findPayload[key] = (value == "true" ? true : false)
      } else if(key !== "xpath-filter" && key !== "fields") {
        findPayload[key] = parseInt(formData.get(key) as string)
      } else {
        findPayload[key] = value
      }
    })
    loadInventoryFindWorker({kind, nsp: findPayload})
  }
</script>

<svelte:head>
	<title>NSP YANG Browser | Query - {basename} ({kind})</title>
</svelte:head>


<Navbar {kind} {basename} {nspIp}/>
<div class="px-6 py-4 font-nunito container mx-auto pt-[85px]">
  <p class="text-gray-800 dark:text-gray-300 pb-1">Connected NSP inventory find:</p>
  <form class="pt-4" method="POST" action="?/inventoryFind" on:submit|preventDefault={inventoryFind}>
    <div class="space-y-4">
      <div>
        <label for="xpath-filter" class="block uppercase text-gray-800 dark:text-gray-200 text-xs mb-2">Target Path*</label>
        <input id="xpath-filter" name="xpath-filter" type="text" required value="{urlPath}" class="font-fira px-3 py-2 rounded-lg w-full text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
      </div>
      <div>
        <label for="fields" class="block uppercase text-gray-800 dark:text-gray-200 text-xs mb-2">Field Filter</label>
        <input id="fields" name="fields" type="text" value="" class="font-fira px-3 py-2 rounded-lg w-full text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
      </div>
    </div>
    <div class="grid md:grid-cols-4 grid-cols-2 gap-4 pt-4">
      <div>
        <label for="include-meta" class="block uppercase text-gray-800 dark:text-gray-200 text-xs mb-2">include-meta</label>
        <select id="include-meta" name="include-meta" class="px-3 py-2 rounded-lg w-full text-sm border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
          <option selected value="false">false</option>
          <option value="true">true</option>
        </select>
      </div>
      <div>
        <label for="depth" class="block uppercase text-gray-800 dark:text-gray-200 text-xs mb-2">Depth*</label>
        <input id="depth" name="depth" type="number" value="2" required class="px-3 py-2 rounded-lg w-full text-sm border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
      </div>
      <div>
        <label for="limit" class="block uppercase text-gray-800 dark:text-gray-200 text-xs mb-2">Limit*</label>
        <input id="limit" name="limit" type="number" value="1" required class="px-3 py-2 rounded-lg w-full text-sm border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
      </div>
      <div>
        <label for="offset" class="block uppercase text-gray-800 dark:text-gray-200 text-xs mb-2">Page / Offset*</label>
        <input id="offset" name="offset" type="number" value="1" required class="px-3 py-2 rounded-lg w-full text-sm border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
      </div>
    </div>
    <div class="flex items-center justify-end pt-4">
      <button type="submit" class="px-4 py-2 rounded-lg text-sm text-white bg-green-600 hover:bg-green-700 {isSubmitting ? 'animate-pulse' : ''}" disabled={isSubmitting}>{isSubmitting ? 'Submitting...' : 'Submit'}</button>
    </div>
  </form>
  <div class="text-sm dark:text-white">
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
  </div>
  <Footer home={false}/>
</div>
