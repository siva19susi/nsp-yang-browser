import { error } from "@sveltejs/kit"

export async function load({ fetch }) {
  const backendActive = await fetch("/api")
  if (!backendActive.ok) {
    throw error(404, "Backend is not active")
  }
}