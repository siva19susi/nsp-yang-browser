import { derived, writable } from "svelte/store"

import type { YangTreePaths } from "./structure"

// WRITABLE STORES
export const searchStore = writable("")
export const stateStore = writable<string[]>([])
export const yangTarget = writable<YangTreePaths>()
export const prefixStore = writable(false)
export const defaultStore = writable(false)

// DERIVED STORES
export const yangTreeArgs = derived([searchStore, stateStore, prefixStore, defaultStore], ([$searchStore, $stateStore, $prefixStore, $defaultStore]) => 
  $searchStore + ";;" + $prefixStore + ";;" + $stateStore.join("#") + ";;" + $defaultStore)
