<script lang="ts">
  export let pathWithDefault: boolean
</script>

<div class="px-2 py-1 flex items-center border hover:border-gray-300 dark:border-gray-500 dark:hover:border-gray-400 rounded-lg">
  <input id="default-checkbox" type="checkbox" class="w-3 h-3" bind:checked={pathWithDefault}>
  <label for="default-checkbox" class="ms-2 text-xs text-nowrap text-gray-900 dark:text-gray-300 cursor-pointer">Has defaults</label>
</div>