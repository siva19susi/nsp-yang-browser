import { derived, writable } from "svelte/store"

import type { YangTreePaths } from "$lib/workers/structure"

// WRITABLE STORES
export const searchStore = writable("")
export const stateStore = writable("")
export const yangTarget = writable<YangTreePaths>()
export const defaultStore = writable(false)

// DERIVED STORES
export const yangTreeArgs = derived([searchStore, stateStore, defaultStore], ([$searchStore, $stateStore, $defaultStore]) => 
  $searchStore + ";;" + $defaultStore + ";;" + $stateStore)
