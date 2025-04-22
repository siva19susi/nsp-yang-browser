import type { PathDef } from "$lib/structure"

export interface FetchPostMessage {
  kind: string
  basename: string
}

export interface FetchResponseMessage {
  success: boolean
  message: string
  paths: PathDef[]
}