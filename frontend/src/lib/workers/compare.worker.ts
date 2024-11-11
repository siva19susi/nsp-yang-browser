import { error, type HttpError } from "@sveltejs/kit"

import type { PathDef } from "$lib/structure"
import { kindView } from "$lib/components/functions"
import type { ComparePostMessage, DiffResponseMessage } from "$lib/workers/structure"

onmessage = async (event: MessageEvent<ComparePostMessage>) => {
  const { xKind, yKind, xBasename, yBasename } = event.data

  let xpaths: PathDef[] = []
  let ypaths: PathDef[] = []

  async function fetchPaths(kind: string, basename: string) {
    const pathResponse = await fetch(`/api/generate/${kind}/${basename}`)

    if(!pathResponse.ok) {
      const errorText = await pathResponse.text();
      throw error(404, errorText);
    }

    const pathJson = await pathResponse.json()
    return pathJson.map((k: PathDef) => ({...k, 
      kind, basename, compareTo: `${basename} (${kindView(kind)})`, 
      "is-state": ("is-state" in k ? "R" : "RW")}))
  }

  try {
    xpaths = await fetchPaths(xKind, xBasename)
    ypaths = await fetchPaths(yKind, yBasename)

    // Start of Compare operation
    const xOnlyPath = xpaths.map((k :PathDef) => k.path)
    const yOnlyPath = ypaths.map((k :PathDef) => k.path)

    const getPathObj = (list: PathDef[], path: string) => list.filter((k :PathDef) => k.path === path)

    const typeChange: DiffResponseMessage[] = []
    const removedFromX: DiffResponseMessage[] = []
    const newInY: DiffResponseMessage[] = []

    const setX = new Set(xOnlyPath)
    const setY = new Set(yOnlyPath)

    for (const item of setX) {
      if (setY.has(item)) {
        const xObj = getPathObj(xpaths, item)[0]
        const yObj = getPathObj(ypaths, item)[0]
        if(xObj.type !== yObj.type) {
          typeChange.push({...yObj, fromType: xObj.type, fromRel: `${xBasename}(${kindView(xKind)})`, compare: "~"})
        }
      } else {
        const xObj = getPathObj(xpaths, item)[0]
        removedFromX.push({...xObj, compare: "-"})
      }
    }

    for (const item of setY) {
      if (!setX.has(item)) {
        const yObj = getPathObj(ypaths, item)[0]
        newInY.push({...yObj, compare: "+"})
      }
    }

    const diff = [...newInY, ...removedFromX, ...typeChange].sort((a, b) => {
      const keyA = a["path"]
      const keyB = b["path"]
      if (keyA < keyB) return -1
      if (keyA > keyB) return 1
      return 0
    })

    postMessage({success: true, message: "", diff})
  } catch(error) {
    postMessage({ success: false, message: (error as HttpError).body.message })
  }
}

export {}