import type { AccessPath, Edge, Node } from '../types'

interface NodeDetailProps {
  node?: Node
  path?: AccessPath
  incoming?: Edge
  labels: Map<string, string>
}

export function NodeDetail({ node, path, incoming, labels }: NodeDetailProps) {
  if (!node) {
    return <section className="node-detail panel empty-detail">Select a node to inspect its access path.</section>
  }
  return (
    <section className="node-detail panel">
      <div className="detail-topline">
        <span className={`node-kind kind-${node.type}`}>{node.type}</span>
        <span className={`sensitivity sensitivity-${node.sensitivity}`}>{node.sensitivity}</span>
      </div>
      <h2>{node.label}</h2>
      <p>{node.description}</p>
      <div className={`reachability ${node.reachable ? 'is-reachable' : 'is-blocked'}`}>
        <span>{node.reachable ? 'reachable' : 'blocked'}</span>
        <strong>{node.reachable ? `${node.depth} delegation hops` : incoming?.blocked_reason ?? 'No active path'}</strong>
      </div>
      {path && (
        <div className="path-reader">
          <span className="section-kicker">Shortest active path</span>
          <div>
            {path.node_ids.map((id, index) => (
              <span className="path-step" key={id}>
                <b>{labels.get(id) ?? id}</b>
                {path.permissions[index] && <i>{path.permissions[index]}</i>}
              </span>
            ))}
          </div>
        </div>
      )}
      {incoming && <p className="rationale"><b>Why:</b> {incoming.rationale}</p>}
    </section>
  )
}
