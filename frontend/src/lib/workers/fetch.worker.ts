import { error, type HttpError } from "@sveltejs/kit"

import type { PathDef } from "$lib/structure"
import type { FetchPostMessage } from "./structure"

onmessage = async (event: MessageEvent<FetchPostMessage>) => {
  const {kind, basename} = event.data
  
  try {
    let paths: PathDef[] = []
    const response = await fetch(`/api/generate/${kind}/${basename}`)

    if(!response.ok) {
      const errorText = await response.text();
      throw error(404, errorText);
    }

    const jsonPaths = await response.json()
    paths = jsonPaths.map((k: PathDef) => ({...k, "is-state": ("is-state" in k ? "R" : "RW")}))

    postMessage({ success: true, message: "", paths })
    
  } catch(error) {
    postMessage({ success: false, message: (error as HttpError).body.message })
  }
}

export {}