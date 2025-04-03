import { error } from "@sveltejs/kit"

export async function load({ fetch }) {
  try {
    fetch("/api")
  } catch (e) {
    throw error(404, "Backend inactive: " + e);
  }
}