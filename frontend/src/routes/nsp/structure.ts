export interface IntentTypeSearch {
  total: number
  pageCount: number
  intentTypes: string[]
}

export interface IntentTypeSearchPostMessage {
  page: number
  filter: string
}

export interface IntentTypeSearchResponseMessage {
  type: string
  value: number
  success: boolean
  message: string
  intentTypes: IntentTypeSearch
}