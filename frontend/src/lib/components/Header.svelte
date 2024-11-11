<script lang="ts">
	import { page } from '$app/stores'
  import { kindView } from '$lib/components/functions'
  
  import Theme from '$lib/components/Theme.svelte'

  export let kind: string
  export let basename: string
</script>

<!-- NAVBAR -->
<nav class="fixed top-0 z-20 p-4 w-screen select-none font-nunito bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
	<div class="flex justify-between">
		<!-- navbar left item -->
    <div class="flex items-center space-x-2">
      <a href="{$page.url.pathname.includes("tree") ? '../../' : '../'}" class="px-2"><img src="/images/navbar-logo.png" alt="Logo" width="25"/></a>
    </div>
		<!-- navbar centre item -->
    <div class="text-center">
      {#if basename.includes(";")}
        {@const [xKind, yKind] = kind.split(";")}
        {@const [xBasename, yBasename] = basename.split(";")}
        <p class="text-nokia-old-blue dark:text-white font-light text-lg lg:text-2xl">Yang Compare</p>
        <div class="text-gray-800 text-xs lg:text-sm dark:text-white">
          <div class="flex flex-wrap items-center justify-center space-x-1">
            <span class="dropdown">
              <button class="dropdown-button font-nokia-headline underline">{yBasename} ({kindView(yKind)})</button>
              <div class="dropdown-content absolute z-10 hidden bg-gray-100 dark:bg-gray-700 dark:text-white rounded-lg shadow">
                <p class="my-2 max-w-[200px] px-1 text-xs">
                  Changes and filters shown are with respect to this release
                </p>
              </div>
            </span>
            <span>with</span>
            <span class="font-nokia-headline">{xBasename} ({kindView(xKind)})</span>
          </div>
        </div>
      {:else}
      <p class="text-nokia-old-blue dark:text-white font-light text-lg lg:text-2xl">Yang Browser</p>
      <p class="text-gray-800 text-xs lg:text-sm dark:text-white">{basename} ({kindView(kind)})</p>
      {/if}
    </div>
		<!-- navbar right item -->
    <div class="flex items-center">
		  <Theme/>
    </div>
	</div>
</nav>