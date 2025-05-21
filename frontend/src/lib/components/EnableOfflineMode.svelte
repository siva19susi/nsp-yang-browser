<script lang="ts">
	import { error } from "@sveltejs/kit";

  export let kind: string
  export let basename: string

  async function enableOffline() {
		let response: Response

    if(kind === "telemetry") {
      response = await fetch(`/api/nsp/telemetry-type/definition?save=true`, {
        method: "POST", body: JSON.stringify({ name: basename })
      })
    } else {
      response = await fetch(`/api/${kind.replace("-", "/")}/${basename}/paths?save=true`)
    }

    if(!response.ok) {
      const errorText = await response.text();
      throw error(404, errorText);
    } else if(response.ok) {
      window.alert("[Success] Snapshot taken for offline access.")
    }
  }
</script>

<button on:click={enableOffline}
  class="px-2 py-1 rounded-lg cursor-pointer text-xs text-nowrap text-white 
  bg-gray-400 dark:bg-gray-600 hover:bg-gray-500 dark:hover:bg-gray-700">
  Enable offline mode
</button>