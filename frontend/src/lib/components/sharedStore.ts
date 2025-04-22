import { writable } from "svelte/store"
import { browser } from "$app/environment"

export const count = 40

export const stateValues = [
  { label: "All", value: "" },
  { label: "State", value: "R" },
  { label: "Config", value: "RW" }
]

export const compareValues = [
  { label: "All", value: "" },
  { label: "Present", value: "+" },
  { label: "Not Present", value: "-" },
  { label: "Modified", value: "~" }
]

export const pathFocus = writable({})

const LOCAL_KEY = 'nsp-yang-browser-compare'

function createCompareStore() {
  const initial = browser && localStorage.getItem(LOCAL_KEY)
    ? JSON.parse(localStorage.getItem(LOCAL_KEY)!)
    : [];

  const { subscribe, set } = writable<string[]>(initial);

  function updateWith(callback: (items: string[]) => string[]) {
    if (!browser) return;
    const currentRaw = localStorage.getItem(LOCAL_KEY);
    const current = currentRaw ? JSON.parse(currentRaw) : [];

    const updated = callback(current);
    localStorage.setItem(LOCAL_KEY, JSON.stringify(updated));
    set(updated);
  }

  if (browser) {
    window.addEventListener('storage', (event) => {
      if (event.key === LOCAL_KEY && event.newValue) {
        set(JSON.parse(event.newValue));
      }
    });
  }

  return {
    subscribe,
    add(item: string) {
      updateWith((items) => {
        if (!items.includes(item)) {
          return items.length < 2 ? [...items, item] : items;
        }
        return items;
      });
    },
    remove(item: string) {
      updateWith((items) => items.filter(i => i !== item));
    },
    nspDisconnected() {
      updateWith((items) => items.filter(i => !i.includes('nsp-intent-type')));
    },
    clear() {
      updateWith(() => []);
    },
  };
}

export const compare = createCompareStore()
