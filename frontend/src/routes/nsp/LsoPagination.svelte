<script lang="ts">
	import { lsoEnd, lsoStart, lsoTotal } from "./store";

  function updateTable(s: number) {
    if(s >= 0 && s < $lsoTotal) {
      lsoStart.set(s)
    }
  }
</script>

<div class="flex items-center justify-end py-3 text-sm mt-2">
  <p class="mr-2 text-gray-800 dark:text-gray-200">{$lsoStart + 1} - {$lsoEnd > 1 ? $lsoEnd : 0} of {$lsoTotal}</p>
  <button class="ml-2 text-white rounded 
      {$lsoStart == 0 ? 'bg-gray-300 dark:bg-gray-500 opacity-50 cursor-not-allowed' : 'bg-gray-400 hover:bg-gray-600'}" 
      disabled="{$lsoStart == 0}" on:click={() => updateTable($lsoStart - 30)}>
    <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
      <path fill-rule="evenodd" d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clip-rule="evenodd"/>
    </svg>
  </button>
  <button class="ml-2 text-white rounded 
      {$lsoEnd == $lsoTotal ? 'bg-gray-300 dark:bg-gray-500 opacity-50 cursor-not-allowed' : 'bg-gray-400 hover:bg-gray-600'}" 
      disabled="{$lsoEnd == $lsoTotal}" on:click={() => updateTable($lsoEnd)}>
    <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
      <path fill-rule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clip-rule="evenodd"/>
    </svg>
  </button>
</div>