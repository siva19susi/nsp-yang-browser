import { error } from '@sveltejs/kit'

export async function load({ url, params }) {
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

  const [xKind, xBasename] = x.split("@")
  const [yKind, yBasename] = y.split("@")

  const urlPath = url.searchParams.get("path")?.trim() ?? ""
        
  return { x, y, xKind, xBasename, yKind, yBasename, urlPath }
}