import { writable, derived } from "svelte/store"

export const localSearchStore = writable<string[]>([])
export const localRepoStore = writable<string[]>([])

export const localSearchFilter = derived([localRepoStore, localSearchStore], 
  ([$localRepoStore, $localSearchStore]) => $localRepoStore.filter((x: string) => 
    $localSearchStore.some((y: string) => x.includes(y))))
