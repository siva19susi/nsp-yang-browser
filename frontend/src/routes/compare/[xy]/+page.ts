import { error } from '@sveltejs/kit'

export async function load({ url, params, fetch }) {
  const xy = params.xy

  const sep = xy.split("..")
  if(sep.length !== 2) {
    throw error(404, "Unsupported X..Y compare parameter")
  }

  const x = sep[0]
  const y = sep[1]

  if(x === y) {
    throw error(404, "X & Y models cannot be the same")
  }

  let [xKind, xBasename] = x.split("@")
  let [yKind, yBasename] = y.split("@")

  if(xKind === "offline") {
    const response = await fetch(`/api/offline/list/${xBasename}`)
    if(response.ok) {
      const tmp = await response.json()
      xBasename = Object.values(tmp).join("__")
    }
  }

  if(yKind === "offline") {
    const response = await fetch(`/api/offline/list/${yBasename}`)
    if(response.ok) {
      const tmp = await response.json()
      yBasename = Object.values(tmp).join("__")
    }
  }

  const urlPath = url.searchParams.get("path")?.trim() ?? ""
        
  return { x, y, xKind, xBasename, yKind, yBasename, urlPath }
}