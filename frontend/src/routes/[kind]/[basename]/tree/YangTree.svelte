<script lang="ts">
	import { page } from "$app/stores"
	import { goto } from "$app/navigation"
	
	import { pathFocus } from '$lib/components/sharedStore'
	import { markFilter, markRender } from "$lib/components/functions"
	import { decideExpand } from "./expand"

	export let folder
	export let expanded: boolean
	export let withPrefix: boolean

	const toggle = () => expanded = !expanded
	const isCrossLaunched = () => $page.data.crossLaunched
  const getUrlPath = () => $page.data.urlPath
  const isNspConnected = () => $page.data.nspIp != "" ? true : false
  const isUrlTree = () => $page.data.isUrlTree

	function leafClick(details: any) {
		pathFocus.set({
      ...details, isUrlTree: isUrlTree(), 
      nspConnected: isNspConnected()
    })
		$page.url.searchParams.delete("from")
		$page.url.searchParams.set("path", details.path)
		goto(`?${$page.url.searchParams.toString()}`, {noScroll: true})
	}
</script>

<div class="flex items-center space-x-2">
	<button class="flex items-center text-left py-0.5 
			{expanded ? 'text-gray-400 dark:text-gray-500': 'dark:text-gray-300'} 
			hover:text-green-600 hover:dark:text-green-600" on:click={toggle}>
		<span class="flex pr-1">
			<!--minus-circle-->
			<svg class="w-5 h-5 {expanded ? 'minus-circle-active': 'hidden'}" fill="none" 
				stroke-width="1.5" stroke="currentColor" viewBox="0 0 25 25" xmlns="http://www.w3.org/2000/svg">
				<path stroke-linecap="round" stroke-linejoin="round" d="M15 12H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z"/>
			</svg>
			<!--plus-circle-->
			<svg class="w-5 h-5 {expanded ? 'plus-circle-inactive hidden': ''}" fill="none" 
				stroke-width="1.5" stroke="currentColor" viewBox="0 0 25 25" xmlns="http://www.w3.org/2000/svg">
				<path stroke-linecap="round" stroke-linejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z"/>
			</svg>
		</span>
		<div title="{folder.details.path}" use:markRender={markFilter(folder.name, getUrlPath(), "tree")}></div>
	</button>
	<a href="/{$page.data.kind}/{$page.data.basename}/tree/payload?path={folder.details.path}{withPrefix ? '&prefix=true' : ''}&expand=false"  
		class="dropdown-button px-2 py-0.5 rounded-lg text-[10px] text-nowrap bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 text-black dark:text-white inline-flex items-center align-bottom">
		Payload
	</a>
</div>

{#if expanded}
	{#if folder.children && folder.children?.length}
		<ul class="ml-2.5 px-2 list-none border-l dark:border-gray-300">
			{#each folder.children as entry}
				{@const urlPath = getUrlPath()}
				{@const crossLaunched = isCrossLaunched()}
				<li class="pt-1">
					{#if entry.children.length > 0}
						<svelte:self folder={entry} {withPrefix} expanded={decideExpand(entry, crossLaunched, urlPath)} />
					{:else}
						<button class="ml-2.5 px-2 py-0.5 rounded hover:underline 
							hover:bg-gray-200 hover:text-black hover:dark:bg-gray-600 hover:dark:text-gray-200 
							{urlPath === entry.details.path ? 'bg-gray-200 dark:bg-gray-600 dark:text-gray-200' : 'text-blue-600 dark:text-blue-500'}" 
							on:click={() => leafClick(entry.details)}>
							{#if urlPath === entry.details.path}
                <div title="{entry.details.path}">{entry.name}</div>
              {:else}
                <div title="{withPrefix ? entry.details["path-with-prefix"] : entry.details.path}" use:markRender={markFilter(entry.name, urlPath, "tree")}></div>
              {/if}
						</button>
					{/if}
				</li>
			{/each}
		</ul>
	{/if}
{/if}