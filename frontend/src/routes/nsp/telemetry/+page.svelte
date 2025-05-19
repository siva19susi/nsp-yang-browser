<script lang="ts">
  import { onMount } from 'svelte'

  import Navbar from '$lib/components/Navbar.svelte'
  import Footer from '$lib/components/Footer.svelte'
  import SearchInput from '$lib/components/SearchInput.svelte'
  import Pagination from './Pagination.svelte'

  import { markFilter, markRender, toLower } from '$lib/components/functions'
	import { definitionStore, paginated, searchStore, total } from './store'

  // DEFAULTS
  let searchInput = ""

  // ON PAGELOAD
  export let data
  let {type, definition, nspIp} = data
  const nspConnected = (nspIp != "" ? true : false)

  // OTHER BINDING VARIABLES
  $: definitionStore.set(definition)
  $: searchStore.set(toLower(searchInput))
</script>

<svelte:head>
	<title>NSP YANG Path Browser | Telemetry Type Definition</title>
</svelte:head>

<Navbar kind="telemetry" basename={type} {nspIp}/>
<div class="min-w-[280px] overflow-x-auto font-nunito dark:bg-gray-800 pt-[75px]">
  <div class="px-6 pt-6 container mx-auto">
    <div class="flex items-center justify-between">
      <p class="text-gray-800 dark:text-gray-300">Telemetry Type Definition</p>
    </div>
    <SearchInput bind:searchInput />
    <Pagination />
    <div class="overflow-x-auto rounded-t-lg max-w-full mt-2">
      <table class="text-left w-full">
        <colgroup>
          <col span="1" class="w-[25%]">
          <col span="1" class="w-[15%]">
          <col span="1" class="w-[60%]">
        </colgroup>
        <thead class="text-sm text-gray-800 dark:text-gray-300 bg-gray-300 dark:bg-gray-700">
          <tr>
            <th scope="col" class="px-3 py-2">Counter</th>
            <th scope="col" class="px-3 py-2">Type</th>
            <th scope="col" class="px-3 py-2">Device Path</th>
          </tr>
        </thead>
        <tbody>
          {#if $total > 0}
            {#each $paginated as item}
              <tr class="bg-white dark:bg-gray-800 border-b dark:border-gray-700 text-gray-700 dark:text-gray-300 hover:cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-600">
                <td class="px-3 py-1.5 font-fira text-[13px] tracking-tight"><div use:markRender={item.counterName}></div></td>
                <td class="px-3 py-1.5 font-fira text-[13px] tracking-tight group"><div use:markRender={item.dataType}></div></td>
                <td class="px-3 py-1.5 font-fira text-[13px] tracking-tight"><div use:markRender={item.deviceXpath}></td>
              </tr>
            {/each}
          {:else}
            <tr class="bg-white border-b dark:bg-gray-800 dark:border-gray-700">
              <td colspan="3" class="px-3 py-1.5 font-fira text-[13px] text-red-600 text-center">No results found</td>
            </tr>
          {/if}
        </tbody>
      </table>
    </div>
    <Pagination />
  </div>
  <Footer home={false}/>
</div>