export interface Claims {
  iss: string
  sub: string
  act: string
  aud: string
  scope: string[]
  maximum_depth: number
  exp: number
  jti: string
}

export interface Scenario {
  id: string
  name: string
  summary: string
  question: string
  token: string
  claims: Claims
}

export interface Node {
  id: string
  label: string
  type: string
  sensitivity: 'internal' | 'confidential' | 'critical'
  description: string
  reachable: boolean
  depth: number
}

export interface Edge {
  from: string
  to: string
  permission: string
  required_scope: string
  rationale: string
  active: boolean
  blocked_reason?: string
}

export interface AccessPath {
  target: string
  node_ids: string[]
  permissions: string[]
}

export interface Metrics {
  reachable_resources: number
  critical_resources: number
  active_delegations: number
  maximum_depth: number
  risk_score: number
  risk_level: 'contained' | 'elevated' | 'critical'
}

export interface Analysis {
  scenario_id: string
  claims: Claims
  nodes: Node[]
  edges: Edge[]
  paths: AccessPath[]
  metrics: Metrics
  verdict: string
}
