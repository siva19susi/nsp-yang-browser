<script lang="ts">
  import { onMount } from 'svelte'

  import Navbar from '$lib/components/Navbar.svelte'
  import Footer from '$lib/components/Footer.svelte'
  import RepoList from './RepoList.svelte'
  import SwitchMenu from '$lib/components/SwitchMenu.svelte'
	
	import { compare } from '$lib/components/sharedStore'
	import { localRepoStore, localSearchStore, localSearchFilter } from './store'
	import type { RepoListResponse } from './structure'
	import { toLower } from '$lib/components/functions'

  let search = ""
  let files: FileList
  let repoDetail: RepoListResponse

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
      window.alert(`[Success] ${filename} has been uploaded. Page will be reload to take effect.`)
      window.location.reload()
    }
  }

  export let data
  const { localRepo } = data

  $: localSearchStore.set(toLower(search).split(/\s+/))
  $: localRepoStore.set(localRepo)
  $: selected = $compare;
</script>

<svelte:head>
	<title>NSP YANG Browser | Uploads</title>
</svelte:head>

<Navbar nspIp=""/>
<div class="font-nunito text-sm container mx-auto pt-[75px]">
  <div class="px-6 py-4">
    <div class="flex items-center justify-between">
      <p class="text-lg pb-2 text-black dark:text-white">Uploads</p>
      <SwitchMenu isUpload={true} />
    </div>
    <div>
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
    <div class="py-2">
      <input type="text" placeholder="Search..." bind:value={search} 
        class="px-3 py-2 rounded-lg w-full text-[12.5px] text-gray-800 dark:text-gray-200 
          dark:placeholder-gray-400 border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700">
      {#if $localRepoStore.length > 0}
        <ul class="mt-2 mb-2 border-y dark:border-gray-700 overflow-y-auto scroll-light dark:scroll-dark max-h-[495px]">
          {#each $localSearchFilter as {name, files}, i}
            {#if name === "" || name.includes("from-nsp")}
              {@const sharedName = "common .yang files"}
              <li class="text-sm bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 {i > 0 ? 'border-t dark:border-gray-700' : ''}">
                <div class="flex items-center justify-between">
                  <p class="px-4 py-3 w-full overflow-x-auto">{sharedName}</p>
                  <div class="flex items-center mx-4 space-x-5">
                    <button class="text-gray-600 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 rounded-lg p-1" on:click={() => repoDetail = {name: sharedName, files: files ?? []}}>
                      <svg class="w-5 h-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                        <path stroke="currentColor" stroke-linecap="round" stroke-width="2" d="M9 8h10M9 12h10M9 16h10M4.99 8H5m-.02 4h.01m0 4H5"/>
                      </svg>
                    </button>
                  </div>
                </div>
              </li>
            {:else}
              {@const compareValue = "uploaded@" + name}
              {@const isDisabled = !selected.includes(compareValue) && selected.length >= 2}
              <li class="text-sm bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 {i > 0 ? 'border-t dark:border-gray-700' : ''} hover:bg-gray-100">
                <div class="flex items-center justify-between">
                  <a data-sveltekit-reload href="/uploaded/{name}" class="px-4 py-3 w-full overflow-x-auto">{name}</a>
                  <div class="flex items-center mx-4 space-x-5">
                    <div title="Add to compare" class="flex">
                      <input type="checkbox" id="uploaded-{name}-check" class="peer hidden" 
                        disabled={isDisabled}
                        checked={selected.includes(compareValue)} 
                        on:change={(e) => e.currentTarget.checked ? compare.add(compareValue) : compare.remove(compareValue)}
                      />
                      <label for="uploaded-{name}-check" class="select-none p-1 rounded-lg peer-checked:bg-blue-600 peer-checked:hover:bg-blue-700 peer-checked:text-white {isDisabled ? 'cursor-not-allowed text-gray-200 dark:text-gray-600' : 'cursor-pointer hover:bg-blue-600 hover:text-white'}">
                        <svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024" fill="currentColor" stroke="currentColor" stroke-width="10" aria-hidden="true">
                          <path d="M420.266667 832c-17.066667 0-34.133333-6.4-44.8-19.2L104.533333 541.866667c-12.8-12.8-19.2-27.733333-19.2-44.8s6.4-34.133333 19.2-44.8L345.6 211.2c23.466667-23.466667 66.133333-23.466667 89.6 0l270.933333 270.933333c12.8 12.8 19.2 27.733333 19.2 44.8s-6.4 34.133333-19.2 44.8L465.066667 812.8c-10.666667 12.8-27.733333 19.2-44.8 19.2z m-29.866667-597.333333c-6.4 0-10.666667 2.133333-14.933333 6.4L134.4 482.133333c-4.266667 4.266667-6.4 8.533333-6.4 14.933334s2.133333 10.666667 6.4 14.933333L405.333333 782.933333c8.533333 8.533333 21.333333 8.533333 29.866667 0l241.066667-241.066666c4.266667-4.266667 6.4-8.533333 6.4-14.933334s-2.133333-10.666667-6.4-14.933333L405.333333 241.066667c-4.266667-4.266667-8.533333-6.4-14.933333-6.4z" />
                          <path d="M618.666667 832c-17.066667 0-34.133333-6.4-46.933334-19.2L317.866667 558.933333c-12.8-12.8-19.2-29.866667-19.2-46.933333s6.4-34.133333 19.2-46.933333L571.733333 211.2c25.6-25.6 68.266667-25.6 93.866667 0l253.866667 253.866667c25.6 25.6 25.6 68.266667 0 93.866666L665.6 812.8c-12.8 12.8-29.866667 19.2-46.933333 19.2z m0-597.333333c-6.4 0-12.8 2.133333-17.066667 6.4L347.733333 494.933333c-4.266667 4.266667-6.4 10.666667-6.4 17.066667s2.133333 12.8 6.4 17.066667l253.866667 253.866666c8.533333 8.533333 23.466667 8.533333 34.133333 0l253.866667-253.866666c8.533333-8.533333 8.533333-23.466667 0-34.133334L635.733333 241.066667c-4.266667-4.266667-10.666667-6.4-17.066666-6.4zM332.8 480z" />
                        </svg>
                      </label>
                    </div>
                    <button class="text-gray-600 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-800 rounded-lg p-1" on:click={() => repoDetail = {name, files}}>
                      <svg class="w-5 h-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                        <path stroke="currentColor" stroke-linecap="round" stroke-width="2" d="M9 8h10M9 12h10M9 16h10M4.99 8H5m-.02 4h.01m0 4H5"/>
                      </svg>
                    </button>
                  </div>
                </div>
              </li>
            {/if}
          {/each}
        </ul>
      {:else}
        <p class="py-2 text-red-600">No records found</p>
      {/if}
    </div>
  </div>
  <RepoList bind:repoDetail />
  <Footer home={false}/>
</div>