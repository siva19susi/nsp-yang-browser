import type { PathDef } from "$lib/structure"

export interface DiffResponseMessage extends PathDef {
  fromType?: string
  fromRel?: string
  compare: string
}

export interface ComparePostMessage {
  x: string
  y: string
  xKind: string
  xBasename: string
  yKind: string
  yBasename: string
}

export interface CompareResponseMessage {
  success: boolean
  message: string
  diff: DiffResponseMessage[]
}

export interface RepoListResponse {
  name: string
  files: string[]
}

export interface RepoResponseMessage {
  kind: string
  success: boolean
  message: string
  repo: RepoListResponse[]
}

export interface FetchPostMessage {
  kind: string
  basename: string
}

export interface FetchResponseMessage {
  success: boolean
  message: string
  paths: PathDef[]
}

export interface YangTreeContainer {
  path: string
}

export interface YangTreePostMessage {
  kind: string
  basename: string
  searchInput: string
  stateInput: string
  defaultInput: boolean
}

export interface YangTreeResponseMessage {
  success: boolean
  message: string
  node: YangTreePaths
}

export interface YangTreePaths {
  name: string
  type: string
  children: YangTreePaths[]
  details: YangTreeContainer | PathDef
}