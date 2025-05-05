import { error } from '@sveltejs/kit'

export async function load({ params, url, fetch }) {
  const kind = params.kind
  const basename = params.basename
  const urlPath = url.searchParams.get("path")?.trim() ?? ""
  const field = url.searchParams.get("field")?.trim() ?? ""
  let nspInfo = {"ip": ""}

  if(kind !== "uploaded" && kind !== "nsp-module" && kind !== "nsp-intent-type") {
    throw error(404, "Unsupported kind")
  }

  const resp = await fetch("/api/nsp/isConnected")
  if(!resp.ok) {
    throw error(404, "Check NSP connection")
  }
  
  nspInfo = await resp.json()

  return { kind, basename, urlPath, field, nspIp: nspInfo.ip }
}