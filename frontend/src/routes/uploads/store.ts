import { writable, derived } from "svelte/store"

import type { RepoListResponse } from "./structure"

export const localSearchStore = writable<string[]>([])
export const localRepoStore = writable<RepoListResponse[]>([])

export const localSearchFilter = derived([localRepoStore, localSearchStore], 
  ([$localRepoStore, $localSearchStore]) => $localRepoStore.filter((x: RepoListResponse) => 
    $localSearchStore.some((y: string) => x.name.includes(y))))
