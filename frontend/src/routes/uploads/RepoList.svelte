<script lang="ts">
	import { compare } from "$lib/components/sharedStore"
	import type { RepoListResponse } from "./structure"
  
  export let repoDetail: RepoListResponse = { name: "", files: [] }
  let files: FileList

  function closeRepoList() {
    repoDetail = { name: "", files: [] }
  }

  async function removeRepo(repo: string) {
    if (window.confirm(`Are you sure you want to delete: ${repo}`)) {
      const removeOperation = await fetch(`/api/delete/${repo}`, {method: "DELETE"})
      if (removeOperation.ok) {
        compare.clear()
        window.alert(`[Success] ${repo} has been deleted. Page will reload to take effect.`)
        window.location.reload()
      } else {
        window.alert(`[Error] Unable to delete (${repo}) at the moment. Try again later.`)
      }
    }
  }

  async function removeYangEntry(repo: string, yangEntry: string) {
    if (window.confirm(`Are you sure you want to delete: ${yangEntry}`)) {
      let url = "/api/delete"
      if(repoDetail.name.includes(".yang")) {
        url += `/file/${yangEntry}`
      } else {
        url += `/${repo}/file/${yangEntry}`
      }
      const removeOperation = await fetch(url, {method: "DELETE"})
      if (removeOperation.ok) {
        window.alert(`[Success] ${yangEntry} has been deleted from ${repo}. Page will reload to take effect.`)
        window.location.reload()
      } else {
        window.alert(`[Error] Unable to delete (${repo}/${yangEntry}) at the moment. Try again later.`)
      }
    }
  }

  async function handleUploadFile() {
    if(files) {
      const filename = files[0].name
      if(!filename.endsWith(".yang")) {
        window.alert(`[Error] Only .yang files are supported`)
      }
      const formData = new FormData()
      formData.append("file", files[0])

      let url = "/api/upload/file"
      if(!repoDetail.name.includes(".yang")) {
        url += `/${repoDetail.name}`
      }
      const uploadOperation = await fetch(url, {
        method: "POST", body: formData
      })

      if(!uploadOperation.ok) {
        const errorText = await uploadOperation.text()
        window.alert(`[Error] Failed to upload ${filename}: ${errorText}`)
      }

      window.alert(`[Success] ${filename} was uploaded to ${repoDetail.name}. Page will reload to take effect.`)
      window.location.reload()
    }
  }
</script>

<svelte:window on:keyup={({key}) => key === "Escape" ? closeRepoList() : ""}/>

<div id="repoListPopup" class="fixed px-6 py-4 inset-0 z-50 items-center { repoDetail.name !== ""  ? '' : 'hidden'}">
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="fixed inset-0 bg-gray-800 bg-opacity-75 transition-opacity" on:click|stopPropagation={closeRepoList}></div>
  <div id="popupContent" class="flex min-h-full justify-center items-center">
    <div class="relative transform overflow-hidden rounded-lg bg-white dark:bg-gray-700 text-left shadow-xl transition-all sm:my-8 max-w-4xl">
      <div id="popupHeader" class="flex items-center justify-between space-x-2 px-4 py-2 rounded-t bg-gray-200 dark:bg-gray-600 border-b border-gray-200 dark:border-gray-600">
        <div class="flex items-center space-x-2">
          <span class="text-lg text-gray-900 dark:text-gray-300">{repoDetail.name}</span>
          {#if !repoDetail.name.includes(".yang")}
            {#if !repoDetail.name.includes("from-nsp")}
              <a target="_blank" href="/api/download/{repoDetail.name}" class="text-gray-500 dark:text-gray-400 hover:bg-gray-300 dark:hover:bg-gray-800 rounded-lg px-2 py-1">
                <svg class="w-4 h-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                  <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 13V4M7 14H5a1 1 0 0 0-1 1v4a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1v-4a1 1 0 0 0-1-1h-2m-1-5-4 5-4-5m9 8h.01"/>
                </svg>
              </a>
            {/if}
            <button class="hover:bg-red-500 hover:text-white text-red-400 rounded-lg px-2 py-1 text-xs" on:click={() => removeRepo(repoDetail.name)}>
              <svg class="w-4 h-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 7h14m-9 3v8m4-8v8M10 3h4a1 1 0 0 1 1 1v3H9V4a1 1 0 0 1 1-1ZM6 7h12v13a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1V7Z"/>
              </svg>
            </button>
          {/if}
        </div>
        <button type="button" class="text-gray-500 hover:bg-gray-300 hover:text-gray-900 rounded-lg text-sm inline-flex justify-center items-center dark:hover:bg-gray-700 dark:hover:text-white" on:click={closeRepoList}>
          <svg class="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"/>
          </svg>
          <span class="sr-only">Close modal</span>
        </button>
      </div>
      <div id="repoListPopupBody">
        <div class="p-4 border-b dark:border-gray-600">
          <input id="yangDropzone" type="file" accept=".yang" class="peer hidden" bind:files on:change={handleUploadFile} />
          <label for="yangDropzone" class="flex flex-col space-y-2 px-4 py-3 h-full w-full items-center justify-center cursor-pointer rounded-lg text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-600 hover:bg-gray-200 dark:hover:bg-gray-800">
            <p class="text-sm">Click to upload .yang file</p>
            {#if files}
              <p class="text-black dark:text-white text-xs"><span>Selected:</span> {`${files[0].name} (${files[0].size} bytes)`}</p>
              <div class="rounded-md h-4 w-4 border-2 border-blue-300 animate-spin"></div>
            {/if}
          </label>
          {#if repoDetail.name.includes(".yang")}
            <div class="px-4 pt-4 text-xs text-gray-500 dark:text-gray-400 text-center">
              <p>Below files will be automatically resolved as dependencies</p>
              <p>during Uploaded or NSP YANG browsing.</p>
            </div>
          {/if}
          <p class="pt-3 text-xs text-gray-700 dark:text-gray-300 text-right">Total items: {repoDetail.files.length}</p>
        </div>
        <div class="md:w-[500px] max-h-[350px] overflow-auto scroll-light dark:scroll-dark">
          <ul class="mb-1">
            {#each repoDetail.files as filename, i}
              <li class="text-sm text-gray-700 dark:text-gray-300 darl:text-gray-200 {i > 0 ? 'border-t dark:border-gray-600' : ''}">
                <div class="flex items-center justify-between">
                  <p class="px-4 py-3">{filename}</p>
                  <div class="flex items-center mx-4 space-x-4">
                    <a target="_blank" href="/api/download/{repoDetail.name}/file/{filename}" class="text-gray-500 dark:text-gray-400 hover:bg-gray-300 dark:hover:bg-gray-800 rounded-lg px-2 py-1">
                      <svg class="w-4 h-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                        <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 13V4M7 14H5a1 1 0 0 0-1 1v4a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1v-4a1 1 0 0 0-1-1h-2m-1-5-4 5-4-5m9 8h.01"/>
                      </svg>
                    </a>
                    {#if repoDetail.name.includes(".yang") || repoDetail["files"].length > 1}
                      <button class="hover:bg-red-500 hover:text-white text-red-400 rounded-lg p-1"  on:click={() => removeYangEntry(repoDetail.name, filename)}>
                        <svg class="w-4 h-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                          <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 7h14m-9 3v8m4-8v8M10 3h4a1 1 0 0 1 1 1v3H9V4a1 1 0 0 1 1-1ZM6 7h12v13a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1V7Z"/>
                        </svg>
                      </button>
                    {/if}
                  </div>
                </div>
              </li>
            {/each}
          </ul>
        </div>
      </div>
    </div>
  </div>
</div>