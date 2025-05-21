<script lang="ts">
  import Navbar from '$lib/components/Navbar.svelte'
  import Footer from '$lib/components/Footer.svelte'
  import SearchInput from '$lib/components/SearchInput.svelte'
  import Pagination from './Pagination.svelte'
  import EnableOfflineMode from '$lib/components/EnableOfflineMode.svelte'

  import { markRender, toLower } from '$lib/components/functions'
	import { definitionStore, paginated, searchStore, total } from './store'

  // DEFAULTS
  let searchInput = ""

  // ON PAGELOAD
  export let data
  let {kind, type, definition, nspIp} = data
  const nspConnected = (nspIp != "" ? true : false)

  // OTHER BINDING VARIABLES
  $: definitionStore.set(definition)
  $: searchStore.set(toLower(searchInput))
</script>

<svelte:head>
	<title>NSP YANG Path Browser | Telemetry Type Definition</title>
</svelte:head>

<Navbar {kind} basename={type} {nspIp}/>
<div class="min-w-[280px] overflow-x-auto font-nunito dark:bg-gray-800 pt-[75px]">
  <div class="px-6 pt-6 container mx-auto">
    <div class="flex items-end justify-between">
      <p class="text-gray-800 dark:text-gray-300">Telemetry Type Definition</p>
      {#if kind !== "offline"}
        <EnableOfflineMode kind="telemetry" basename={type} />
      {/if}
    </div>
    <SearchInput bind:searchInput />
    <div class="px-1 py-1.5 text-xs text-gray-900 dark:text-gray-300">
      <p>
        Creating the object filter for NSP telemetry subscription can be framed  
        by prepending the below NE identifier to the Device Path:
      </p>
      <p class="font-fira text-[11px] tracking-tight pt-1">
        /network-device-mgr:network-devices/network-device[name='IPv4/IPv6-NE-ID']/root
      </p>
    </div>
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
              <tr class="bg-white dark:bg-gray-800 border-b dark:border-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600">
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