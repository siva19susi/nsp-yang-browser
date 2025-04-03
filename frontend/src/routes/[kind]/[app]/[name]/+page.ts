import { error } from '@sveltejs/kit'

export async function load({ url, params }) {
  const kind = params.kind

  if(kind != "local" && kind != "nsp") {
    throw error(404, "Unsupported kind")
  }

  const app = params.app ?? ""

  if(kind == "nsp") {
    const nspIsConnected = await fetch(`/api/nsp/isConnected`)
    if (nspIsConnected.ok) {
      if(app != "modules" && app != "intent-types") {
        throw error(404, "Unsupported NSP module")
      }
    }
    else {
      throw error(404, "NSP is not connected")
    }
  }

  const name = params.name ?? ""
  const urlPath = url.searchParams.get("path")?.trim() ?? ""

  return { kind, app, name, urlPath }
}