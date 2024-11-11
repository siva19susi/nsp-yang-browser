<script lang="ts">
  import { onMount } from 'svelte'
 
  import RepoList from '$lib/components/RepoList.svelte'
	import Footer from '$lib/components/Footer.svelte'
  import Theme from '$lib/components/Theme.svelte';

	import { kindView, toLower } from '$lib/components/functions'
  import type { RepoListResponse, RepoResponseMessage } from '$lib/workers/structure'
	import { localRepoStore, localSearchFilter, localSearchStore, nspRepoStore, nspSearchFilter, nspSearchStore } from './store'
	
  // DEFAULTS
  let files: FileList
  let repoDetail: RepoListResponse
  let localRepo: RepoListResponse[] = []
  let nspRepo: RepoListResponse[] = []
  let compare: string[] = []
  let currentPanel = "local"
  let nspPanelReady = false
  let isSubmitting = false
  let nspConnectedUser = ""
  let nspConnectedIp = ""

  // BASENAME WORKER
  let basenameWorker: Worker | undefined = undefined
  async function loadRepoWorker(kind: string) {
    const RepoWorker = await import('$lib/workers/repo.worker?worker')
    const repoWorker = new RepoWorker.default()
    repoWorker.postMessage(kind)
    repoWorker.onmessage = onWorkerMessage
  }
  function onWorkerMessage(event: MessageEvent<RepoResponseMessage>) {
    const response = event.data
    if(response.success) {
      if(response.kind == "local") {
        localRepo = response.repo
      } else if(response.kind == "nsp") {
        [nspConnectedUser, nspConnectedIp] = response.message.split("@")
        nspRepo = response.repo
        nspPanelReady = true
      }
    } else {
      if(response.kind == "local") {
        window.alert(response.message)
      } else if(response.kind == "nsp") {
        nspPanelReady = true
      }
    }
  }

  // ON PAGELOAD
  onMount(() => {
    loadRepoWorker("local")
    loadRepoWorker("nsp")
  })

  // DEFAULTS
  let localSearch = ""
  let nspSearch = ""
  $: localSearchStore.set(toLower(localSearch).split(/\s+/))
  $: nspSearchStore.set(toLower(nspSearch).split(/\s+/))
  $: localRepoStore.set(localRepo)
  $: nspRepoStore.set(nspRepo)

  async function handleUpload() {
    if(files) {
      const filename = files[0].name
      const basename = `${filename.replace(".zip", "")}`
      const formData = new FormData()
      formData.append("file", files[0])

      const uploadOperation = await fetch("/api/upload", {
        method: "POST", body: formData
      })

      if(!uploadOperation.ok) {
        const errorText = await uploadOperation.text()
        window.alert(`[Error] Failed to upload ${filename}: ${errorText}`)
      }

      window.alert(`[Success] ${filename} has been uploaded. Page will be redirected to take effect.`)
      window.location.href = `local/${basename}`
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

  async function nspDisconnect(event: SubmitEvent) {
    isSubmitting = true

		const nspDisconnected = await fetch("/api/nsp/disconnect", {method: "POST"})
    if(nspDisconnected.ok) {
      window.alert(`[Success] NSP disconnected. Page will reload to take effect.`)
      window.location.reload()
    } else {
      window.alert(`[Error] Failed to disconnect NSP. Page will reload to take effect.`)
    }
  }

  function fetchCompareKey(key: string, fetch: string) {
    const [kind, basename] = key.split("@")
    if(fetch === "kind") {
      return kindView(kind)
    }
    if(fetch === "basename") {
      return basename
    }
  }
</script>
  
<svelte:head>
	<title>Yang Browser</title>
</svelte:head>

<div class="flex flex-col items-center min-h-screen pt-5 has-header-img font-nunito">
  <div class="flex-grow-0 flex-shrink-0">
    <p class="px-4"><img src="/images/nwhite.svg" width="100" alt="Logo"/></p>
  </div>
  <div class="flex-grow-1 flex-shrink-0 m-auto px-4 py-10 md:w-[600px]">
    <div>
      <div class="px-4 pt-3 flex items-center justify-between rounded-t-lg bg-gray-50 dark:bg-gray-700">
        <p class="text-black dark:text-white">Yang Browser</p>
        <Theme/>
      </div>
      <div class="p-4 bg-gray-50 dark:bg-gray-700 border-b dark:border-gray-700">
        <ul class="flex items-center space-x-2 text-sm text-center text-gray-500 dark:text-gray-400">
          <li>
            <button
              on:click={() => currentPanel = "local"}
              class="inline-block px-2 py-1 rounded-lg cursor-pointer {currentPanel === "local"
                ? 'bg-blue-600 text-white'
                : 'hover:bg-gray-200 hover:text-gray-800 dark:hover:text-gray-800 border border-gray-200 dark:border-gray-600'}">Uploads</button>
          </li>
          <li class="{nspPanelReady ? '' : 'pointer-events-none'}">
            {#if nspPanelReady}
              <button on:click={() => currentPanel = "nsp"}
                class="inline-flex items-center space-x-2 px-2 py-1 rounded-lg cursor-pointer {currentPanel === "nsp"
                ? 'bg-blue-600 text-white'
                : 'hover:bg-gray-200 hover:text-gray-800 dark:hover:text-gray-800 border border-gray-200 dark:border-gray-600'}">
                  <span>NSP</span>
                  {#if nspConnectedIp !== "" && nspConnectedUser !== ""}
                    <svg class="w-4 h-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                      <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 14v3m-3-6V7a3 3 0 1 1 6 0v4m-8 0h10a1 1 0 0 1 1 1v7a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1v-7a1 1 0 0 1 1-1Z"/>
                    </svg>
                  {:else}
                    <svg class="w-4 h-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                      <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14v3m4-6V7a3 3 0 1 1 6 0v4M5 11h10a1 1 0 0 1 1 1v7a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1v-7a1 1 0 0 1 1-1Z"/>
                    </svg>
                  {/if}
                </button>
            {:else}
              <button class="inline-flex items-center space-x-2 px-2 py-1 rounded-lg cursor-pointer 
                  hover:bg-gray-200 hover:text-gray-800 border border-gray-200 dark:border-gray-600 animate-pulse bg-gray-200 dark:bg-gray-600">
                <span>NSP</span>
                <svg class="animate-spin -ml-1 mr-3 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </button>
            {/if}
          </li>
          <li class="relative inline-flex select-none">
            <button on:click={() => currentPanel = "compare"} 
              class="inline-block px-2 py-1 rounded-lg cursor-pointer {currentPanel === "compare"
                ? 'bg-blue-600 text-white'
                : 'hover:bg-gray-200 hover:text-gray-800 dark:hover:text-gray-800 border border-gray-200 dark:border-gray-600'}">
              <div class="flex items-center space-x-1">
                <span>Compare</span>
                <svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024" fill="currentColor" stroke="currentColor" stroke-width="10" aria-hidden="true">
                  <path d="M420.266667 832c-17.066667 0-34.133333-6.4-44.8-19.2L104.533333 541.866667c-12.8-12.8-19.2-27.733333-19.2-44.8s6.4-34.133333 19.2-44.8L345.6 211.2c23.466667-23.466667 66.133333-23.466667 89.6 0l270.933333 270.933333c12.8 12.8 19.2 27.733333 19.2 44.8s-6.4 34.133333-19.2 44.8L465.066667 812.8c-10.666667 12.8-27.733333 19.2-44.8 19.2z m-29.866667-597.333333c-6.4 0-10.666667 2.133333-14.933333 6.4L134.4 482.133333c-4.266667 4.266667-6.4 8.533333-6.4 14.933334s2.133333 10.666667 6.4 14.933333L405.333333 782.933333c8.533333 8.533333 21.333333 8.533333 29.866667 0l241.066667-241.066666c4.266667-4.266667 6.4-8.533333 6.4-14.933334s-2.133333-10.666667-6.4-14.933333L405.333333 241.066667c-4.266667-4.266667-8.533333-6.4-14.933333-6.4z" />
                  <path d="M618.666667 832c-17.066667 0-34.133333-6.4-46.933334-19.2L317.866667 558.933333c-12.8-12.8-19.2-29.866667-19.2-46.933333s6.4-34.133333 19.2-46.933333L571.733333 211.2c25.6-25.6 68.266667-25.6 93.866667 0l253.866667 253.866667c25.6 25.6 25.6 68.266667 0 93.866666L665.6 812.8c-12.8 12.8-29.866667 19.2-46.933333 19.2z m0-597.333333c-6.4 0-12.8 2.133333-17.066667 6.4L347.733333 494.933333c-4.266667 4.266667-6.4 10.666667-6.4 17.066667s2.133333 12.8 6.4 17.066667l253.866667 253.866666c8.533333 8.533333 23.466667 8.533333 34.133333 0l253.866667-253.866666c8.533333-8.533333 8.533333-23.466667 0-34.133334L635.733333 241.066667c-4.266667-4.266667-10.666667-6.4-17.066666-6.4zM332.8 480z" />
                </svg>
              </div>
            </button>
          </li>
        </ul>
      </div>
    </div>
    <div class="bg-white dark:bg-gray-800 rounded-b-lg shadow-xl">
      <div class="h-[480px]">
        <!--UPLOADS-->
        <div class="flex flex-col h-full {currentPanel === "local" ? 'block' : 'hidden'}">
          <div class="{$localRepoStore.length === 0 ? 'flex-grow' : ''} p-4">
            <input id="dropzone" type="file" class="peer hidden" accept="application/zip" bind:files on:change={handleUpload} />
            <label for="dropzone" class="flex flex-col space-y-2 px-4 py-3 h-full w-full items-center justify-center cursor-pointer rounded-lg text-gray-500 dark:text-gray-400 bg-gray-100 hover:bg-gray-200 dark:bg-gray-800 dark:hover:bg-gray-700 border-2 border-dashed border-gray-200 dark:border-gray-700">
              <svg class="w-5 h-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v9m-5 0H5a1 1 0 0 0-1 1v4a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1v-4a1 1 0 0 0-1-1h-2M8 9l4-5 4 5m1 8h.01"/>
              </svg>                
              <p class="text-sm">Click to upload a YANG repo</p>
              {#if files}
                <p class="text-black dark:text-white text-xs"><span>Selected:</span> {`${files[0].name} (${files[0].size} bytes)`}</p>
                <div class="rounded-md h-4 w-4 border-2 border-blue-300 animate-spin"></div>
              {:else}
                <p class="text-xs">Supported file format .zip (max 10 MB)</p>
              {/if}
            </label>
          </div>
          {#if $localRepoStore.length > 0}
            <div class="px-4 py-2 border-y dark:border-gray-700">
              <input type="text" bind:value={localSearch} placeholder="Search..." 
              class="px-3 py-2 rounded-lg w-full text-[12.5px] text-gray-800 dark:text-gray-200 
                dark:placeholder-gray-400 border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700">
            </div>
          {/if}
          {#if $localSearchFilter.length > 0}
            <ul class="mb-1 overflow-y-auto scroll-light dark:scroll-dark">
              {#each $localSearchFilter as {name, files}, i}
                {@const compareValue = "local@" + name}
                {@const isDisabled = compare.length === 2 && !compare.includes(compareValue)}
                <li class="text-sm bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 {i > 0 ? 'border-t dark:border-gray-700' : ''}">
                  <div class="flex items-center justify-between">
                    <a data-sveltekit-reload href="/local/{name}" class="px-4 py-3 w-full overflow-x-auto">{name}</a>
                    <div class="flex items-center mx-4 space-x-5">
                      <button class="text-gray-600 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-800 rounded-lg p-1" on:click={() => repoDetail = {name, files}}>
                        <svg class="w-5 h-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                          <path stroke="currentColor" stroke-linecap="round" stroke-width="2" d="M9 8h10M9 12h10M9 16h10M4.99 8H5m-.02 4h.01m0 4H5"/>
                        </svg>
                      </button>
                      <div title="Add to compare" class="flex">
                        <input type="checkbox" id="local-{name}-check" bind:group={compare} value="{compareValue}" disabled={isDisabled} class="peer hidden" />
                        <label for="local-{name}-check" class="select-none p-1 rounded-lg peer-checked:bg-blue-600 peer-checked:hover:bg-blue-700 peer-checked:text-white {isDisabled ? 'cursor-not-allowed text-gray-200 dark:text-gray-600' : 'cursor-pointer hover:bg-blue-600 hover:text-white'}">
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
          {/if}
        </div>
        <!--NSP-->
        <div class="flex flex-col h-full {currentPanel === "nsp" ? 'block' : 'hidden'}">
          {#if nspConnectedIp === ""}
            <div class="flex-grow p-4 bg-gray-200 dark:bg-gray-800 rounded-b-lg">
              <div class="flex items-center justify-center h-full">
                <form class="p-6 w-full" method="POST" action="?/nspConnect" on:submit|preventDefault={nspConnect}>
                  <div class="w-full mb-6">
                    <label for="nsp-ip" class="block uppercase text-gray-800 dark:text-gray-200 text-xs font-bold mb-2">IP / FQDN</label>
                    <input id="nsp-ip" name="ip" type="text" class="px-3 py-2 rounded-lg w-full text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
                  </div>
                  <div class="flex mb-10">
                    <div class="w-full md:w-1/2 pr-2 mb-6 md:mb-0">
                      <label for="nsp-user" class="block uppercase text-gray-800 dark:text-gray-200 text-xs font-bold mb-2">Username</label>
                      <input id="nsp-user" name="user" type="text" class="px-3 py-2 rounded-lg w-full text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
                    </div>
                    <div class="w-full md:w-1/2 pl-2">
                      <label for="nsp-pass" class="block uppercase text-gray-800 dark:text-gray-200 text-xs font-bold mb-2">Password</label>
                      <input id="nsp-pass" name="pass" type="password" class="px-3 py-2 rounded-lg w-full text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 {isSubmitting ? 'bg-gray-300' : 'bg-gray-100'}" disabled={isSubmitting}>
                    </div>
                  </div>
                  <div class="flex items-center justify-end">
                    <button type="submit" class="px-4 py-1 rounded-lg text-sm text-white bg-green-600 hover:bg-green-700 {isSubmitting ? 'animate-pulse' : ''}" disabled={isSubmitting}>{isSubmitting ? 'Connecting...' : 'Connect'}</button>
                  </div>
                </form>
              </div>
            </div>
          {:else}
            <div class="p-4">
              <form class="p-4 rounded-lg bg-gray-100 dark:bg-gray-800 border-2 border-dashed border-gray-200 dark:border-gray-700" method="POST" action="?/nspDisconnect" on:submit|preventDefault={nspDisconnect}>
                <div class="flex mb-4">
                  <div class="w-full md:w-1/2 pr-2 mb-6 md:mb-0">
                    <label for="nsp-connect-ip" class="block uppercase text-gray-700 dark:text-gray-300 text-xs font-bold mb-2">IP / FQDN</label>
                    <input id="nsp-connect-ip" bind:value={nspConnectedIp} type="text" class="px-3 py-2 rounded-lg w-full text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-800 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 select-none pointer-events-none" readonly>
                  </div>
                  <div class="w-full md:w-1/2 pl-2">
                    <label for="nsp-connect-user" class="block uppercase text-gray-700 dark:text-gray-300 text-xs font-bold mb-2">Username</label>
                    <input id="nsp-connect-user" bind:value={nspConnectedUser} type="text" class="px-3 py-2 rounded-lg w-full text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-800 dark:text-gray-200 bg-gray-50 dark:bg-gray-700 select-none pointer-events-none" readonly>
                  </div>
                </div>
                <div class="flex items-center justify-end">
                  <button type="submit" class="px-4 py-1 rounded-lg text-xs text-white bg-red-600 hover:bg-red-700 {isSubmitting ? 'animate-pulse' : ''}" disabled={isSubmitting}>{isSubmitting ? 'Disconnecting...' : 'Disconnect'}</button>
                </div>
              </form>
            </div>
          {/if}
          {#if $nspRepoStore.length > 0}
            <div class="px-4 py-2 border-y dark:border-gray-700">
              <input type="text" bind:value={nspSearch} placeholder="Search..." 
              class="px-3 py-2 rounded-lg w-full text-[12.5px] text-gray-800 dark:text-gray-200 
                dark:placeholder-gray-400 border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700">
            </div>
          {/if}
          {#if $nspSearchFilter.length > 0}
            <ul class="mb-1 overflow-y-auto scroll-light">
              {#each $nspSearchFilter as {name}, i}
                {@const compareValue = "nsp@" + name}
                {@const isDisabled = compare.length === 2 && !compare.includes(compareValue)}
                <li class="text-sm bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 {i > 0 ? 'border-t dark:border-gray-700' : ''} hover:bg-gray-100">
                  <div class="flex items-center justify-between">
                    <a data-sveltekit-reload href="/nsp/{name}" class="px-4 py-3 w-full overflow-x-auto">{name}</a>
                    <div class="flex items-center mx-4 space-x-5">
                      <div title="Add to compare" class="flex">
                        <input type="checkbox" id="nsp-{name}-check" bind:group={compare} value="{compareValue}" disabled={isDisabled} class="peer hidden" />
                        <label for="nsp-{name}-check" class="select-none p-1 rounded-lg peer-checked:bg-blue-600 peer-checked:hover:bg-blue-700 peer-checked:text-white {isDisabled ? 'cursor-not-allowed text-gray-200 dark:text-gray-600' : 'cursor-pointer hover:bg-blue-600 hover:text-white'}">
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
          {/if}
        </div>
        <!--COMPARE-->
        <div class="flex flex-col h-full {currentPanel === "compare" ? 'block' : 'hidden'}">
          <div class="flex-grow bg-gray-100 dark:bg-gray-800">
            <div class="flex items-center p-4 text-sm space-x-3 select-none h-full">
              <div class="w-full text-center">
                <p class="px-3 py-20 rounded-t-lg border-x-2 border-t-2 border-dashed border-gray-400 dark:border-gray-500 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200">
                  {compare.length > 0 ? fetchCompareKey(compare[0], "basename") : 'X'}
                </p>
                <p class="p-3 {compare.length > 0 ? 'bg-gray-600' : ''} bg-gray-400 dark:bg-gray-500 text-white rounded-b-lg">
                  {compare.length > 0 ? fetchCompareKey(compare[0], "kind") : '{{ source }}'}
                </p>
              </div>
              <p class="rounded-full px-3 py-2 bg-blue-600 text-white">to</p>
              <div class="w-full text-center">
                <p class="px-3 py-20 rounded-t-lg border-x-2 border-t-2 border-dashed border-gray-400 dark:border-gray-500 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200">
                  {compare.length > 1 ? fetchCompareKey(compare[1], "basename") : 'Y'}
                </p>
                <p class="p-3 {compare.length > 1 ? 'bg-gray-600' : ''} bg-gray-400 dark:bg-gray-500 text-white rounded-b-lg">
                  {compare.length > 1 ? fetchCompareKey(compare[1], "kind") : '{{ source }}'}
                </p>
              </div>
            </div>
          </div>
          <div class="text-center px-3 py-4 bg-gray-50 dark:bg-gray-700 border-t dark:border-gray-700 rounded-b-lg">
            <a href="/compare/{compare[0]}..{compare[1]}" 
              class="px-4 py-1 rounded-lg text-sm text-white select-none 
                {compare.length === 2 ? 
                  'bg-green-600 hover:bg-green-700' : 
                  'bg-gray-300 dark:bg-gray-600 text-gray-500 dark:text-gray-500 pointer-events-none'}">
              Click to Compare
            </a>
            <div class="px-4 pt-6 text-[10px] text-gray-800 dark:text-gray-200 italic text-left">
              <p>1) Order or compare selection determines the value of X (first selected) and Y (second selected).</p>
              <p>2) The compare provides added paths, removed paths and paths with modified type definition. 
                For more detailed info regarding each repo, contact the repo source Admin.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <RepoList bind:repoDetail />
  <div class="flex-grow-0 flex-shrink-0">
    <Footer home={true} />
  </div>
</div>