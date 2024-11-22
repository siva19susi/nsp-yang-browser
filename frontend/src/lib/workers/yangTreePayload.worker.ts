import { error, type HttpError } from "@sveltejs/kit"

import type { PathDef } from "$lib/structure"
import type { FetchPostMessage, SamplePayload } from "./structure"

onmessage = async (event: MessageEvent<FetchPostMessage>) => {
  const {kind, basename} = event.data

  try {
    let paths: PathDef[] = []
    const pathResponse = await fetch(`/api/generate/${kind}/${basename}`)

    if(!pathResponse.ok) {
      const errorText = await pathResponse.text();
      throw error(404, errorText);
    }

    const pathJson = await pathResponse.json()
    paths = pathJson.map((k: PathDef) => ({...k, "is-state": ("is-state" in k ? "R" : "RW")}))

    function getSampleValue(item: PathDef) {
      if (item.default !== undefined) return item.default;
      switch (item.type) {
        case "enumeration":
          return item["enum-values"] ? item["enum-values"][0] : null
        case "uint16":
          return 0
        case "uint32":
          return 0
        case "int16":
          return 0
        case "int32":
          return 0
        case "int64":
          return 0
        case "string":
          return ""
        case "decimal64": 
          return 0
        case "boolean":
          return false
        default:
          return null
      }
    }

    const tree: SamplePayload = {}
    for(const item of paths) {
      const parts = item.path.split("/").filter(Boolean)
      let current = tree

      parts.forEach((part, index) => {
        if(part.includes("[") && part.includes("=*]")) {
          const listContainer = part.split("[")[0]
          const listKeys = part.match(/\[.*?=\*\]/g).map(p => p.split("=")[0].slice(1))

          if (!(listContainer in current)) {
            current[listContainer] = [{}]
          }
          
          for (const k of listKeys) {
            current[listContainer][0][k] = ""
          }

          current = current[listContainer][0]
        } else {
          if(!(part in current)) {
            current[part] = (index === parts.length - 1)
              ? getSampleValue(item)
              : {}
          }
          current = current[part];
        }
      })
    }

    postMessage({ success: true, message: "", tree })

  } catch(error) {
    postMessage({ success: false, message: (error as HttpError).body.message })
  }
}

export {}