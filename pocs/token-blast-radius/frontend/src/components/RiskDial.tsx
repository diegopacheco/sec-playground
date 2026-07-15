import type { Metrics } from '../types'

interface RiskDialProps {
  metrics: Metrics
}

export function RiskDial({ metrics }: RiskDialProps) {
  const angle = Math.round(metrics.risk_score * 3.6)
  return (
    <section className={`risk-dial risk-${metrics.risk_level}`} aria-label={`Risk score ${metrics.risk_score} out of 100`}>
      <div className="dial-orbit" style={{ '--risk-angle': `${angle}deg` } as React.CSSProperties}>
        <div className="dial-core">
          <strong>{metrics.risk_score}</strong>
          <span>/ 100</span>
        </div>
      </div>
      <div>
        <span className="section-kicker">Exposure rating</span>
        <h2>{metrics.risk_level}</h2>
      </div>
    </section>
  )
}
