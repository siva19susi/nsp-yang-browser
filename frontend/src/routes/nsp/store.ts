import { writable, derived } from "svelte/store"

const count = 30

export const intentTypeStore = writable<string[]>([])
export const start = writable(0)
export const total = writable(0)
export const pageCount = writable(0)

export const end = derived([start, total, pageCount], ([$start, $total, $pageCount]) => 
  ($start + $pageCount) <= $total ? ($start + $pageCount) : $total)

//--------------------------------------------------------------------------

export const lsoStore = writable<string[]>([])
export const lsoSearch = writable("")
export const lsoStart = writable(0)

export const lsoTotal = derived(lsoStore, ($lsoStore) => { 
  lsoStart.set(0)
  return $lsoStore.length
})

export const lsoSearchFilter = derived([lsoSearch, lsoStore], ([$lsoSearch, $lsoStore]) => 
  {
    console.log($lsoSearch)
    return $lsoStore.filter(x => $lsoSearch.split(/\s+/).every(y => x.includes(y)))
  })

export const lsoEnd = derived([lsoStart, lsoTotal], ([$lsoStart, $lsoTotal]) => 
  ($lsoStart + count) <= $lsoTotal ? ($lsoStart + count) : $lsoTotal)

export const lsoPaginated = derived([lsoStart, lsoEnd, lsoSearchFilter], ([$lsoStart, $lsoEnd, $lsoSearchFilter]) => 
  $lsoSearchFilter.slice($lsoStart, $lsoEnd))

//--------------------------------------------------------------------------

export const telemetryStore = writable<string[]>([])
export const telemetrySearch = writable("")
export const telemetryStart = writable(0)

export const telemetryTotal = derived(telemetryStore, ($telemetryStore) => { 
  lsoStart.set(0)
  return $telemetryStore.length
})

export const telemetrySearchFilter = derived([telemetrySearch, telemetryStore], ([$telemetrySearch, $telemetryStore]) => 
  $telemetryStore.filter(x => $telemetrySearch.split(/\s+/).every(y => x.includes(y))))

export const telemetryEnd = derived([telemetryStart, telemetryTotal], ([$telemetryStart, $telemetryTotal]) => 
  ($telemetryStart + count) <= $telemetryTotal ? ($telemetryStart + count) : $telemetryTotal)

export const telemetryPaginated = derived([telemetryStart, telemetryEnd, telemetrySearchFilter], ([$telemetryStart, $telemetryEnd, $telemetrySearchFilter]) => 
  $telemetrySearchFilter.slice($telemetryStart, $telemetryEnd))
