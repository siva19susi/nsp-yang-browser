import type { RepoResponseMessage } from "./structure"

onmessage = async (event: MessageEvent<string>) => {
  const kind = event.data
  const response: RepoResponseMessage = {kind, success: false, message: "", repo: []}

  async function getRepo() {
    const repoResponse = await fetch(`/api/list/${kind}`)
    if (repoResponse.ok) {
      response.success = true
      response.repo = await repoResponse.json()
    } else {
      response.message = `Error fetching ${kind} yang repo`
    }
  }

  if(kind == "local") {
    await getRepo()
  }
  else if(kind == "nsp") {
    const nspIsConnected = await fetch(`/api/nsp/isConnected`)
    if (nspIsConnected.ok) {
      const nspInfo = await nspIsConnected.json()
      response.message = `${nspInfo.user}@${nspInfo.ip}`
      await getRepo()
    }
    else {
      response.message = `NSP is not conected`
    }
  }

  postMessage(response)
}

export {}