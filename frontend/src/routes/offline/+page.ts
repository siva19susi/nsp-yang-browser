import { error } from "@sveltejs/kit"

export async function load({ fetch }) {
  const resp = await fetch("/api/offline/list")
  if(!resp.ok) {
    throw error(404, "Error fetching uploaded yang repo")
  }

  return { localRepo: await resp.json() }
}