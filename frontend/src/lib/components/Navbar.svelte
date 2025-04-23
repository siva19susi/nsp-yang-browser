<script lang="ts">
	import Theme from "$lib/components/Theme.svelte"
	import Compare from "./Compare.svelte"

	import { compare } from "./sharedStore"

  export let kind: string = ""
  export let basename: string = ""
  export let nspIp: string

  let visualiseCompare = false
  const isNspUrl = kind.includes("nsp") || nspIp != ""

  async function nspDisconnect() {
    if (window.confirm("Are you sure you want to disconnect from NSP?")) {
      const nspDisconnected = await fetch("/api/nsp/disconnect", {method: "POST"})
      if(nspDisconnected.ok) {
        compare.clear()
        window.alert(`[Success] NSP disconnected. Redirecting to home page.`)
        window.location.href = "/"
      } else {
        window.alert(`[Error] Failed to disconnect NSP. Page will reload to take effect.`)
      }
    }
  }
</script>

<nav class="fixed top-0 z-20 px-6 py-4 w-screen font-nunito text-black dark:text-white backdrop-filter backdrop-blur-lg bg-opacity-50 border-b dark:border-gray-700">
  <div class="flex justify-between">
    <div class="flex items-center space-x-4">
      <a href="/"><img src="/images/navbar-logo.png" alt="Logo" width="25"/></a>
      <div class="flex flex-col whitespace-nowrap overflow-x-auto scroll-light dark:scroll-dark w-44 sm:w-fit">
        <div class="flex items-center space-x-1 {kind != "" ? 'text-xs text-gray-500 dark:text-gray-400' : 'text-sm'}">
          <p class="">NSP YANG Browser</p>
          {#if isNspUrl}
            <span>|</span>
            {#if kind.includes(";")}
              <span>Compare</span>
            {:else}
              <button class="text-blue-700 dark:text-blue-400 hover:underline" on:click={nspDisconnect}>{nspIp}</button>
            {/if}
          {/if}
        </div>
        {#if kind != ""}
          {#if kind.includes(";")}
            {@const [xKind, yKind] = kind.split(";")}
            {@const [xBasename, yBasename] = basename.split(";")}
            <div class="text-gray-800 text-xs lg:text-sm dark:text-white">
              <div class="flex flex-wrap items-center justify-center space-x-1">
                <span class="dropdown">
                  <button class="dropdown-button underline">{yBasename} ({yKind})</button>
                  <div class="dropdown-content absolute z-10 hidden bg-gray-100 dark:bg-gray-700 dark:text-white rounded-lg shadow">
                    <p class="my-2 max-w-[200px] px-2 text-xs text-wrap">
                      Changes and filters shown are with respect to this selection.
                    </p>
                  </div>
                </span>
                <span>with</span>
                <span>{xBasename} ({xKind})</span>
              </div>
            </div>
          {:else}
            <p class="text-sm">{basename} ({kind})</p>
          {/if}
        {/if}
      </div>
    </div>
    <div class="flex items-center">
      {#if !kind.includes(";")}
        <button class="inline-block cursor-pointer {$compare.length ? 'animate-pulse text-blue-600 hover:text-blue-800 dark:text-blue-400 hover:dark:text-blue-600' : 'text-gray-600 dark:text-gray-400'}" on:click={() => visualiseCompare = !visualiseCompare}>
          <svg class="w-6 h-6" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024" fill="currentColor" stroke="currentColor" stroke-width="10" aria-hidden="true">
            <path d="M420.266667 832c-17.066667 0-34.133333-6.4-44.8-19.2L104.533333 541.866667c-12.8-12.8-19.2-27.733333-19.2-44.8s6.4-34.133333 19.2-44.8L345.6 211.2c23.466667-23.466667 66.133333-23.466667 89.6 0l270.933333 270.933333c12.8 12.8 19.2 27.733333 19.2 44.8s-6.4 34.133333-19.2 44.8L465.066667 812.8c-10.666667 12.8-27.733333 19.2-44.8 19.2z m-29.866667-597.333333c-6.4 0-10.666667 2.133333-14.933333 6.4L134.4 482.133333c-4.266667 4.266667-6.4 8.533333-6.4 14.933334s2.133333 10.666667 6.4 14.933333L405.333333 782.933333c8.533333 8.533333 21.333333 8.533333 29.866667 0l241.066667-241.066666c4.266667-4.266667 6.4-8.533333 6.4-14.933334s-2.133333-10.666667-6.4-14.933333L405.333333 241.066667c-4.266667-4.266667-8.533333-6.4-14.933333-6.4z" />
            <path d="M618.666667 832c-17.066667 0-34.133333-6.4-46.933334-19.2L317.866667 558.933333c-12.8-12.8-19.2-29.866667-19.2-46.933333s6.4-34.133333 19.2-46.933333L571.733333 211.2c25.6-25.6 68.266667-25.6 93.866667 0l253.866667 253.866667c25.6 25.6 25.6 68.266667 0 93.866666L665.6 812.8c-12.8 12.8-29.866667 19.2-46.933333 19.2z m0-597.333333c-6.4 0-12.8 2.133333-17.066667 6.4L347.733333 494.933333c-4.266667 4.266667-6.4 10.666667-6.4 17.066667s2.133333 12.8 6.4 17.066667l253.866667 253.866666c8.533333 8.533333 23.466667 8.533333 34.133333 0l253.866667-253.866666c8.533333-8.533333 8.533333-23.466667 0-34.133334L635.733333 241.066667c-4.266667-4.266667-10.666667-6.4-17.066666-6.4zM332.8 480z" />
          </svg>
        </button>
      {/if}
      <Theme/>
    </div>
  </div>
</nav>
{#if !kind.includes(";")}
  <Compare bind:visualiseCompare/>
{/if}