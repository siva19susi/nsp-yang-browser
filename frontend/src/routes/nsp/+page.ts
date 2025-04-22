export async function load({ fetch }) {
  let nspInfo = {"ip": "", "user": ""}
  let modules = []

  const resp1 = await fetch("/api/nsp/isConnected")
  if(resp1.ok) {
    nspInfo = await resp1.json()

    if(nspInfo.ip !== "") {
      const resp2 = await fetch("/api/nsp/modules")
      modules = await resp2.json()
    }
  }

  return { nspInfo, modules }
}