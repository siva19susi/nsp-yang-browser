import { writable, derived } from "svelte/store"

import { count } from "$lib/components/sharedStore"
import type { TelemetryTypeDefinition } from "./structure"

function definitionSearchFilter(x: TelemetryTypeDefinition, searchTerm: string) {
  const keys = searchTerm.split(/\s+/)
  const searchStr = `${x.counterName};${x.dataType};${x.deviceXpath}`
  return keys.every(x => searchStr.includes(x))
}

// WRITABLE STORES
export const searchStore = writable("")
export const definitionStore = writable<TelemetryTypeDefinition[]>([])
export const start = writable(0)

// DERIVED STORES
export const total = derived(definitionStore, ($definitionStore) => { 
  start.set(0)
  return $definitionStore.length
})

export const searchFilter = derived([searchStore, definitionStore], ([$searchStore, $definitionStore]) => 
  $definitionStore.filter((x: TelemetryTypeDefinition) => definitionSearchFilter(x, $searchStore)))

export const end = derived([start, total], ([$start, $total]) => 
  ($start + count) <= $total ? ($start + count) : $total)

export const paginated = derived([start, end, searchFilter], ([$start, $end, $searchFilter]) => 
  $searchFilter.slice($start, $end))
