import { error } from '@sveltejs/kit'

export async function load({ params, url, fetch }) {
  const kind = params.kind
  const basename = params.basename
  const urlPath = url.searchParams.get("path")?.trim() ?? ""
  let nspInfo = {"ip": ""}
  let intents: string[] = []

  if(kind !== "uploaded" && kind !== "nsp-intent-type") {
    throw error(404, "Unsupported kind")
  }

  const resp = await fetch("/api/nsp/isConnected")
  if(resp.ok) {
    nspInfo = await resp.json()

    if(nspInfo.ip !== "") {
      const resp1 = await fetch(`/api/nsp/intent-type/${basename}/intents`)
      intents = await resp1.json()
    }
  } else {
    throw error(404, "Check NSP connection")
  }

  return { kind, basename, urlPath, nspIp: nspInfo.ip, intents }
}