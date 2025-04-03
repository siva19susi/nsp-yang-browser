import { error } from '@sveltejs/kit'

export async function load({ params, fetch }) {
  const kind = params.kind

  if (kind !== "local" && kind !== "nsp") {
    throw error(404, "Unsupported kind")
  }

  let modules: string[] = []
  let message = ""

  if (kind === "nsp") {
    try {
      await fetch("http://localhost:8080/api/nsp/isConnected")
      .then(response => response.json())
      .then(nspInfo => {
        message = `${nspInfo.user}@${nspInfo.ip}`
        fetch("http://localhost:8080/api/nsp/modules")
        .then(response => response.json())
        .then(response => modules = response)
      })
    } catch (e) {
      throw error(404, "Internal Error: " + e);
    }
  }

  return { kind, message, modules }
}