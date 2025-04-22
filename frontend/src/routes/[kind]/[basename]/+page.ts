import { error } from '@sveltejs/kit'

export async function load({ url, params, fetch }) {
  const kind = params.kind
  const basename = params.basename
  const urlPath = url.searchParams.get("path")?.trim() ?? "" 
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

  return { kind, basename, urlPath, nspIp: nspInfo.ip }
}