import { error, type HttpError } from "@sveltejs/kit"

import type { PathDef } from "$lib/structure"
import type { YangTreePayloadPostMessage, SamplePayload } from "./structure"

onmessage = async (event: MessageEvent<YangTreePayloadPostMessage>) => {
  const {kind, basename, urlPath, withPrefix, expandFull} = event.data
  const urlPathTranform = urlPath.replaceAll("]", "=*]")

  try {
    const pathResponse = await fetch(`/api/generate/${kind}/${basename}`)

    if(!pathResponse.ok) {
      const errorText = await pathResponse.text();
      throw error(404, errorText);
    }

    function urlPathFilter(x: PathDef, searchTerm: string, showPrefix: boolean) {
      const keys = searchTerm.split(/\s+/)
      const searchStr = `${showPrefix ? x["path-with-prefix"] : x.path}`
      return keys.every(x => searchStr.includes(x))
    }

    const pathJson = await pathResponse.json()
    const paths: PathDef[] = pathJson.map((k: PathDef) => ({...k, "is-state": ("is-state" in k ? "R" : "RW")}))
    const pathFilter = urlPath !== "" ? paths.filter((x: PathDef) => urlPathFilter(x, urlPathTranform, withPrefix)) : paths


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
    for(const item of pathFilter) {
      const loopPath = (withPrefix ? item["path-with-prefix"] : item.path)
      const tranformPath = urlPath !== "" && !expandFull ? loopPath.replace(urlPathTranform, "") : loopPath
      const parts = tranformPath.split("/").filter(Boolean)
      let current = tree
      
      for (const [index, part] of parts.entries()) {
        if(part.includes("[") && part.includes("=*]")) {
          const listContainer = part.split("[")[0]
          const listKeys = part.match(/\[.*?=\*\]/g).map(p => p.split("=")[0].slice(1))

          if (!(listContainer in current)) {
            current[listContainer] = [{}]
          }
          
          for (const k of listKeys) {
            current[listContainer][0][k] = "{{ mandatory_key }}"
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
        
        if(urlPath !== "" && !expandFull) break
      }
    }

    postMessage({ success: true, message: "", tree })

  } catch(error) {
    postMessage({ success: false, message: (error as HttpError).body.message })
  }
}

export {}