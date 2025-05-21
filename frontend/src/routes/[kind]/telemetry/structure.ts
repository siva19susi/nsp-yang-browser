export interface NodeSupportKind {
  node: string
  release: string[]
}

export interface TelemetryTypeDefinition {
  telemetryType: string
  path: string
  counterName: string
  dataType: string
  deviceXpath: string
  nodeSupport: NodeSupportKind[]
}