import { useEffect, useMemo, useState } from "react"

type Claims = {
  iss: string
  sub: string
  aud: string
  action: string
  resource: string
  max_amount_cents: number
  nbf: number
  exp: number
  jti: string
}

type IssueInput = {
  subject: string
  audience: string
  action: string
  resource: string
  max_amount_cents: number
  valid_for_seconds: number
}

type IssueResponse = {
  token: string
  claims: Claims
}

type VerifyInput = {
  token: string
  audience: string
  action: string
  resource: string
  amount_cents: number
}

type Decision = {
  allowed: boolean
  code: string
  reason: string
  claims?: Claims
}

type AuditRecord = {
  time: string
  token_id?: string
  subject?: string
  action: string
  resource: string
  allowed: boolean
  code: string
  reason: string
}

type Trace = Decision & {
  label: string
}

const initialIssue: IssueInput = {
  subject: "agent:buyer",
  audience: "payments-api",
  action: "transfer",
  resource: "account:operations",
  max_amount_cents: 50000,
  valid_for_seconds: 60,
}

const money = new Intl.NumberFormat("en-US", {
  style: "currency",
  currency: "USD",
})

const time = new Intl.DateTimeFormat("en-US", {
  hour: "2-digit",
  minute: "2-digit",
  second: "2-digit",
})

async function readJSON<T>(response: Response): Promise<T> {
  const body = (await response.json()) as T & { error?: string }
  if (!response.ok && body.error) {
    throw new Error(body.error)
  }
  return body
}

async function issueToken(input: IssueInput): Promise<IssueResponse> {
  const response = await fetch("/api/tokens", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  })
  return readJSON<IssueResponse>(response)
}

async function verifyIntent(input: VerifyInput): Promise<Decision> {
  const response = await fetch("/api/verify", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  })
  return readJSON<Decision>(response)
}

function fingerprint(value: string): string {
  if (!value) return "No token issued"
  return `${value.slice(0, 14)} ··· ${value.slice(-12)}`
}

function shortID(value?: string): string {
  if (!value) return "unresolved"
  return `${value.slice(0, 8)}…${value.slice(-6)}`
}

function App() {
  const [issueInput, setIssueInput] = useState<IssueInput>(initialIssue)
  const [issued, setIssued] = useState<IssueResponse | null>(null)
  const [verifyInput, setVerifyInput] = useState<Omit<VerifyInput, "token">>({
    audience: initialIssue.audience,
    action: initialIssue.action,
    resource: initialIssue.resource,
    amount_cents: 49999,
  })
  const [decision, setDecision] = useState<Decision | null>(null)
  const [audit, setAudit] = useState<AuditRecord[]>([])
  const [trace, setTrace] = useState<Trace[]>([])
  const [online, setOnline] = useState(false)
  const [busy, setBusy] = useState("")
  const [error, setError] = useState("")
  const [now, setNow] = useState(Date.now())

  const remaining = useMemo(() => {
    if (!issued) return 0
    return Math.max(0, issued.claims.exp - Math.floor(now / 1000))
  }, [issued, now])

  const refreshAudit = async () => {
    const response = await fetch("/api/audit")
    setAudit(await readJSON<AuditRecord[]>(response))
  }

  useEffect(() => {
    fetch("/health")
      .then((response) => response.ok)
      .then(setOnline)
      .catch(() => setOnline(false))
    refreshAudit().catch(() => undefined)
    const timer = window.setInterval(() => setNow(Date.now()), 1000)
    return () => window.clearInterval(timer)
  }, [])

  const updateIssue = <K extends keyof IssueInput>(key: K, value: IssueInput[K]) => {
    setIssueInput((current) => ({ ...current, [key]: value }))
  }

  const handleIssue = async () => {
    setBusy("issue")
    setError("")
    setDecision(null)
    setTrace([])
    try {
      const result = await issueToken(issueInput)
      setIssued(result)
      setVerifyInput({
        audience: result.claims.aud,
        action: result.claims.action,
        resource: result.claims.resource,
        amount_cents: Math.max(0, result.claims.max_amount_cents - 1),
      })
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "Token issuance failed")
    } finally {
      setBusy("")
    }
  }

  const handleVerify = async () => {
    if (!issued) return
    setBusy("verify")
    setError("")
    try {
      const result = await verifyIntent({ token: issued.token, ...verifyInput })
      setDecision(result)
      setTrace((current) => [...current, { ...result, label: "Manual request" }])
      await refreshAudit()
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "Intent verification failed")
    } finally {
      setBusy("")
    }
  }

  const runSequence = async () => {
    setBusy("sequence")
    setError("")
    setDecision(null)
    setTrace([])
    try {
      const fresh = await issueToken(issueInput)
      setIssued(fresh)
      const base: VerifyInput = {
        token: fresh.token,
        audience: fresh.claims.aud,
        action: fresh.claims.action,
        resource: fresh.claims.resource,
        amount_cents: fresh.claims.max_amount_cents,
      }
      const checks: Array<{ label: string; input: VerifyInput }> = [
        { label: "Scope escalation", input: { ...base, action: "delete" } },
        {
          label: "Amount overflow",
          input: { ...base, amount_cents: fresh.claims.max_amount_cents + 1 },
        },
        { label: "Exact intent", input: base },
        { label: "Replay attempt", input: base },
      ]
      const outcomes: Trace[] = []
      for (const check of checks) {
        const result = await verifyIntent(check.input)
        outcomes.push({ ...result, label: check.label })
        setTrace([...outcomes])
      }
      setDecision(outcomes.at(-1) ?? null)
      setVerifyInput({
        audience: base.audience,
        action: base.action,
        resource: base.resource,
        amount_cents: base.amount_cents,
      })
      await refreshAudit()
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "Sequence failed")
    } finally {
      setBusy("")
    }
  }

  return (
    <main>
      <div className="signal-line" />
      <header className="masthead">
        <a className="brand" href="#top" aria-label="Intent Token Ledger home">
          <span className="brand-mark">IT</span>
          <span>
            <strong>Intent Token Ledger</strong>
            <small>Authorization evidence console</small>
          </span>
        </a>
        <div className={`service-state ${online ? "online" : "offline"}`}>
          <span />
          <div>
            <small>Issuer status</small>
            <strong>{online ? "READY · :8081" : "UNREACHABLE"}</strong>
          </div>
        </div>
      </header>

      <section className="hero" id="top">
        <div className="hero-copy">
          <p className="eyebrow">Least privilege / one action / one use</p>
          <h1>Bind the agent<br />before it acts.</h1>
          <p className="lede">
            Mint a signed authorization envelope, pressure-test every constraint,
            then consume it once. Every decision lands in the evidence ledger.
          </p>
          <div className="hero-actions">
            <button className="primary" onClick={runSequence} disabled={!online || Boolean(busy)}>
              {busy === "sequence" ? "Running sequence…" : "Run full sequence"}
            </button>
            <a href="#workbench">Configure intent</a>
          </div>
        </div>
        <div className="trust-card">
          <div className="trust-card-top">
            <span>Trust boundary 01</span>
            <span>HMAC–SHA256</span>
          </div>
          <div className="seal">
            <span>1×</span>
            <small>USE</small>
          </div>
          <dl>
            <div><dt>Identity</dt><dd>Subject bound</dd></div>
            <div><dt>Capability</dt><dd>Exact action</dd></div>
            <div><dt>Blast radius</dt><dd>Resource + cap</dd></div>
            <div><dt>Window</dt><dd>Short lived</dd></div>
          </dl>
        </div>
      </section>

      <section className="principles" aria-label="Security controls">
        <article><span>01</span><strong>Declare</strong><p>State the actor, target, action, resource, and ceiling.</p></article>
        <article><span>02</span><strong>Constrain</strong><p>Sign the smallest useful authority for the shortest useful time.</p></article>
        <article><span>03</span><strong>Consume</strong><p>Accept one matching request and reject every replay.</p></article>
        <article><span>04</span><strong>Account</strong><p>Retain an allow or deny record for every decision.</p></article>
      </section>

      <section className="workbench" id="workbench">
        <div className="section-heading">
          <div><p className="eyebrow">Control surface</p><h2>Intent workbench</h2></div>
          <p>Change a signed claim or requested value. The verifier reports exactly which boundary held.</p>
        </div>

        {error && <div className="error-banner" role="alert"><strong>Request failed</strong><span>{error}</span></div>}

        <div className="workbench-grid">
          <form className="panel issue-panel" onSubmit={(event) => { event.preventDefault(); handleIssue() }}>
            <div className="panel-number">A</div>
            <div className="panel-title"><span>Signed envelope</span><h3>Issue authority</h3></div>
            <label>Agent subject<input value={issueInput.subject} onChange={(event) => updateIssue("subject", event.target.value)} required /></label>
            <label>Target audience<input value={issueInput.audience} onChange={(event) => updateIssue("audience", event.target.value)} required /></label>
            <div className="field-row">
              <label>Allowed action<input value={issueInput.action} onChange={(event) => updateIssue("action", event.target.value)} required /></label>
              <label>Bound resource<input value={issueInput.resource} onChange={(event) => updateIssue("resource", event.target.value)} required /></label>
            </div>
            <div className="field-row">
              <label>Amount ceiling · cents<input type="number" min="0" value={issueInput.max_amount_cents} onChange={(event) => updateIssue("max_amount_cents", Number(event.target.value))} required /></label>
              <label>Lifetime · seconds<input type="number" min="1" max="3600" value={issueInput.valid_for_seconds} onChange={(event) => updateIssue("valid_for_seconds", Number(event.target.value))} required /></label>
            </div>
            <button className="primary full" disabled={!online || Boolean(busy)}>{busy === "issue" ? "Signing…" : "Issue intent token"}</button>
          </form>

          <div className="token-column">
            <div className={`token-card ${issued ? "active" : "empty"}`}>
              <div className="token-head"><span>Bearer artifact</span><span>{issued ? `${remaining}s left` : "awaiting issue"}</span></div>
              <p className="fingerprint">{fingerprint(issued?.token ?? "")}</p>
              {issued ? (
                <>
                  <div className="claim-pair"><span>Subject</span><strong>{issued.claims.sub}</strong></div>
                  <div className="claim-pair"><span>Audience</span><strong>{issued.claims.aud}</strong></div>
                  <div className="claim-pair"><span>Capability</span><strong>{issued.claims.action} → {issued.claims.resource}</strong></div>
                  <div className="claim-pair"><span>Ceiling</span><strong>{money.format(issued.claims.max_amount_cents / 100)}</strong></div>
                  <div className="claim-pair"><span>Token ID</span><strong>{shortID(issued.claims.jti)}</strong></div>
                  <div className="ttl-track"><span style={{ width: `${Math.min(100, remaining / issueInput.valid_for_seconds * 100)}%` }} /></div>
                </>
              ) : <p className="empty-copy">Issue a token to reveal its signed constraint set.</p>}
            </div>
            <div className="caution-note"><strong>Treat as a secret.</strong><span>The full bearer value stays out of the visible interface; only a fingerprint is rendered.</span></div>
          </div>

          <form className="panel verify-panel" onSubmit={(event) => { event.preventDefault(); handleVerify() }}>
            <div className="panel-number">B</div>
            <div className="panel-title"><span>Runtime request</span><h3>Challenge intent</h3></div>
            <label>Presented audience<input value={verifyInput.audience} onChange={(event) => setVerifyInput((current) => ({ ...current, audience: event.target.value }))} required /></label>
            <div className="field-row">
              <label>Requested action<input value={verifyInput.action} onChange={(event) => setVerifyInput((current) => ({ ...current, action: event.target.value }))} required /></label>
              <label>Requested resource<input value={verifyInput.resource} onChange={(event) => setVerifyInput((current) => ({ ...current, resource: event.target.value }))} required /></label>
            </div>
            <label>Requested amount · cents<input type="number" min="0" value={verifyInput.amount_cents} onChange={(event) => setVerifyInput((current) => ({ ...current, amount_cents: Number(event.target.value) }))} required /></label>
            <div className="probe-row">
              <button type="button" onClick={() => setVerifyInput((current) => ({ ...current, action: "delete" }))}>Escalate action</button>
              <button type="button" onClick={() => setVerifyInput((current) => ({ ...current, amount_cents: issueInput.max_amount_cents + 1 }))}>Exceed ceiling</button>
              <button type="button" onClick={() => setVerifyInput({ audience: issueInput.audience, action: issueInput.action, resource: issueInput.resource, amount_cents: issueInput.max_amount_cents })}>Match claims</button>
            </div>
            <button className="dark full" disabled={!issued || Boolean(busy)}>{busy === "verify" ? "Evaluating…" : "Verify request"}</button>
          </form>
        </div>
      </section>

      <section className="decision-section">
        <div className="section-heading">
          <div><p className="eyebrow">Policy result</p><h2>Decision record</h2></div>
          <p>The signature opens the door. Exact intent matching decides whether the request may cross it.</p>
        </div>
        <div className="decision-grid">
          <div className={`decision-stamp ${decision?.allowed ? "allow" : decision ? "deny" : "pending"}`}>
            <small>Latest verdict</small>
            <strong>{decision ? (decision.allowed ? "ALLOW" : "DENY") : "PENDING"}</strong>
            <span>{decision?.code ?? "issue and verify an intent"}</span>
            <p>{decision?.reason ?? "No policy decision has been recorded in this session."}</p>
          </div>
          <div className="trace-list">
            {trace.length ? trace.map((item, index) => (
              <div className="trace-row" key={`${item.label}-${index}`}>
                <span className={item.allowed ? "pass" : "block"}>{String(index + 1).padStart(2, "0")}</span>
                <div><strong>{item.label}</strong><small>{item.code.replaceAll("_", " ")}</small></div>
                <b>{item.allowed ? "passed" : "blocked"}</b>
              </div>
            )) : <div className="trace-empty">Run the full sequence to compare escalation, overflow, exact intent, and replay outcomes.</div>}
          </div>
        </div>
      </section>

      <section className="audit-section">
        <div className="section-heading">
          <div><p className="eyebrow">Append-only view</p><h2>Evidence ledger</h2></div>
          <button className="text-button" onClick={() => refreshAudit()} disabled={!online}>Refresh records</button>
        </div>
        <div className="audit-table">
          <div className="audit-header"><span>Time</span><span>Subject / token</span><span>Requested capability</span><span>Decision</span></div>
          {audit.length ? [...audit].reverse().map((record, index) => (
            <div className="audit-row" key={`${record.time}-${index}`}>
              <span>{time.format(new Date(record.time))}</span>
              <span><strong>{record.subject ?? "unknown subject"}</strong><small>{shortID(record.token_id)}</small></span>
              <span><strong>{record.action}</strong><small>{record.resource}</small></span>
              <span className={record.allowed ? "audit-allow" : "audit-deny"}><strong>{record.allowed ? "ALLOW" : "DENY"}</strong><small>{record.code}</small></span>
            </div>
          )) : <div className="audit-empty">No decisions recorded. The first verification will appear here.</div>}
        </div>
      </section>

      <footer>
        <span>Intent Token Ledger · local trust laboratory</span>
        <span>GET /health · POST /api/tokens · POST /api/verify · GET /api/audit</span>
      </footer>
    </main>
  )
}

export default App
