
<script lang="ts">
  import { onMount } from 'svelte'

	import Navbar from '$lib/components/Navbar.svelte'
  import Footer from '$lib/components/Footer.svelte'
	import Popup from '$lib/components/Popup.svelte'
  import Loading from '$lib/components/Loading.svelte'
	import NothingToCompare from '$lib/components/NothingToCompare.svelte'
  import ErrorNotification from '$lib/components/ErrorNotification.svelte'

  import SearchInput from '$lib/components/SearchInput.svelte'
  import StateButton from '$lib/components/StateButton.svelte'
  import WithDefaultCheck from '$lib/components/WithDefaultCheck.svelte'
	import ChangesButton from '$lib/components/ChangesButton.svelte'
  import Pagination from './Pagination.svelte'
	
  import type { ComparePayLoad } from '$lib/structure'
  import type { CompareResponseMessage, DiffResponseMessage } from './structure'
  import { compareStore, defaultStore, paginated, searchStore, stateStore, total, yangPaths } from './store'
  import { markFilter, markRender, toLower } from '$lib/components/functions'
	
  // DEFAULTS
  let popupDetail = {}
  let diff: DiffResponseMessage[] = []
  let workerComplete = false
  let workerStatus = {status: 404, error: {message: "Unknown Error"}}

  // COMPARE WORKER
  let compareWorker: Worker | undefined = undefined
  async function loadWorker(data: ComparePayLoad) {
    const CompareWorker = await import('./compare.worker?worker')
    compareWorker = new CompareWorker.default()
    compareWorker.postMessage(data)
    compareWorker.onmessage = onWorkerMessage
  }
  function onWorkerMessage(event: MessageEvent<CompareResponseMessage>) {
    const response = event.data
    workerStatus.error.message = response.message
    if(event.data.success) {
      diff = response.diff
      workerStatus.status = 200
    }
    workerComplete = true
  }

  // ON PAGELOAD
  export let data: ComparePayLoad
  const {x, y, xKind, xBasename, yKind, yBasename, urlPath} = data
  onMount(() => loadWorker(data))

  // OTHER BINDING VARIABLES
  let searchInput = urlPath
  let compareInput = ""
  let stateInput: string[] = ["R", "RW"]
  let pathWithDefault = false

  $: searchStore.set(toLower(searchInput))
  $: compareStore.set(compareInput)
  $: stateStore.set(stateInput)
  $: defaultStore.set(pathWithDefault)
  $: yangPaths.set(diff)
</script>

<svelte:head>
	<title>NSP YANG Browser | Compare - {y} ({yKind}) with {x} ({xKind})</title>
</svelte:head>

{#if !workerComplete}
  <Loading/>
{:else}
  {#if workerStatus.status === 200}
    {#if $yangPaths.length > 0}
      <Navbar kind={xKind + ";" + yKind} basename={xBasename + ";" + yBasename} nspIp="" />
      <div class="min-w-[280px] overflow-x-auto font-nunito dark:bg-gray-800 pt-[75px]">
        <div class="px-6 pt-6 container mx-auto">
          <p class="text-gray-800 dark:text-gray-300">Compare</p>
          <SearchInput bind:searchInput />
          <div class="overflow-x-auto scroll-light dark:scroll-dark">
            <div class="py-2 space-x-2 flex items-center text-sm">
              <ChangesButton bind:compareInput />
              <StateButton bind:stateInput />
              <WithDefaultCheck bind:pathWithDefault />
            </div>
          </div>
          <Pagination />
          <div class="overflow-x-auto rounded-t-lg max-w-full mt-2">
            <table class="text-left w-full text-xs">
              <colgroup>
                <col span="1" class="w-[3%]">
                <col span="1" class="w-[2%]">
                <col span="1" class="w-[75%]">
                <col span="1" class="w-[20%]">
              </colgroup>
              <thead class="text-sm text-gray-800 dark:text-gray-300 bg-gray-300 dark:bg-gray-700">
                <tr>
                  <th scope="col" class="px-3 py-2"></th>
                  <th scope="col" class="px-3 py-2"></th>
                  <th scope="col" class="px-3 py-2">Path</th>
                  <th scope="col" class="px-3 py-2">Type</th>
                </tr>
              </thead>
              <tbody>
                {#if $total > 0}
                  {#each $paginated as item}
                    {@const path = markFilter(item.path, $searchStore)}
                    {@const type = markFilter(item.type, $searchStore)}
                    <tr class="bg-white dark:bg-gray-800 border-b dark:border-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600 hover:cursor-pointer" on:click={() => popupDetail = item}>
                      <td class="px-3 py-1.5 font-fira text-[13px] tracking-tight">{item.compare}</td>
                      <td class="px-3 py-1.5 font-fira text-[13px] tracking-tight">{item["added-filter"]}</td>
                      <td class="px-3 py-1.5 font-fira text-[13px] tracking-tight"><div use:markRender={path}></div></td>
                      <td class="px-3 py-1.5 font-fira text-[13px] tracking-tight">
                        {#if item.compare === "~"}
                          {@const fromType = markFilter(item.fromType || "", searchInput)}
                          <div class="inline-flex text-gray-400 dark:text-gray-500">from: <div class="ml-1" use:markRender={fromType}></div></div>
                          <div use:markRender={type}></div>
                        {:else}
                          <div use:markRender={type}></div>
                        {/if}
                      </td>
                    </tr>
                  {/each}
                {:else}
                  <tr class="bg-white border-b dark:bg-gray-800 dark:border-gray-700">
                    <td colspan="4" class="px-3 py-1.5 font-fira text-[13px] text-gray-400 dark:text-gray-500 text-center">{workerComplete ? 'No results found' : 'Yang compare under process...'}</td>
                  </tr>
                {/if}
              </tbody>
            </table>
          </div>
          <Pagination />
          <Popup {popupDetail} />
        </div>
        <Footer home={false}/>
      </div>
    {:else}
      <NothingToCompare {x} {y}/>
    {/if}
  {:else}
    <ErrorNotification pageError={workerStatus} />
  {/if}
{/if}