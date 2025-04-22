import type { RepoResponseMessage } from "./structure"

onmessage = async () => {
  const response: RepoResponseMessage = {success: false, message: "", repo: []}

  const repoResponse = await fetch("/api/uploaded")
  if (repoResponse.ok) {
    response.success = true
    response.repo = await repoResponse.json()
  } else {
    response.message = `Error fetching uploaded yang repo`
  }

  postMessage(response)
}

export {}