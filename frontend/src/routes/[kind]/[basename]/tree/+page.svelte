<script lang="ts">
  import { onMount } from "svelte"
  import { page } from "$app/stores"
  import { goto } from "$app/navigation"

  import Header from '$lib/components/Header.svelte'
  import Footer from '$lib/components/Footer.svelte'
  import SearchInput from "$lib/components/SearchInput.svelte"
  import StateButton from '$lib/components/StateButton.svelte'
  import ShowPrefixCheck from '$lib/components/ShowPrefixCheck.svelte'
  import WithDefaultCheck from '$lib/components/WithDefaultCheck.svelte'
  import CrossBrowser from '$lib/components/crossBrowser.svelte'
  import ErrorNotification from "$lib/components/ErrorNotification.svelte"
  
  import Popup from '$lib/components/Popup.svelte'
  import Loading from '$lib/components/Loading.svelte'
  import YangTree from './YangTree.svelte'

  import { decideExpand } from "./expand"
  import { toLower, kindView } from "$lib/components/functions"
  import { pathFocus } from '$lib/components/sharedStore'
  import { defaultStore, prefixStore, searchStore, stateStore, yangTarget, yangTreeArgs } from "./store"

  import type { TreePayLoad } from '$lib/structure'
  import type { YangTreeResponseMessage, YangTreePaths } from "$lib/workers/structure"

  const getUrlPath = () => $page.data.urlPath
  const isCrossLaunched = () => $page.data.crossLaunched
  
  // DEFAULTS
  let popupDetail = {}
  let treePaths: YangTreePaths
  let workerComplete = false
  let workerStatus = {status: 404, error: {message: "Unknown Error"}}

  // YANGTREE WORKER
  let yangTreeWorker: Worker | undefined = undefined
  async function loadYangTreeWorker (kind: string, basename: string, searchInput: string, prefixInput: boolean, stateInput: string, defaultInput: boolean) {
    const YangTreeWorker = await import('$lib/workers/yangTree.worker?worker')
    yangTreeWorker = new YangTreeWorker.default()
    yangTreeWorker.postMessage({ kind, basename, searchInput, prefixInput, stateInput, defaultInput })
    yangTreeWorker.onmessage = onYangTreeWorkerMessage
  }
  function onYangTreeWorkerMessage(event: MessageEvent<YangTreeResponseMessage>) {
    const response = event.data
    workerStatus.error.message = response.message
    if(event.data.success) {
      treePaths = response.node
      workerStatus.status = 200
    }
    workerComplete = true
  }

  // ON PAGELOAD
	export let data
  let {kind, basename} = data  

  // OTHER BINDING VARIABLES
  let searchInput: string = isCrossLaunched() ? "" : getUrlPath()
  let stateInput = ""
  let showPathPrefix = false
  let pathWithDefault = false

  let pastYangTreeArgs = `${searchInput};;${showPathPrefix};;;;${pathWithDefault}`

  onMount(() => loadYangTreeWorker(kind, basename, searchInput, showPathPrefix, "", pathWithDefault))

  pathFocus.set({})
	pathFocus.subscribe((value) => {
    popupDetail = value
  })

  $: {
    searchStore.set(toLower(searchInput))
    stateStore.set(stateInput)
    prefixStore.set(showPathPrefix)
    defaultStore.set(pathWithDefault)
  }
  $: yangTarget.set(treePaths)

  // TRIGGER SEARCH FILTERS
  function triggerApply() {
    if(pastYangTreeArgs !== $yangTreeArgs) {
      pastYangTreeArgs = $yangTreeArgs
      $page.url.searchParams.delete("from")
      if($searchStore != "") {
        $page.url.searchParams.set("path", $searchStore)
      } else {
        $page.url.searchParams.delete("path")
      }
      goto(`?${$page.url.searchParams.toString()}`, {invalidateAll: true})
      loadYangTreeWorker(kind, basename, $searchStore, $prefixStore, $stateStore, $defaultStore)
    }
	}
</script>

<svelte:head>
	<title>Yang Tree Browser {basename} ({kindView(kind)})</title>
</svelte:head>

<svelte:window on:keyup={({key}) => key === "Enter" ? triggerApply() : ""} />

{#if !workerComplete}
  <Loading/>
{:else}
  {#if workerStatus.status === 200}
    <Header {kind} {basename} />
    <div class="min-w-[280px] overflow-x-auto font-nunito dark:bg-gray-800 pt-[75px] lg:pt-[85px]">
      <div class="px-6 py-7 container mx-auto">
        <div class="flex items-center justify-between">
          <p class="text-gray-800 dark:text-gray-300">Tree Browser</p>
          <CrossBrowser {kind} {basename} isTree={true} />
        </div>
        <SearchInput bind:searchInput />
        <div class="flex py-1 items-center space-x-2">
          <StateButton bind:stateInput />
          <ShowPrefixCheck bind:showPathPrefix />
          <WithDefaultCheck bind:pathWithDefault />
          <div class="dropdown">
            <a href="/{kind}/{basename}/tree/payload" target="_blank" class="dropdown-button px-2 py-1 rounded-lg text-xs text-nowrap bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-white inline-flex items-center align-bottom">
              <svg class="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 14v4.833A1.166 1.166 0 0 1 16.833 20H5.167A1.167 1.167 0 0 1 4 18.833V7.167A1.166 1.166 0 0 1 5.167 6h4.618m4.447-2H20v5.768m-7.889 2.121 7.778-7.778"/>
              </svg>
              <span class="ml-1">Generate Payload</span>
            </a>
            <div class="dropdown-content absolute z-10 hidden bg-gray-100 dark:bg-gray-700 dark:text-white rounded-lg shadow">
              <p class="my-2 max-w-[300px] px-3 text-xs">
                The payload generated is purely a sample based on default values and type definition of the yang model. Contact repo source Admin for more details.
              </p>
            </div>
          </div>
        </div>
        <div class="text-right mt-6">
          <button class="px-4 py-2 rounded-lg text-xs 
            {pastYangTreeArgs === $yangTreeArgs ? 'bg-green-100 dark:bg-green-900 text-gray-500 dark:text-gray-500 cursor-not-allowed' : 'text-white bg-green-600 hover:bg-green-700 dark:bg-green-700 dark:hover:bg-green-800'}" 
            disabled={pastYangTreeArgs === $yangTreeArgs} on:click={triggerApply}>Apply
          </button>
        </div>
      </div>
      {#if Object.keys($yangTarget)?.length}
        <div class="px-5 py-4 container mx-auto border-t dark:border-gray-600">
          <div class="font-fira text-xs tracking-tight">
            {#each $yangTarget.children as folder}
              <YangTree {folder} withPrefix={showPathPrefix} expanded={decideExpand(folder, isCrossLaunched(), getUrlPath())} />
            {/each}
          </div>
          <Popup {kind} {basename} {popupDetail} />
          <Footer home={false}/>
        </div>
      {/if}
    </div>
  {:else}
    <ErrorNotification pageError={workerStatus} />
  {/if}
{/if}
