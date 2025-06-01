import { error } from '@sveltejs/kit'

export async function load({ url, params, fetch }) {
  const kind = params.kind
  let basename = params.basename
  const urlPath = url.searchParams.get("path")?.trim() ?? "" 
  const isUrlTree = url.pathname.includes("tree") ? true : false
  let nspInfo = {"ip": ""}

  if(kind !== "offline" && kind !== "nsp-module" && kind !== "nsp-intent-type" && kind !== "nsp-lso-operation") {
    throw error(404, "Unsupported kind")
  }

  if(kind.includes("nsp")) {
    const resp = await fetch("/api/nsp/isConnected")
    if(!resp.ok) {
      throw error(404, "Check NSP connection")
    } else if(resp.ok) {
      nspInfo = await resp.json()
    }
  }

  if(kind === "offline") {
    const response = await fetch(`/api/offline/list/${basename}`)
    if(response.ok) {
      const tmp = await response.json()
      basename = Object.values(tmp).join("__")
    }
  }

  return { kind, basename, urlPath, nspIp: nspInfo.ip, isUrlTree }
}