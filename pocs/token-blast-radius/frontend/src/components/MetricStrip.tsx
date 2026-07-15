import type { Metrics } from '../types'

interface MetricStripProps {
  metrics: Metrics
}

export function MetricStrip({ metrics }: MetricStripProps) {
  const values = [
    ['Reachable', metrics.reachable_resources, 'resources'],
    ['Critical', metrics.critical_resources, 'targets'],
    ['Live paths', metrics.active_delegations, 'relationships'],
    ['Depth', metrics.maximum_depth, 'hops']
  ]
  return (
    <section className="metric-strip" aria-label="Reachability metrics">
      {values.map(([label, value, unit]) => (
        <div className="metric" key={label}>
          <span>{label}</span>
          <strong>{value}</strong>
          <small>{unit}</small>
        </div>
      ))}
    </section>
  )
}
