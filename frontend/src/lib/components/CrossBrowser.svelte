<script lang="ts">
  export let kind: string
  export let basename: string
  export let isTree: boolean

  const newBasename = kind === "offline" ? basename.split("__")[0] : basename
</script>

<a data-sveltekit-reload  
  href="{isTree ? `/${kind}/${newBasename}` : `/${kind}/${newBasename}/tree`}" 
  class="px-2 py-1 rounded-lg cursor-pointer text-xs text-nowrap text-white 
  bg-blue-600 dark:bg-blue-700 hover:bg-blue-700 dark:hover:bg-blue-600">
  {isTree ? 'Path' : 'Tree'} Browser
</a>
