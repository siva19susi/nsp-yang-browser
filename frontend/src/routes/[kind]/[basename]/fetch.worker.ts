import { error, type HttpError } from "@sveltejs/kit"

import type { PathDef } from "$lib/structure"
import type { FetchPostMessage } from "./structure"

onmessage = async (event: MessageEvent<FetchPostMessage>) => {
  const {kind, basename} = event.data
  
  try {
    let paths: PathDef[] = []

    const response = await fetch(`/api/${kind.replace("-", "/")}/${basename}/paths`)
    if(!response.ok) {
      const errorText = await response.text();
      throw error(404, errorText);
    }

    const jsonPaths = await response.json()
    paths = jsonPaths.map((k: PathDef) => {
      let value = "RW"
    
      if ("is-state" in k) value = "R"
      else if ("is-rpc" in k) value = "RPC"
      else if ("is-action" in k) value = "A"
      else if ("is-notification" in k) value = "N"
    
      return {
        ...k,
        "added-filter": value
      }
    })

    postMessage({ success: true, message: "", paths })
    
  } catch(error) {
    postMessage({ success: false, message: (error as HttpError).body.message })
  }
}

export {}