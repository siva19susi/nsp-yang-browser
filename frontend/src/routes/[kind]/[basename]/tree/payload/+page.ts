import { error } from '@sveltejs/kit'

export async function load({ params, url, fetch }) {
  const kind = params.kind
  const basename = params.basename
  const urlPath = url.searchParams.get("path")?.trim() ?? "" 
  const withPrefix = url.searchParams.get("prefix")?.trim() === "true" ? true : false
  const expandFull = url.searchParams.get("expand")?.trim() === "false" ? (urlPath !== "" ? false : true) : true
  let nspInfo = {"ip": ""}

  if(kind !== "uploaded" && kind !== "nsp-module" && kind !== "nsp-intent-type") {
    throw error(404, "Unsupported kind")
  }

  if(kind.includes("nsp")) {
    const resp = await fetch("/api/nsp/isConnected")
    nspInfo = await resp.json()
    if(!resp.ok) {
      throw error(404, "Check NSP connection")
    }
  }

  return { kind, basename, urlPath, withPrefix, expandFull, nspIp: nspInfo.ip }
}