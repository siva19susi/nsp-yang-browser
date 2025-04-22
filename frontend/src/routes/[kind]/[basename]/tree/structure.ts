import type { PathDef } from "$lib/structure"

export interface YangTreeContainer {
  path: string
}

export interface YangTreePostMessage {
  kind: string
  basename: string
  searchInput: string
  stateInput: string
  prefixInput: boolean
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