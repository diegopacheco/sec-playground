import { useMemo } from 'react'
import type { Analysis, Edge, Node } from '../types'

interface ReachabilityGraphProps {
  analysis: Analysis
  selectedNode?: string
  onSelectNode: (node: string) => void
}

interface Position {
  x: number
  y: number
}

const width = 1280
const height = 680
const nodeWidth = 166
const nodeHeight = 66

export function ReachabilityGraph({ analysis, selectedNode, onSelectNode }: ReachabilityGraphProps) {
  const positions = useMemo(() => layout(analysis), [analysis])
  return (
    <section className="graph-shell panel">
      <div className="graph-heading">
        <div>
          <span className="section-kicker">Live relationship map</span>
          <h2>Credential reachability</h2>
        </div>
        <div className="legend" aria-label="Graph legend">
          <span><i className="legend-line active-line" />reachable</span>
          <span><i className="legend-line blocked-line" />blocked</span>
          <span><i className="legend-dot critical-dot" />critical</span>
        </div>
      </div>
      <div className="graph-stage">
        <svg viewBox={`0 0 ${width} ${height}`} role="img" aria-label="Token relationship reachability graph">
          <defs>
            <marker id="arrow-active" markerWidth="8" markerHeight="8" refX="7" refY="4" orient="auto">
              <path d="M0,0 L8,4 L0,8 z" fill="#b8ff3d" />
            </marker>
            <marker id="arrow-blocked" markerWidth="8" markerHeight="8" refX="7" refY="4" orient="auto">
              <path d="M0,0 L8,4 L0,8 z" fill="#526069" />
            </marker>
            <filter id="active-glow" x="-50%" y="-50%" width="200%" height="200%">
              <feGaussianBlur stdDeviation="4" result="blur" />
              <feMerge><feMergeNode in="blur" /><feMergeNode in="SourceGraphic" /></feMerge>
            </filter>
          </defs>
          <g className="grid-lines">
            {Array.from({ length: 13 }, (_, i) => <line key={`v${i}`} x1={i * 100} y1="0" x2={i * 100} y2={height} />)}
            {Array.from({ length: 8 }, (_, i) => <line key={`h${i}`} x1="0" y1={i * 100} x2={width} y2={i * 100} />)}
          </g>
          <g className="graph-edges">
            {analysis.edges.map((edge, index) => {
              const from = positions.get(edge.from)
              const to = positions.get(edge.to)
              if (!from || !to) return null
              const path = curve(from, to)
              return (
                <g className={edge.active ? 'edge-active' : 'edge-blocked'} key={`${edge.from}-${edge.to}-${index}`}>
                  <path d={path} markerEnd={`url(#arrow-${edge.active ? 'active' : 'blocked'})`} />
                  <EdgeLabel from={from} to={to} edge={edge} />
                </g>
              )
            })}
          </g>
          <g className="graph-nodes">
            {analysis.nodes.map(node => {
              const position = positions.get(node.id)
              if (!position) return null
              const selected = node.id === selectedNode
              return (
                <g
                  key={node.id}
                  className={`graph-node ${node.reachable ? 'node-reachable' : 'node-blocked'} sensitivity-${node.sensitivity} ${selected ? 'node-selected' : ''}`}
                  transform={`translate(${position.x - nodeWidth / 2}, ${position.y - nodeHeight / 2})`}
                  onClick={() => onSelectNode(node.id)}
                  role="button"
                  tabIndex={0}
                  onKeyDown={event => { if (event.key === 'Enter' || event.key === ' ') onSelectNode(node.id) }}
                  aria-label={`${node.label}, ${node.reachable ? 'reachable' : 'blocked'}`}
                >
                  <rect width={nodeWidth} height={nodeHeight} rx="5" />
                  <circle cx="16" cy="17" r="4" />
                  <text className="node-type" x="28" y="21">{node.type}</text>
                  <text className="node-label" x="16" y="47">{fit(node.label)}</text>
                </g>
              )
            })}
          </g>
        </svg>
      </div>
    </section>
  )
}

function layout(analysis: Analysis): Map<string, Position> {
  const depths = structuralDepths(analysis)
  const groups = new Map<number, Node[]>()
  analysis.nodes.forEach(node => {
    const depth = depths.get(node.id) ?? 0
    groups.set(depth, [...(groups.get(depth) ?? []), node])
  })
  const maxDepth = Math.max(...groups.keys(), 1)
  const result = new Map<string, Position>()
  Array.from(groups.entries()).sort(([a], [b]) => a - b).forEach(([depth, nodes]) => {
    const x = 105 + depth * ((width - 210) / maxDepth)
    const ordered = [...nodes].sort((a, b) => Number(b.reachable) - Number(a.reachable) || a.label.localeCompare(b.label))
    ordered.forEach((node, index) => {
      const gap = height / (ordered.length + 1)
      result.set(node.id, { x, y: gap * (index + 1) })
    })
  })
  return result
}

function structuralDepths(analysis: Analysis): Map<string, number> {
  const depths = new Map<string, number>([[analysis.claims.act, 0], [analysis.claims.sub, 1]])
  for (let round = 0; round < analysis.nodes.length; round += 1) {
    analysis.edges.forEach(edge => {
      const from = depths.get(edge.from)
      if (from === undefined) return
      const candidate = from + 1
      const current = depths.get(edge.to)
      if (current === undefined || candidate < current) depths.set(edge.to, candidate)
    })
  }
  return depths
}

function curve(from: Position, to: Position): string {
  const startX = from.x + nodeWidth / 2
  const endX = to.x - nodeWidth / 2 - 4
  const control = Math.max(40, (endX - startX) * 0.52)
  return `M ${startX} ${from.y} C ${startX + control} ${from.y}, ${endX - control} ${to.y}, ${endX} ${to.y}`
}

function EdgeLabel({ from, to, edge }: { from: Position, to: Position, edge: Edge }) {
  const x = (from.x + to.x) / 2
  const y = (from.y + to.y) / 2 - 7
  return <text className="edge-label" x={x} y={y} textAnchor="middle">{edge.permission}</text>
}

function fit(label: string): string {
  return label.length > 21 ? `${label.slice(0, 19)}…` : label
}
