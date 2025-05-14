export interface YangTreePayloadPostMessage {
  kind: string
  basename: string
  urlPath: string
  withPrefix: boolean
  expandFull: boolean
}

export interface YangTreePayloadResponseMessage {
  success: boolean
  message: string
  tree: SamplePayload
}

export interface SamplePayload {
  [key: string]: unknown
}