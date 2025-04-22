export interface RepoListResponse {
  name: string
  files: string[]
}

export interface RepoResponseMessage {
  success: boolean
  message: string
  repo: RepoListResponse[]
}