export async function load({ fetch }) {
  let nspInfo = {"ip": ""}
  
  const resp1 = await fetch("/api/nsp/isConnected")
  if(resp1.ok) {
    nspInfo = await resp1.json()
  }

  return {nspConnected: nspInfo.ip !== "" ? true : false}
}