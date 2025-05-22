<script lang="ts">
  export let kind: string
  export let basename: string
  export let isTree: boolean
</script>

<a data-sveltekit-reload  
  href="{isTree ? `/${kind}/${basename}` : `/${kind}/${basename}/tree`}" 
  class="px-2 py-1 rounded-lg cursor-pointer text-xs text-nowrap text-white 
  bg-blue-600 dark:bg-blue-700 hover:bg-blue-700 dark:hover:bg-blue-600">
  {isTree ? 'Path' : 'Tree'} Browser
</a>
