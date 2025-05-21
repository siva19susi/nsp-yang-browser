import { error } from '@sveltejs/kit'

export async function load({ params, url, fetch }) {
  const kind = params.kind
  const type = url.searchParams.get("type")?.trim() ?? ""

  if(type === "") {
    throw error(404, "Type parameter missing")
  }

  let nspInfo = {"ip": ""}
  let definition = []

  if(kind === "nsp") {
    const resp = await fetch("/api/nsp/isConnected")
    if(!resp.ok) {
      throw error(404, "Check NSP connection")
    } else if(resp.ok) {
      nspInfo = await resp.json()

      if(nspInfo.ip !== "") {
        const resp1 = await fetch("/api/nsp/telemetry-types")

        if(!resp1.ok) {
          throw error(404, "Unable to fetch supported Telemetry Types")
        } else if(resp1.ok) {
          const telemetryTypes = await resp1.json()

          if(!telemetryTypes.includes(type)) {
            throw error(404, "Unsupported Telemetry Type")
          } else {
            const resp2 = await fetch(`/api/nsp/telemetry-type/definition`, {
              method: "POST", body: JSON.stringify({ name: type })
            })
            
            if(!resp2.ok) {
              throw error(404, "Unable to fetch Telemetry Type definition")
            } else if(resp2.ok) {
              definition = await resp2.json()
            }
          }
        }
      }
    }
  } else if(kind === "offline") {
    const resp1 = await fetch(`/api/offline/${type}/paths`)
    
    if(!resp1.ok) {
      throw error(404, "Unable to fetch Telemetry Type definition")
    } else if(resp1.ok) {
      definition = await resp1.json()
    }
  }

  return { kind, type, definition, nspIp: nspInfo.ip }
}