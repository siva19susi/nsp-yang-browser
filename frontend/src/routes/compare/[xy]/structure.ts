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
