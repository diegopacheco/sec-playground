import type { Claims } from '../types'

interface TokenInspectorProps {
  claims: Claims
  token: string
  originalToken: string
  onTokenChange: (token: string) => void
}

export function TokenInspector({ claims, token, originalToken, onTokenChange }: TokenInspectorProps) {
  const modified = token !== originalToken
  const corrupt = () => {
    const final = token.at(-1) === 'A' ? 'B' : 'A'
    onTokenChange(token.slice(0, -1) + final)
  }
  return (
    <section className="inspector panel">
      <div className="panel-heading">
        <div>
          <span className="section-kicker">Credential anatomy</span>
          <h2>Signed claims</h2>
        </div>
        <span className={`integrity-pill ${modified ? 'is-broken' : ''}`}>{modified ? 'signature changed' : 'signature intact'}</span>
      </div>
      <div className="claim-grid">
        <Claim label="sub" value={claims.sub} tone="cyan" />
        <Claim label="act" value={claims.act} tone="lime" />
        <Claim label="aud" value={claims.aud} tone="amber" />
        <Claim label="depth" value={String(claims.maximum_depth)} tone="plain" />
      </div>
      <div className="scope-row">
        <span>scope</span>
        <div>{claims.scope.map(scope => <b key={scope}>{scope}</b>)}</div>
      </div>
      <div className="token-string" aria-label="Compact signed token">{token}</div>
      <div className="inspector-actions">
        <button className="text-button danger" onClick={corrupt}>Alter signature</button>
        <button className="text-button" onClick={() => onTokenChange(originalToken)} disabled={!modified}>Restore token</button>
      </div>
    </section>
  )
}

function Claim({ label, value, tone }: { label: string, value: string, tone: string }) {
  return (
    <div className={`claim claim-${tone}`}>
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  )
}
