interface SamplePayload {
  [key: string]: unknown
}

export interface IntentQueryPostMessage {
  url: string
  target: string
  "intent-key": string
}

export interface IntentQueryResponseMessage {
  type: string
  value: number
  success: boolean
  message: string
  output: SamplePayload
}