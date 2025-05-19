<script lang="ts">
  import { stateValues } from '$lib/components/sharedStore'

  export let stateInput: string[]
</script>

<div class="dropdown">
  <button class="dropdown-button px-2 py-1 text-xs text-center text-nowrap rounded-lg inline-flex items-center dark:text-white 
    {stateInput?.length ? 'bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600' : ''} border border-gray-200 dark:border-gray-600">
      Filters
      <svg class="w-2.5 h-2.5 ms-3" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 10 6">
        <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m1 1 4 4 4-4"/>
      </svg>
  </button>
  <div class="dropdown-content absolute z-10 hidden dark:text-white rounded-lg shadow
    bg-gray-100 dark:bg-gray-700 border border-gray-200 dark:border-gray-600">
    <div class="my-1 overflow-y-auto scroll-light dark:scroll-dark">
      <ul>
        {#each stateValues as entry}
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
          <li class="px-4 py-2 flex items-center text-xs hover:bg-gray-200 dark:hover:bg-gray-600 cursor-pointer">
            <input id="{entry.label}-checkbox" type="checkbox" class="w-3 h-3" value="{entry.value}" bind:group={stateInput}>
            <label for="{entry.label}-checkbox" class="ms-2 text-xs text-nowrap text-gray-900 dark:text-gray-300 cursor-pointer w-full">{entry.label}</label>
          </li>
        {/each}
      </ul>
    </div>
  </div>
</div>