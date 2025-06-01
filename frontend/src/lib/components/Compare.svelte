<script lang="ts">
  import { compare } from "./sharedStore"

  export let visualiseCompare = false

  function closeCompare() {
    if(visualiseCompare) {
      visualiseCompare = false
    }
  }

  function fetchCompareKey(key: string, fetch: string) {
    const [kind, basename] = key.split("@")
    if(kind === "offline") {
      const [id, nspIp, module, name ] = basename.split("__")
      if(fetch === "kind") {
        return `offline (${nspIp})`
      } else if(fetch === "basename") {
        return (module === "telemetry-type" ? "/" + name.replaceAll("_", "/") : name) + ` (${module})`
      }
    } else {
      if(fetch === "kind") {
        return kind
      } else if(fetch === "basename") {
        return basename
      }
    }
  }

  function fetchCompareId(key: string) {
    console.log(key)
    const [kind, basename] = key.split("@")
    if(kind === "offline") {
      const [id, nspIp, module, name ] = basename.split("__")
      return `${kind}@${id}`
    } else {
      return key
    }
  }
</script>

<svelte:window on:keyup={({key}) => key === "Escape" ? closeCompare() : ""}/>

<div id="comparePopup" class="fixed px-6 py-4 inset-0 z-50 items-center font-nunito { visualiseCompare  ? '' : 'hidden'}">
  <div class="fixed inset-0 bg-gray-800 bg-opacity-75 transition-opacity"></div>
  <div id="popupContent" class="flex min-h-full justify-center items-center">
    <div class="relative transform overflow-hidden rounded-lg bg-white dark:bg-gray-700 text-left shadow-xl transition-all sm:my-8 max-w-4xl">
      <div id="popupHeader" class="flex items-center justify-between space-x-2 px-4 py-2 rounded-t bg-gray-200 dark:bg-gray-600 border-b border-gray-200 dark:border-gray-600">
        <div class="flex items-center space-x-2">
          <span class="text-lg text-gray-900 dark:text-gray-300">Compare</span>
        </div>
        <button type="button" class="text-gray-500 hover:bg-gray-300 hover:text-gray-900 rounded-lg text-sm inline-flex justify-center items-center dark:hover:bg-gray-700 dark:hover:text-white" on:click={closeCompare}>
          <svg class="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"/>
          </svg>
          <span class="sr-only">Close modal</span>
        </button>
      </div>
      <div id="comparePopupBody">
        <div class="overflow-auto scroll-light dark:scroll-dark">
          <div class="flex flex-col h-full">
            <div class="flex-grow bg-gray-100 dark:bg-gray-800">
              <div class="flex items-center p-4 text-sm space-x-3 select-none h-full">
                <div class="w-full text-center">
                  <p class="px-3 py-20 rounded-t-lg border-x-2 border-t-2 border-dashed border-gray-400 dark:border-gray-500 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200">
                    {$compare.length > 0 ? fetchCompareKey($compare[0], "basename") : 'X'}
                  </p>
                  <p class="p-3 {$compare.length > 0 ? 'bg-gray-600' : ''} bg-gray-400 dark:bg-gray-500 text-white rounded-b-lg">
                    {$compare.length > 0 ? fetchCompareKey($compare[0], "kind") : '{{ source }}'}
                  </p>
                </div>
                <p class="rounded-full px-3 py-2 bg-blue-600 text-white">to</p>
                <div class="w-full text-center">
                  <p class="px-3 py-20 rounded-t-lg border-x-2 border-t-2 border-dashed border-gray-400 dark:border-gray-500 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200">
                    {$compare.length > 1 ? fetchCompareKey($compare[1], "basename") : 'Y'}
                  </p>
                  <p class="p-3 {$compare.length > 1 ? 'bg-gray-600' : ''} bg-gray-400 dark:bg-gray-500 text-white rounded-b-lg">
                    {$compare.length > 1 ? fetchCompareKey($compare[1], "kind") : '{{ source }}'}
                  </p>
                </div>
              </div>
            </div>
            <div class="text-center px-4 py-4 bg-gray-50 dark:bg-gray-700 border-t dark:border-gray-700 rounded-b-lg">
              <div class="flex items-center justify-between">
                <button class="px-3 py-1 text-sm bg-blue-200 rounded-lg" on:click={() => compare.clear()}>Reset</button>
                <a href="/compare/{$compare.length === 2 ? `${fetchCompareId($compare[0])}..${fetchCompareId($compare[1])}` : '' }" 
                  class="px-4 py-1 rounded-lg text-sm text-white select-none 
                  {$compare.length === 2 ? 
                    'bg-green-600 hover:bg-green-700' : 
                    'bg-gray-300 dark:bg-gray-600 text-gray-500 dark:text-gray-500 pointer-events-none'}">
                  Click to Compare
                </a>
              </div>
              <div class="pt-6 text-[12px] text-left text-gray-800 dark:text-gray-200">
                <p>Note:</p>
                <ul class="px-4 list-disc list-outside">
                  <li>Order or compare selection determines the value of X (first selected) and Y (second selected).</li>
                  <li>The compare provides added paths, removed paths and paths with modified type definition.</li>
                  <li>For more detailed info regarding each repo, contact the repo source Admin.</li>
                  <li>Disconnecting NSP or deleting the selected repo (uploaded) clears the compare selection.</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>