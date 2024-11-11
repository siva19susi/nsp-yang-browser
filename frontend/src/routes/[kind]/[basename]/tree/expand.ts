import { removeKeyDefault } from "$lib/components/functions"

import type { YangTreePaths } from "$lib/workers/structure"

export function nameMatchesTerm(urlPath: string, targetName: string) {
  const searchKeys = urlPath.split(/(\s+)/).map(x => x.toLowerCase())
  return searchKeys.some(v => targetName.indexOf(v) !== -1)
}

export function decideExpand(folder: YangTreePaths, crossLaunched: boolean = false, urlPath: string = ""): boolean {
  let result = false
  urlPath = removeKeyDefault(urlPath)

  if(crossLaunched && urlPath !== "") {
    return urlPath.indexOf(folder.details.path) !== -1
  } 
  if(urlPath !== "") {
    // temporary until we have a good solution
    if(urlPath.split(/(\s+)/).length === 1 && urlPath.indexOf("/") != -1) {
      return urlPath.indexOf(folder.details.path) !== -1
    }
    urlPath = urlPath.replaceAll("/", " ")

    result = result || nameMatchesTerm(urlPath, folder.name)
    
    if(folder.children) {
      result = result || folder.children.some((item: YangTreePaths) => {
        if (item.type === "folder") {
          return decideExpand(item, crossLaunched, urlPath);
        } else {
          return nameMatchesTerm(urlPath, item.name);
        }
      })
    }
  }
  return result
}