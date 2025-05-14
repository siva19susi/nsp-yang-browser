import { writable, derived } from "svelte/store"

export const intentTypeStore = writable<string[]>([])
export const start = writable(0)
export const total = writable(0)
export const pageCount = writable(0)

export const end = derived([start, total, pageCount], ([$start, $total, $pageCount]) => 
  ($start + $pageCount) <= $total ? ($start + $pageCount) : $total)
