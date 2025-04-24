import { error } from '@sveltejs/kit'

export async function load({ params, url, fetch }) {
  const kind = params.kind
  const basename = params.basename
  let urlPath = url.searchParams.get("path")?.trim() ?? ""
  let nspInfo = {"ip": ""}

  if(kind !== "uploaded" && kind !== "nsp-module" && kind !== "nsp-intent-type") {
    throw error(404, "Unsupported kind")
  }

  const resp = await fetch("/api/nsp/isConnected")
  if(!resp.ok) {
    throw error(404, "Check NSP connection")
  }
  
  nspInfo = await resp.json()
  const cleanPath = (path: string) => {
    return decodeURIComponent(path)
      .replace(/\[[^\]]*=[^\]]*\]/g, '') // remove [anything=*]
      .replace(/\/[^/]+$/, '');         // remove last /segment
  }
  urlPath = cleanPath(urlPath)

  return { kind, basename, urlPath, nspIp: nspInfo.ip }
}