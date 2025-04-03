import { error } from '@sveltejs/kit'

export async function load({ params, url }) {
  const kind = params.kind
  const basename = params.basename

  if(kind !== "local" && kind !== "nsp") {
    throw error(404, "Unsupported kind")
  }

  const urlPath = url.searchParams.get("path")?.trim() ?? "" 
  const withPrefix = url.searchParams.get("prefix")?.trim() === "true" ? true : false
  const expandFull = url.searchParams.get("expand")?.trim() === "false" ? (urlPath !== "" ? false : true) : true

  return { kind, basename, urlPath, withPrefix, expandFull }
}