<script lang="ts">
  import { onMount } from 'svelte'

  import Navbar from '$lib/components/Navbar.svelte'
  import Footer from '$lib/components/Footer.svelte'
  import LsoPagination from './LsoPagination.svelte'
  import TelemetryPagination from './TelemetryPagination.svelte'
	
  import { intentTypeStore, total, pageCount, start, end, lsoStore, telemetryStore, lsoSearch, telemetrySearch, telemetryPaginated, lsoPaginated } from './store'
	import type { IntentTypeSearch, IntentTypeSearchResponseMessage } from './structure'
	import { compare } from '$lib/components/sharedStore'

  let search = ""
  let isSubmitting = false
  let typingTimer: number | undefined
  let doneTypingDelay = 500

  let intentTypes: IntentTypeSearch = {
	  total: 0,
	  pageCount: 0,
	  intentTypes: []
  }

  let workerStatus = {
    complete: false, 
    progress: 30, 
    progressText: "(if available)", 
    error: {
      message: "Unknown Error"
    }
  }

  let intentTypeWorker: Worker | undefined = undefined
  async function loadWorker(page: number = 1, filter: string = "") {
    const IntentTypeWorker = await import('./intentType.worker?worker')
    intentTypeWorker = new IntentTypeWorker.default()
    intentTypeWorker.postMessage({page, filter})
    intentTypeWorker.onmessage = onWorkerMessage
  }
  function onWorkerMessage(event: MessageEvent<IntentTypeSearchResponseMessage>) {
    const response = event.data
    if (response.type === "progress") {
      workerStatus.progress = response.value
    }
    if (response.type === "complete") {
      workerStatus.progress = 100
      workerStatus.error.message = response.message
      workerStatus.complete = true
      if(response.success) {
        intentTypes = response.intentTypes
      }
    }
  }

  async function nspConnect(event: SubmitEvent) {
    isSubmitting = true
    const data = new FormData(event.currentTarget as HTMLFormElement)
		const ip = data.get("ip")
		const user = data.get("user")
		const pass = data.get("pass")
		const nspConnected = await fetch("/api/nsp/connect", {
      method: "POST", body: JSON.stringify({ip, user, pass})
    })
    if(nspConnected.ok) {
      window.alert(`[Success] NSP connected. Page will reload to take effect.`)
      window.location.reload()
    } else {
      window.alert(`[Error] Failed to connect to NSP`)
    }
    isSubmitting = false
  }

  function onSearchKeyUp(f: string) {
    clearTimeout(typingTimer)
    typingTimer = setTimeout(() => {
      updateTable(0, f)
    }, doneTypingDelay)
  }

  function onSearchKeyDown() {
    clearTimeout(typingTimer)
  }

  function updateTable(s: number, f: string = "") {
    const pageNumber = s/30+1
    workerStatus.complete = false
    workerStatus.progressText = (f != "" ? "(filtering by search)" : `(page ${pageNumber})`)
    if(s >= 0 && s < $total) {
      start.set(s)
    }
    loadWorker(pageNumber, f)
  }

  export let data
  const { nspInfo, modules, lsoOperations, telemetryTypes } = data
  const {ip: nspIp, user: nspUser} = nspInfo
  const nspConnected = (nspIp !== "" ? true : false)
  onMount(() => {
    start.set(0)
    if(nspConnected) {
      loadWorker()
    }
  })

  $: intentTypeStore.set(intentTypes.intentTypes)
  $: total.set(intentTypes.total)
  $: pageCount.set(intentTypes.pageCount)
  $: lsoStore.set(lsoOperations)
  $: telemetryStore.set(telemetryTypes)
  $: selected = $compare
</script>

<svelte:head>
	<title>NSP YANG Browser | Telemetry Type Definition</title>
</svelte:head>

<Navbar {nspIp}/>
<div class="font-nunito text-sm container mx-auto pt-[75px]">
  {#if !nspConnected}
    <div class="px-6 py-4">
      <form method="POST" action="?/nspConnect" on:submit|preventDefault={nspConnect}>
        <div class="grid md:grid-cols-3 gap-4">
          <div>
            <label for="nsp-ip" class="block uppercase text-gray-800 dark:text-gray-200 text-xs mb-2">IP / FQDN</label>
            <input id="nsp-ip" name="ip" type="text" class="px-3 py-2 rounded-lg w-full text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
          </div>
          <div>
            <label for="nsp-user" class="block uppercase text-gray-800 dark:text-gray-200 text-xs mb-2">Username</label>
            <input id="nsp-user" name="user" type="text" class="px-3 py-2 rounded-lg w-full text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
          </div>
          <div>
            <label for="nsp-pass" class="block uppercase text-gray-800 dark:text-gray-200 text-xs mb-2">Password</label>
            <input id="nsp-pass" name="pass" type="password" class="px-3 py-2 rounded-lg w-full text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
          </div>
        </div>
        <div class="flex items-center justify-end pt-4">
          <button type="submit" class="px-4 py-2 rounded-lg text-sm text-white bg-green-600 hover:bg-green-700 {isSubmitting ? 'animate-pulse' : ''}" disabled={isSubmitting}>{isSubmitting ? 'Connecting...' : 'Connect'}</button>
        </div>
      </form>
    </div>
  {:else}
    <div class="px-6 pt-2 pb-4">
      <p class="text-lg pb-2 text-black dark:text-white">Modules</p>
      <div class="py-2 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:lg:grid-cols-5 2xl:lg:grid-cols-6 gap-4 pb-8">
        {#each modules.sort() as module}
          <a data-sveltekit-reload href="/nsp-module/{module}" class="font-medium rounded-lg text-sm px-3 py-2 text-gray-900 dark:text-white bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 border border-gray-300 dark:border-gray-600 dark:hover:border-gray-600">{module}</a>
        {/each}
      </div>
    </div>
    {#if $lsoStore?.length || $telemetryStore?.length }
      <div class="md:flex md:items-center md:justify-between border-t dark:border-gray-700">
        {#if $lsoStore?.length}
          <div class="px-6 pt-6 pb-4 w-full">
            <div class="flex items-center justify-between">
              <p class="text-lg pb-2 text-nokia dark:text-white">LSO Operations</p>
            </div>
            <div class="py-2">
              <input type="text" placeholder="Search..." bind:value={$lsoSearch} 
                class="px-3 py-2 rounded-lg w-full text-[12.5px] text-gray-800 dark:text-gray-200 
                  dark:placeholder-gray-400 border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700">
              <LsoPagination/>
              <ul class="mb-2 border-t dark:border-gray-700 overflow-y-auto scroll-light dark:scroll-dark max-h-[495px]">
                {#each $lsoPaginated.sort() as name, i}
                  {@const compareValue = "nsp-lso-operation@" + name}
                  {@const isDisabled = selected.length === 2 && !selected.includes(compareValue)}
                  <li class="text-sm bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 {i > 0 ? 'border-t dark:border-gray-700' : ''} hover:bg-gray-100">
                    <div class="flex items-center justify-between">
                      <a data-sveltekit-reload href="/nsp-lso-operation/{name}" class="px-4 py-3 w-full overflow-x-auto">{name}</a>
                      <div class="flex items-center mx-4 space-x-5">
                        <div title="Add to compare" class="flex">
                          <input type="checkbox" id="nsp-lso-operation-{name}-check" class="peer hidden" 
                            disabled={isDisabled}
                            checked={selected.includes(compareValue)} 
                            on:change={(e) => e.currentTarget.checked ? compare.add(compareValue) : compare.remove(compareValue)} 
                          />
                          <label for="nsp-lso-operation-{name}-check" class="select-none p-1 rounded-lg peer-checked:bg-blue-600 peer-checked:hover:bg-blue-700 peer-checked:text-white {isDisabled ? 'cursor-not-allowed text-gray-200 dark:text-gray-600' : 'cursor-pointer hover:bg-blue-600 hover:text-white'}">
                            <svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024" fill="currentColor" stroke="currentColor" stroke-width="10" aria-hidden="true">
                              <path d="M420.266667 832c-17.066667 0-34.133333-6.4-44.8-19.2L104.533333 541.866667c-12.8-12.8-19.2-27.733333-19.2-44.8s6.4-34.133333 19.2-44.8L345.6 211.2c23.466667-23.466667 66.133333-23.466667 89.6 0l270.933333 270.933333c12.8 12.8 19.2 27.733333 19.2 44.8s-6.4 34.133333-19.2 44.8L465.066667 812.8c-10.666667 12.8-27.733333 19.2-44.8 19.2z m-29.866667-597.333333c-6.4 0-10.666667 2.133333-14.933333 6.4L134.4 482.133333c-4.266667 4.266667-6.4 8.533333-6.4 14.933334s2.133333 10.666667 6.4 14.933333L405.333333 782.933333c8.533333 8.533333 21.333333 8.533333 29.866667 0l241.066667-241.066666c4.266667-4.266667 6.4-8.533333 6.4-14.933334s-2.133333-10.666667-6.4-14.933333L405.333333 241.066667c-4.266667-4.266667-8.533333-6.4-14.933333-6.4z" />
                              <path d="M618.666667 832c-17.066667 0-34.133333-6.4-46.933334-19.2L317.866667 558.933333c-12.8-12.8-19.2-29.866667-19.2-46.933333s6.4-34.133333 19.2-46.933333L571.733333 211.2c25.6-25.6 68.266667-25.6 93.866667 0l253.866667 253.866667c25.6 25.6 25.6 68.266667 0 93.866666L665.6 812.8c-12.8 12.8-29.866667 19.2-46.933333 19.2z m0-597.333333c-6.4 0-12.8 2.133333-17.066667 6.4L347.733333 494.933333c-4.266667 4.266667-6.4 10.666667-6.4 17.066667s2.133333 12.8 6.4 17.066667l253.866667 253.866666c8.533333 8.533333 23.466667 8.533333 34.133333 0l253.866667-253.866666c8.533333-8.533333 8.533333-23.466667 0-34.133334L635.733333 241.066667c-4.266667-4.266667-10.666667-6.4-17.066666-6.4zM332.8 480z" />
                            </svg>
                          </label>
                        </div>
                      </div>
                    </div>
                  </li>
                {/each}
              </ul>
            </div>
          </div>
        {/if}
        {#if $telemetryStore?.length}
          <div class="px-6 pt-6 pb-4 w-full">
            <div class="flex items-center justify-between">
              <p class="text-lg pb-2 text-black dark:text-white">Telemetry Types</p>
            </div>
            <div class="py-2">
              <input type="text" placeholder="Search..." bind:value={$telemetrySearch} 
                class="px-3 py-2 rounded-lg w-full text-[12.5px] text-gray-800 dark:text-gray-200 
                  dark:placeholder-gray-400 border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700">
              <TelemetryPagination/>
              <ul class="mb-2 border-t dark:border-gray-700 overflow-y-auto scroll-light dark:scroll-dark max-h-[495px]">
                {#each $telemetryPaginated as name, i}
                  <li class="text-sm bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 {i > 0 ? 'border-t dark:border-gray-700' : ''} hover:bg-gray-100">
                    <div class="flex items-center justify-between">
                      <a data-sveltekit-reload href="/nsp/telemetry?type={name}" class="px-4 py-3 w-full overflow-x-auto">{name}</a>
                    </div>
                  </li>
                {/each}
              </ul>
            </div>
          </div>
        {/if}
      </div>
    {/if}
    <div class="px-6 pt-6 pb-4 border-t dark:border-gray-700">
      <div class="flex items-center justify-between">
        <p class="text-lg pb-2 text-black dark:text-white">Intent Types</p>
      </div>
      <div class="py-2">
        <input type="text" placeholder="Search..." bind:value={search} on:keyup={() => onSearchKeyUp(search)} on:keydown={onSearchKeyDown}
          class="px-3 py-2 rounded-lg w-full text-[12.5px] text-gray-800 dark:text-gray-200 
            dark:placeholder-gray-400 border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700">

        {#if !workerStatus.complete}
          <div class="px-1 py-2">
            <div class="w-full bg-gray-200 rounded-full h-1 dark:bg-gray-700">
              <div class="bg-blue-600 h-1 rounded-full" style="width: {workerStatus.progress}%"></div>
            </div>
            <p class="pt-2 text-black dark:text-white">Loading Intent Types {workerStatus.progressText}...</p>
          </div>
        {:else}
          {#if $intentTypeStore.length > 0}
            <div class="flex items-center justify-end py-3 text-sm mt-2">
              <p class="mr-2 text-gray-800 dark:text-gray-200">{$start + 1} - {$end > 1 ? $end : 0} of {$total}</p>
              <button class="ml-2 text-white rounded 
                  {$start == 0 ? 'bg-gray-300 dark:bg-gray-500 opacity-50 cursor-not-allowed' : 'bg-gray-400 hover:bg-gray-600'}" 
                  disabled="{$start == 0}" on:click={() => updateTable($start - 30)}>
                <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clip-rule="evenodd"/>
                </svg>
              </button>
              <button class="ml-2 text-white rounded 
                  {$end == $total ? 'bg-gray-300 dark:bg-gray-500 opacity-50 cursor-not-allowed' : 'bg-gray-400 hover:bg-gray-600'}" 
                  disabled="{$end == $total}" on:click={() => updateTable($end)}>
                <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clip-rule="evenodd"/>
                </svg>
              </button>
            </div>
            <ul class="mb-2 border-t dark:border-gray-700 overflow-y-auto scroll-light dark:scroll-dark max-h-[495px]">
              {#each $intentTypeStore.sort() as name, i}
                {@const compareValue = "nsp-intent-type@" + name}
                {@const isDisabled = selected.length === 2 && !selected.includes(compareValue)}
                <li class="text-sm bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 {i > 0 ? 'border-t dark:border-gray-700' : ''} hover:bg-gray-100">
                  <div class="flex items-center justify-between">
                    <a data-sveltekit-reload href="/nsp-intent-type/{name}" class="px-4 py-3 w-full overflow-x-auto">{name}</a>
                    <div class="flex items-center mx-4 space-x-5">
                      <div title="Add to compare" class="flex">
                        <input type="checkbox" id="nsp-intent-type-{name}-check" class="peer hidden" 
                          disabled={isDisabled}
                          checked={selected.includes(compareValue)} 
                          on:change={(e) => e.currentTarget.checked ? compare.add(compareValue) : compare.remove(compareValue)} 
                        />
                        <label for="nsp-intent-type-{name}-check" class="select-none p-1 rounded-lg peer-checked:bg-blue-600 peer-checked:hover:bg-blue-700 peer-checked:text-white {isDisabled ? 'cursor-not-allowed text-gray-200 dark:text-gray-600' : 'cursor-pointer hover:bg-blue-600 hover:text-white'}">
                          <svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024" fill="currentColor" stroke="currentColor" stroke-width="10" aria-hidden="true">
                            <path d="M420.266667 832c-17.066667 0-34.133333-6.4-44.8-19.2L104.533333 541.866667c-12.8-12.8-19.2-27.733333-19.2-44.8s6.4-34.133333 19.2-44.8L345.6 211.2c23.466667-23.466667 66.133333-23.466667 89.6 0l270.933333 270.933333c12.8 12.8 19.2 27.733333 19.2 44.8s-6.4 34.133333-19.2 44.8L465.066667 812.8c-10.666667 12.8-27.733333 19.2-44.8 19.2z m-29.866667-597.333333c-6.4 0-10.666667 2.133333-14.933333 6.4L134.4 482.133333c-4.266667 4.266667-6.4 8.533333-6.4 14.933334s2.133333 10.666667 6.4 14.933333L405.333333 782.933333c8.533333 8.533333 21.333333 8.533333 29.866667 0l241.066667-241.066666c4.266667-4.266667 6.4-8.533333 6.4-14.933334s-2.133333-10.666667-6.4-14.933333L405.333333 241.066667c-4.266667-4.266667-8.533333-6.4-14.933333-6.4z" />
                            <path d="M618.666667 832c-17.066667 0-34.133333-6.4-46.933334-19.2L317.866667 558.933333c-12.8-12.8-19.2-29.866667-19.2-46.933333s6.4-34.133333 19.2-46.933333L571.733333 211.2c25.6-25.6 68.266667-25.6 93.866667 0l253.866667 253.866667c25.6 25.6 25.6 68.266667 0 93.866666L665.6 812.8c-12.8 12.8-29.866667 19.2-46.933333 19.2z m0-597.333333c-6.4 0-12.8 2.133333-17.066667 6.4L347.733333 494.933333c-4.266667 4.266667-6.4 10.666667-6.4 17.066667s2.133333 12.8 6.4 17.066667l253.866667 253.866666c8.533333 8.533333 23.466667 8.533333 34.133333 0l253.866667-253.866666c8.533333-8.533333 8.533333-23.466667 0-34.133334L635.733333 241.066667c-4.266667-4.266667-10.666667-6.4-17.066666-6.4zM332.8 480z" />
                          </svg>
                        </label>
                      </div>
                    </div>
                  </div>
                </li>
              {/each}
            </ul>
          {:else}
            <p class="py-2 text-red-600">No records found</p>
          {/if}
        {/if}
      </div>
    </div>
  {/if}
  <Footer home={false}/>
</div>