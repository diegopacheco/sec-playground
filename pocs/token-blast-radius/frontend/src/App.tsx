import { useEffect, useMemo, useState } from 'react'
import { analyzeToken, loadScenarios } from './api'
import { MetricStrip } from './components/MetricStrip'
import { NodeDetail } from './components/NodeDetail'
import { ReachabilityGraph } from './components/ReachabilityGraph'
import { RiskDial } from './components/RiskDial'
import { TokenInspector } from './components/TokenInspector'
import type { Analysis, Scenario } from './types'

export default function App() {
  const [scenarios, setScenarios] = useState<Scenario[]>([])
  const [scenarioId, setScenarioId] = useState('')
  const [token, setToken] = useState('')
  const [analysis, setAnalysis] = useState<Analysis>()
  const [selectedNode, setSelectedNode] = useState<string>()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const scenario = scenarios.find(item => item.id === scenarioId)
  const labels = useMemo(() => new Map(analysis?.nodes.map(node => [node.id, node.label]) ?? []), [analysis])
  const node = analysis?.nodes.find(item => item.id === selectedNode)
  const path = analysis?.paths.find(item => item.target === selectedNode)
  const incoming = analysis?.edges.find(edge => edge.to === selectedNode)

  useEffect(() => {
    loadScenarios()
      .then(items => {
        setScenarios(items)
        if (items[0]) {
          setScenarioId(items[0].id)
          setToken(items[0].token)
          return runAnalysis(items[0].id, items[0].token)
        }
      })
      .catch(cause => setError(message(cause)))
      .finally(() => setLoading(false))
  }, [])

  const runAnalysis = async (id: string, value: string) => {
    setLoading(true)
    setError('')
    try {
      const result = await analyzeToken(id, value)
      setAnalysis(result)
      const firstCritical = result.nodes.find(item => item.reachable && item.sensitivity === 'critical')
      const firstResource = result.nodes.find(item => item.reachable && !['human', 'agent', 'workspace'].includes(item.type))
      setSelectedNode((firstCritical ?? firstResource ?? result.nodes[0])?.id)
    } catch (cause) {
      setError(message(cause))
    } finally {
      setLoading(false)
    }
  }

  const chooseScenario = (next: Scenario) => {
    setScenarioId(next.id)
    setToken(next.token)
    void runAnalysis(next.id, next.token)
  }

  return (
    <main>
      <header className="topbar">
        <a className="brand" href="#top" aria-label="Token Blast Radius home">
          <span className="brand-mark">TBR</span>
          <span>Token Blast Radius</span>
        </a>
        <div className="system-status"><i /> analysis engine online</div>
      </header>

      <section className="hero" id="top">
        <div className="hero-copy">
          <span className="eyebrow">Identity security observatory · 01</span>
          <h1>See the breach<br /><em>before it happens.</em></h1>
          <p>Turn one signed credential into a complete map of every resource, trust boundary, and delegated identity it can reach.</p>
        </div>
        <div className="hero-index" aria-hidden="true">
          <span>TRACE</span>
          <strong>{analysis?.metrics.reachable_resources.toString().padStart(2, '0') ?? '00'}</strong>
          <small>reachable resources</small>
        </div>
      </section>

      <section className="scenario-section">
        <div className="section-intro">
          <span className="section-kicker">Choose a credential</span>
          <p>Each token carries a different identity chain, scope set, and delegation ceiling.</p>
        </div>
        <div className="scenario-grid">
          {scenarios.map((item, index) => (
            <button className={`scenario-card ${item.id === scenarioId ? 'is-active' : ''}`} key={item.id} onClick={() => chooseScenario(item)}>
              <span>0{index + 1}</span>
              <strong>{item.name}</strong>
              <p>{item.summary}</p>
              <i>{item.id === scenarioId ? 'tracing now' : 'trace credential'} →</i>
            </button>
          ))}
        </div>
      </section>

      {scenario && (
        <section className="question-band">
          <span>Question under investigation</span>
          <strong>{scenario.question}</strong>
          <button className="trace-button" onClick={() => void runAnalysis(scenario.id, token)} disabled={loading}>
            {loading ? 'Tracing paths' : 'Trace blast radius'}
          </button>
        </section>
      )}

      {error && <section className="error-banner"><span>Analysis stopped</span><strong>{error}</strong><p>A changed signature cannot be trusted. Restore the issued token to continue.</p></section>}

      {analysis && scenario && (
        <div className={`analysis-workspace ${loading ? 'is-loading' : ''}`}>
          <section className="summary-grid">
            <RiskDial metrics={analysis.metrics} />
            <div className="verdict panel">
              <span className="section-kicker">Machine verdict</span>
              <p>{analysis.verdict}</p>
            </div>
            <MetricStrip metrics={analysis.metrics} />
          </section>

          <ReachabilityGraph analysis={analysis} selectedNode={selectedNode} onSelectNode={setSelectedNode} />

          <section className="lower-grid">
            <TokenInspector claims={analysis.claims} token={token} originalToken={scenario.token} onTokenChange={setToken} />
            <NodeDetail node={node} path={path} incoming={incoming} labels={labels} />
          </section>

          <section className="learning-rail">
            <article><span>01</span><div><strong>Subject</strong><p>The `sub` claim names the agent holding the credential. Graph traversal begins here.</p></div></article>
            <article><span>02</span><div><strong>Actor</strong><p>The `act` claim preserves the accountable human behind the non-human identity.</p></div></article>
            <article><span>03</span><div><strong>Scopes</strong><p>Each relationship requires a scope. Missing scopes turn potential paths into blocked edges.</p></div></article>
            <article><span>04</span><div><strong>Depth</strong><p>The delegation ceiling limits how many relationship hops a credential may cross.</p></div></article>
          </section>
        </div>
      )}

      <footer>
        <span>Token Blast Radius</span>
        <p>Signed claims × authorization relationships = effective reach</p>
        <span>Go / React / TypeScript</span>
      </footer>
    </main>
  )
}

function message(value: unknown): string {
  return value instanceof Error ? value.message : 'Unexpected analysis failure'
}
