<script lang="ts">
	import Theme from "$lib/components/Theme.svelte"

  let isSubmitting = false

  export let data
  const { kind, modules, message } = data
  let [ nspIp, nspUser ] = message.split("@")

  async function nspDisconnect(event: SubmitEvent) {
    isSubmitting = true

		const nspDisconnected = await fetch("/api/nsp/disconnect", {method: "POST"})
    if(nspDisconnected.ok) {
      window.alert(`[Success] NSP disconnected. Page will reload to take effect.`)
      window.location.reload()
    } else {
      window.alert(`[Error] Failed to disconnect NSP. Page will reload to take effect.`)
    }
  }
</script>

<!-- NAVBAR -->
<nav class="fixed top-0 z-20 p-4 w-screen select-none font-nunito bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
	<div class="flex justify-between">
		<!-- navbar left item -->
		<div class="flex items-center space-x-2">
			<a href="/" class="flex px-2"><img src="/images/navbar-logo.png" alt="Logo" width="25"/></a>
		</div>
		<!-- navbar centre item -->
    <div class="text-center">
      <p class="text-nokia-old-blue dark:text-white text-lg lg:text-xl">NSP Modules</p>
    </div>
		<!-- navbar right item -->
    <div class="flex items-center">
		  <Theme/>
    </div>
	</div>
</nav>

<div class="px-6 pb-6 pt-[75px] lg:pt-[85px]">
  <!-- NSP Connected setting -->
  <form class="p-4 space-y-6 flex flex-col items-center justify-center" method="POST" action="?/nspDisconnect" on:submit|preventDefault={nspDisconnect}>
    <div class="">
      <label for="nsp-connect-ip" class="block uppercase text-gray-700 dark:text-gray-300 text-xs mb-2">IP / FQDN</label>
      <input id="nsp-connect-ip" bind:value={nspIp} type="text" class="px-3 py-2 rounded-lg text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-400 dark:text-gray-400 bg-gray-50 dark:bg-gray-700 select-none pointer-events-none" readonly>
    </div>
    <div class="">
      <label for="nsp-connect-user" class="block uppercase text-gray-700 dark:text-gray-300 text-xs mb-2">Username</label>
      <input id="nsp-connect-user" bind:value={nspUser} type="text" class="px-3 py-2 rounded-lg text-[12.5px] border border-gray-300 dark:border-gray-600 text-gray-400 dark:text-gray-400 bg-gray-50 dark:bg-gray-700 select-none pointer-events-none" readonly>
    </div>
    <button type="submit" class="px-4 py-1.5 rounded-lg text-xs text-white bg-red-600 hover:bg-red-700 {isSubmitting ? 'animate-pulse' : ''}" disabled={isSubmitting}>{isSubmitting ? 'Disconnecting...' : 'Disconnect'}</button>
  </form>

  <!-- NSP Applications -->
  <div class="container mx-auto">
    <ul class="pt-6 flex flex-wrap items-center justify-between text-sm text-center">
      <li class="me-2">
        <a href="#" class="px-3 py-2 text-sm rounded-lg text-white bg-blue-700 hover:bg-blue-800 dark:bg-blue-600 dark:hover:bg-blue-700">Intent Types</a>
      </li>
      {#each modules as module}
        <li class="me-2">
          <a href="#" class="px-3 py-2 text-sm rounded-lg text-white bg-blue-700 hover:bg-blue-800 dark:bg-blue-600 dark:hover:bg-blue-700">{module}</a>
        </li>
      {/each}
    </ul>
  </div>
</div>