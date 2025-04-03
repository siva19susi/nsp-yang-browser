import { derived, writable } from "svelte/store"

import type { YangTreePaths } from "$lib/workers/structure"

// WRITABLE STORES
export const searchStore = writable("")
export const stateStore = writable("")
export const yangTarget = writable<YangTreePaths>()
export const prefixStore = writable(false)
export const defaultStore = writable(false)

// DERIVED STORES
export const yangTreeArgs = derived([searchStore, stateStore, prefixStore, defaultStore], ([$searchStore, $stateStore, $prefixStore, $defaultStore]) => 
  $searchStore + ";;" + $prefixStore + ";;" + $stateStore + ";;" + $defaultStore)
