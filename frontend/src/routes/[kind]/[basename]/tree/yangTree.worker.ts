import { error, type HttpError } from "@sveltejs/kit"

import type { PathDef } from "$lib/structure"
import type { YangTreeContainer, YangTreePostMessage, YangTreePaths } from "./structure"
import { removeKeyDefault, searchBasedFilter } from "$lib/components/functions"

onmessage = async (event: MessageEvent<YangTreePostMessage>) => {
  const {kind, basename, searchInput, prefixInput, stateInput, defaultInput} = event.data

  try {
    let paths: PathDef[] = []
    const pathResponse = await fetch(`/api/${kind.replace("-", "/")}/${kind === "offline" ? basename.split("__")[0] : basename}/paths`)

    if(!pathResponse.ok) {
      const errorText = await pathResponse.text();
      throw error(404, errorText);
    }

    const pathJson = await pathResponse.json()
    paths = pathJson.map((k: PathDef) => {
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

    const stateFilter = paths.filter((x: PathDef) => stateInput.includes(x["added-filter"]) ? true : false)
    const searchFilter = stateFilter.filter((x: PathDef) => searchBasedFilter(x, searchInput))
    const defaultFilter = searchFilter.filter((x: PathDef) => defaultInput ? "default" in x : x)

    // Tree Builder
    class TreeNode {
      name: string
      type: string
      children: YangTreePaths[]
      details: YangTreeContainer | PathDef
      constructor(name: string, isKey: boolean, details: YangTreeContainer | PathDef, type: string) {
        this.name = isKey ? name + "*" : name
        this.type = type
        this.children = []
        this.details = details
      }
    }

    const node = new TreeNode(basename, false, {path: ""}, "folder")
    const extractBetween = (str: string) => {
      const regex = /\[(.*?)\]/g
      const matches = []
      let match
      while ((match = regex.exec(str)) !== null) {
        matches.push(match[1])
      }
      return matches
    }

    let keys: string[] = []
    for (const entry of defaultFilter) {
      let currentNode = node

      const xpath = prefixInput ? entry["path-with-prefix"] : entry["path"]
      const clean = removeKeyDefault(xpath)
      const segments = clean.split("/").slice(1)
      const segLen = segments.length

      const containerPath: string[] = []

      segments.forEach((segment: string, i: number) => {
        containerPath.push(segment)
        if(segment.includes("[")) keys = extractBetween(segment)
        let childNode = currentNode.children.find((node: { name: string }) => node.name === segment)

        if (!childNode) {
          let isKey = false
          const isLast = (i == (segLen - 1))

          const paramPath = (isLast ? entry : {"path" : "/" + containerPath.join("/")})
          if(keys.length > 0 && keys.includes(segment)) isKey = true
          const nodeType = (isLast ? "file" : "folder")

          childNode = new TreeNode(segment, isKey, paramPath, nodeType)
          if(isKey) {
            currentNode.children = [childNode].concat(currentNode.children)
          }
          else currentNode.children.push(childNode)

          if(isLast) containerPath.pop()
        }

        currentNode = childNode
      })
    }

    postMessage({ success: true, message: "", node })

  } catch(error) {
    postMessage({ success: false, message: (error as HttpError).body.message })
  }
}

export {}