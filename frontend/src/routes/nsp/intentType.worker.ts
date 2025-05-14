import { error, type HttpError } from "@sveltejs/kit"

import type { IntentTypeSearchPostMessage } from "./structure"

onmessage = async (event: MessageEvent<IntentTypeSearchPostMessage>) => {
  const {page, filter} = event.data
  
  postMessage({ type: "progress", value: 50 })
  try {
    const response = await fetch(`/api/nsp/intent-types?limit=30&page=${page}&filter=${filter}`)
    if(!response.ok) {
      const errorText = await response.text();
      throw error(404, errorText);
    }

    postMessage({ type: "progress", value: 85 })
    await new Promise(res => setTimeout(res, 500))
    
    const intentTypes = await response.json()
    if(!("intentTypes" in intentTypes)) {
      intentTypes.intentTypes = []
    }

    postMessage({
      type: "complete",
      success: true,
      message: "",
      intentTypes,
    })
    
  } catch(error) {
    postMessage({ type: "complete", success: false, message: (error as HttpError).body.message })
  }
}

export {}