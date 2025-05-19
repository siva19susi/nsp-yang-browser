import { error } from '@sveltejs/kit'

export async function load({ url, params, fetch }) {
  const kind = params.kind
  const basename = params.basename
  const urlPath = url.searchParams.get("path")?.trim() ?? "" 
  const isUrlTree = url.pathname.includes("tree") ? true : false
  let nspInfo = {"ip": ""}

  if(kind !== "uploaded" && kind !== "nsp-module" && kind !== "nsp-intent-type" && kind !== "nsp-lso-operation") {
    throw error(404, "Unsupported kind")
  }

  const resp = await fetch("/api/nsp/isConnected")
  if((kind.includes("nsp") && !resp.ok)) {
    throw error(404, "Check NSP connection")
  } else if(resp.ok) {
    nspInfo = await resp.json()
  }

  return { kind, basename, urlPath, nspIp: nspInfo.ip, isUrlTree }
}