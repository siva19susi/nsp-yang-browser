export async function load({ fetch }) {
  let nspInfo = {"ip": "", "user": ""}
  let modules = []
  let lsoOperations = []
  let telemetryTypes = []

  const resp1 = await fetch("/api/nsp/isConnected")
  if(resp1.ok) {
    nspInfo = await resp1.json()

    if(nspInfo.ip !== "") {
      const resp2 = await fetch("/api/nsp/modules")
      modules = await resp2.json()

      const resp3 = await fetch("/api/nsp/lso-operations")
      lsoOperations = await resp3.json()

      const resp4 = await fetch("/api/nsp/telemetry-types")
      telemetryTypes = await resp4.json()
    }
  }

  return { nspInfo, modules, lsoOperations, telemetryTypes  }
}