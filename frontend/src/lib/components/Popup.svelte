<script lang="ts">
  import { page } from '$app/stores'
	import { copy } from 'svelte-copy'

  import { markRender } from '$lib/components/functions'
	import { stateValues } from './sharedStore'

  export let kind: string = ""
  export let basename: string = ""
  export let popupDetail: any = {}
  
  const enableQueryNsp = () => {
    if(kind !== "uploaded") {
      return popupDetail.nspConnected && !popupDetail["is-rpc"] && !popupDetail["is-notification"] && !popupDetail["is-action"]
    }
    return false
  }

  function closePopup() {
    if(Object.keys(popupDetail).length !== 0) {
      popupDetail = {}
    }
  }

  function closeSidebarPopup(event: any) {
    if(!document.getElementById("popupContent")?.contains(event.target)) {
      closePopup()
    }
  }

  function crossLaunch(path: any) {
    const toTree = (popupDetail.isUrlTree ? "" : `/tree`)
    const fromParam = (popupDetail.isUrlTree ? "" : "&from=pb")
    if(kind === "") kind = path.kind
    if(basename === "") basename = path.basename
    return `/${kind}/${basename}${toTree}?path=${encodeURIComponent(path.path)}${fromParam}`
  }

  function queryNsp(path: any) {
    const urlPath = path["path-with-prefix"]

    const removeListKey = (path: string) => {
      // Remove all [key=value] style segments
      return path.replace(/\[[^\]]*=[^\]]*\]/g, '')
    }
    const removeLastSegment = (path: string) => {
      const segments = path.split('/')

      if (segments.length > 2) { // more than one "/"
        const removed = segments.pop() // remove last
        const newPath = segments.join('/')
        return { newPath, removedSegment: removed }
      }

      return { newPath: path, removedSegment: null }
    }

    const { newPath, removedSegment } = removeLastSegment(urlPath)

    if(kind === "nsp-intent-type") {
      return `/${kind}/${basename}/intent-query?path=${encodeURIComponent(newPath)}`
    }
    return `/${kind}/${basename}/query?path=${encodeURIComponent(removeListKey(newPath))}`
    //return `/${kind}/${basename}/query?path=${encodeURIComponent(removeListKey(newPath))}&field=${removedSegment}`
  }
</script>

<svelte:window on:keyup={({key}) => key === "Escape" ? closePopup() : ""}/>

{#if Object.keys(popupDetail).length !== 0}
  <div id="popup" class="fixed p-4 inset-0 z-50 items-center { Object.keys(popupDetail).length !== 0  ? '' : 'hidden'}">
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div class="fixed inset-0 bg-gray-800 bg-opacity-75 transition-opacity" on:click|stopPropagation={closeSidebarPopup}></div>
    <div id="popupContent" class="flex min-h-full justify-center items-center">
      <div class="relative transform overflow-hidden rounded-lg bg-white dark:bg-gray-700 text-left shadow-xl transition-all sm:my-8 max-w-4xl">
        <div id="popupHeader" class="flex items-center justify-between px-4 py-2 rounded-t bg-gray-200 dark:bg-gray-600 border-b border-gray-200 dark:border-gray-600">
          <div class="flex items-center">
            <span class="text-lg text-gray-900 dark:text-gray-300">Path Details</span>
          </div>
          <button type="button" class="text-gray-500 hover:bg-gray-300 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center dark:hover:bg-gray-700 dark:hover:text-white" on:click={closePopup}>
            <svg class="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
                <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"/>
            </svg>
            <span class="sr-only">Close modal</span>
          </button>
        </div>
        <div id="popupBody" class="p-4 text-left">
          <div class="overflow-x-auto max-w-full">
            <table class="w-full">
              <tbody>
                {#if "compare" in popupDetail}
                  <tr>
                    {#if popupDetail.compare === "~"}
                      <td colspan="2" class="pt-1 pb-3 text-sm text-gray-400 dark:text-gray-400">MODIFIED in {popupDetail.compareTo}</td>
                    {:else if popupDetail.compare === "+"}
                      <td colspan="2" class="pt-1 pb-3 text-sm text-green-600 dark:text-green-300">PRESENT in {popupDetail.compareTo}</td>
                    {:else if popupDetail.compare === "-"}
                      <td colspan="2" class="pt-1 pb-3 text-sm text-red-600 dark:text-red-300">NOT PRESENT in {popupDetail.compareTo}</td>
                    {/if}
                  </tr>
                {/if}
                <tr>
                  <th scope="row" class="py-1 whitespace-nowrap text-sm dark:text-gray-400">Data:</th>
                  <td class="py-1 px-2 dark:text-gray-300 font-fira text-[13px] tracking-tight">{stateValues.filter(x => x.value == popupDetail["added-filter"])[0].label}</td>
                </tr>
                <tr class="border-t border-gray-200 dark:border-gray-600">
                  <th scope="row" class="py-1 whitespace-nowrap text-sm dark:text-gray-400">Type:</th>
                  <td class="py-1 px-2 dark:text-gray-300 font-fira text-[13px] tracking-tight">
                    {#if "fromType" in popupDetail}
                      <p class="text-gray-400">from ({popupDetail.fromRel}): {popupDetail.fromType}</p>
                    {/if}
                    <p>{popupDetail.type}</p>
                  </td>
                </tr>
                {#if popupDetail["type"] === "enumeration" && "enum-values" in popupDetail}
                  <tr class="border-t border-gray-200 dark:border-gray-600">
                    <th scope="row" class="py-1 whitespace-nowrap text-sm dark:text-gray-400">Enum Values:</th>
                    <td class="py-1 px-2 dark:text-gray-300 font-fira text-[13px] tracking-tight">{popupDetail["enum-values"].join(", ")}</td>
                  </tr>
                {/if}
                {#if "default" in popupDetail}
                  <tr class="border-t border-gray-200 dark:border-gray-600">
                    <th scope="row" class="py-1 whitespace-nowrap text-sm dark:text-gray-400">Default:</th>
                    <td class="py-1 px-2 dark:text-gray-300 font-fira text-[13px] tracking-tight">{popupDetail["default"]}</td>
                  </tr>
                {/if}
                <tr class="border-t border-gray-200 dark:border-gray-600">
                  <th scope="row" class="py-1 whitespace-nowrap text-sm dark:text-gray-400">Path:</th>
                  <td class="py-1 px-2 dark:text-gray-300 font-fira text-[13px] tracking-tight">{popupDetail.path}</td>
                </tr>
                <tr class="border-t border-gray-200 dark:border-gray-600">
                  <th scope="row" class="py-1 whitespace-nowrap text-sm dark:text-gray-400">Path with Prefix:</th>
                  <td class="py-1 px-2 dark:text-gray-300 font-fira text-[13px] tracking-tight">{popupDetail["path-with-prefix"]}</td>
                </tr>
                <tr class="border-t border-gray-200 dark:border-gray-600">
                  <th scope="row" class="py-1 whitespace-nowrap text-sm dark:text-gray-400">Description:</th>
                  <td class="py-1 px-2 dark:text-gray-300 font-fira text-[13px] tracking-tight">
                    <div class="overflow-y-auto max-h-40 scroll-light dark:scroll-dark" use:markRender={"description" in popupDetail ? popupDetail.description.replaceAll("\n", "<br>") : ''}></div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div id="popupFooter" class="flex items-center {enableQueryNsp() ? 'justify-between space-x-2' : 'justify-end'} px-4 py-2 border-t border-gray-200 rounded-b dark:border-gray-600">
          {#if enableQueryNsp()}
            <a href="{queryNsp(popupDetail)}" class="text-sm px-3 py-1 rounded text-white bg-gray-500 hover:bg-gray-600 dark:bg-gray-600 dark:hover:bg-gray-800">
              Query NSP
            </a>
          {/if}
          <a href="{crossLaunch(popupDetail)}" class="text-sm px-2 py-1 rounded text-white bg-blue-500 hover:bg-blue-600 dark:bg-blue-600 dark:hover:bg-blue-700">
            {popupDetail.isUrlTree ? 'Path' : 'Tree'} Browser View
          </a>
        </div>
      </div>
    </div>
  </div>
{/if}