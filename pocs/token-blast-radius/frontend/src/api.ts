import type { Analysis, Scenario } from './types'

export async function loadScenarios(): Promise<Scenario[]> {
  const response = await fetch('/api/scenarios')
  if (!response.ok) {
    throw new Error('Could not load token scenarios')
  }
  return response.json() as Promise<Scenario[]>
}

export async function analyzeToken(scenarioId: string, token: string): Promise<Analysis> {
  const response = await fetch('/api/analyze', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ scenario_id: scenarioId, token })
  })
  const body = await response.json() as Analysis | { error: string }
  if (!response.ok) {
    throw new Error('error' in body ? body.error : 'Token analysis failed')
  }
  return body as Analysis
}
