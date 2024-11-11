import { writable, derived } from "svelte/store"

import type { RepoListResponse } from "$lib/workers/structure"

export const localSearchStore = writable<string[]>([])
export const nspSearchStore = writable<string[]>([])
export const localRepoStore = writable<RepoListResponse[]>([])
export const nspRepoStore = writable<RepoListResponse[]>([])

export const localSearchFilter = derived([localRepoStore, localSearchStore], 
  ([$localRepoStore, $localSearchStore]) => $localRepoStore.filter((x: RepoListResponse) => 
    $localSearchStore.some((y: string) => x.name.includes(y))))

  export const nspSearchFilter = derived([nspRepoStore, nspSearchStore], 
  ([$nspRepoStore, $nspSearchStore]) => $nspRepoStore.filter((x: RepoListResponse) => 
    $nspSearchStore.some((y: string) => x.name.includes(y))))
