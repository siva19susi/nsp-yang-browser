import type { PathDef } from '$lib/structure'

export function kindView(kind: string) {
  if(kind === "local") return "uploaded"
  else if(kind === "nsp") return "from NSP"
}

export function toLower(str: string) {
  return str.trim().toLowerCase()
}

export function toUpper(str: string) {
  return str.trim().toUpperCase()
}

export function toggleSidebar() {
  document.getElementById('sidebar')?.classList.toggle('-translate-x-0')
  document.getElementById('sidebar')?.classList.toggle('-translate-x-full')
  document.getElementById('open-sidebar')?.classList.toggle('hidden')
  document.getElementById('close-sidebar')?.classList.toggle('hidden')
}

export function closeSidebar() {
  if (document.getElementById('open-sidebar')?.classList.contains("hidden")) {
    toggleSidebar();
  }
}

export function escapeText(text: string) {
  //return text.replace(/[-[\]{}()*+?.,\\^$|#\s]/g, '\\$&')
  return text.replace(/[\[\]\*]/g, '\\$&')
}

export function removeKeyDefault(text: string) {
  return text.replaceAll("=*", "")
}

export function searchBasedFilter(x: PathDef, searchTerm: string, showPrefix: boolean = false) {
  const keys = searchTerm.split(/\s+/)
  const searchStr = `${showPrefix ? x["path-with-prefix"] : x.path};${x.type}`
  return keys.every(x => searchStr.includes(x))
}

export function markFilter(target: string, term: string, from: string = "table") {
  if(term != "") {
    const keys = term.split(/\s+/)
    const pattern = (new RegExp(escapeText(keys.join('|')), 'g'))
    let markClass = "text-nokia-blue dark:text-yellow-400 bg-white dark:bg-gray-800 font-bold"
    if(from === "tree") markClass = "bg-green-300 dark:bg-green-400"
    const markTerm = (str: string) => str.replace(pattern, (match: any) => `<mark class="${markClass}">${match}</mark>`)
    return markTerm(target)
  }
  return target
}

export function markRender (node: HTMLSpanElement, text:string) {
  const action = () => node.innerHTML = text
  action()
  return {
    update(obj: string) {
      text = obj
      action()
    },
  }
}
