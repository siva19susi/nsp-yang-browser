import { error } from '@sveltejs/kit'

export async function load({ params }) {
  const kind = params.kind
  const basename = params.basename

  if(kind !== "local" && kind !== "nsp") {
    throw error(404, "Unsupported kind")
  }

  return { kind, basename }
}