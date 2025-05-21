export interface PayLoad {
  kind: string
  basename: string
  urlPath: string
}

export interface TreePayLoad {
  kind: string
  basename: string
  urlPath: string
  crossLaunched: boolean
}

export interface ComparePayLoad {
  x: string
  y: string
  xKind: string
  xBasename: string
  yKind: string
  yBasename: string
  urlPath: string
}

export interface PathDef {
  kind?: string
  basename?: string
  compare?: string
  path: string
  "path-with-prefix": string
  type: string
  description: string
  default?: string
  namespace?: string
  "is-state"?: string
  "is-rpc"?: boolean
  "is-notification"?: boolean
  "is-action"?: boolean
  "enum-values"?: string[]
  "if-features"?: string[]
  "added-filter" : string
}

export interface OfflineInfo {
  id: string
  nspIp: string
  timestamp: string
  module: string
  name: string
}