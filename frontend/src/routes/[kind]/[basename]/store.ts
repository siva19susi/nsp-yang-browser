import { writable, derived } from "svelte/store"

import { count } from "$lib/components/sharedStore"
import type { PathDef } from "$lib/structure"
import { searchBasedFilter } from "$lib/components/functions"

// WRITABLE STORES
export const searchStore = writable("")
export const stateStore = writable("")
export const prefixStore = writable(false)
export const defaultStore = writable(false)

export const yangPaths = writable<PathDef[]>([])
export const start = writable(0)

// DERIVED STORES
export const stateFilter = derived([stateStore, yangPaths], ([$stateStore, $yangPaths]) => 
  $yangPaths.filter((x: PathDef) => $stateStore == "" ? true : x["is-state"] == $stateStore))

export const searchFilter = derived([searchStore, stateFilter, prefixStore], ([$searchStore, $stateFilter, $prefixStore]) => 
  $stateFilter.filter((x: PathDef) => searchBasedFilter(x, $searchStore, $prefixStore)))

export const withDefaultFilter = derived([searchFilter, defaultStore], ([$searchFilter, $defaultStore]) => 
  $searchFilter.filter((x: PathDef) => $defaultStore ? "default" in x : x))

export const total = derived(withDefaultFilter, ($withDefaultFilter) => { 
  start.set(0)
  return $withDefaultFilter.length
})

export const end = derived([start, total], ([$start, $total]) => 
  ($start + count) <= $total ? ($start + count) : $total)

export const paginated = derived([start, end, withDefaultFilter], ([$start, $end, $withDefaultFilter]) => 
  $withDefaultFilter.slice($start, $end))
