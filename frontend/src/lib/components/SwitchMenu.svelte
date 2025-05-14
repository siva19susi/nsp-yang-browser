<script lang="ts">
  export let isUpload: boolean
</script>

<a data-sveltekit-reload  
  href="{isUpload ? `/nsp` : `/uploads`}" 
  class="px-2 py-1 rounded-lg cursor-pointer text-xs text-nowrap text-white 
  bg-blue-600 dark:bg-blue-700 hover:bg-blue-700 dark:hover:bg-blue-600">
  {isUpload ? 'NSP Connect' : 'Uploads'}
</a>
