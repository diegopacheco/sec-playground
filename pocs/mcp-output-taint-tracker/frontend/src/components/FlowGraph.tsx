import type { Label, TaintEvent } from "../types";

const ROW = 44;
const WIDTH = 760;
const BOX_W = 250;
const BOX_H = 32;
const COLORS: Record<Label, string> = { untrusted: "#d97706", secret: "#7c3aed" };

interface Props {
  events: TaintEvent[];
  onSelect: (seq: number) => void;
}

export function FlowGraph({ events, onSelect }: Props) {
  const sources = events.filter((e) => e.kind === "result" && e.labels?.length);
  const calls = events.filter((e) => e.kind === "call");
  if (sources.length === 0 && calls.length === 0) return null;
  const rows = Math.max(sources.length, calls.length);
  const height = rows * ROW + 40;
  const leftY = new Map(sources.map((s, i) => [s.seq, 30 + i * ROW]));
  const rightY = new Map(calls.map((c, i) => [c.seq, 30 + i * ROW]));
  const edges = calls.flatMap((c) =>
    (c.decision?.flows ?? []).map((f) => ({ from: f.sourceId, to: c.seq, label: f.label, key: `${c.seq}-${f.sourceId}-${f.label}-${f.arg}` })),
  );
  return (
    <div className="graph">
      <div className="graph-titles">
        <span>sources</span>
        <span>tool calls</span>
      </div>
      <svg viewBox={`0 0 ${WIDTH} ${height}`} role="img" aria-label="taint flow graph">
        {edges.map((e) => {
          const y1 = (leftY.get(e.from) ?? 0) + BOX_H / 2;
          const y2 = (rightY.get(e.to) ?? 0) + BOX_H / 2;
          const x1 = BOX_W + 10;
          const x2 = WIDTH - BOX_W - 10;
          const mid = (x1 + x2) / 2;
          return (
            <path key={e.key} d={`M${x1},${y1} C${mid},${y1} ${mid},${y2} ${x2},${y2}`} stroke={COLORS[e.label]} strokeWidth={2.5} fill="none" markerEnd={`url(#arrow-${e.label})`} />
          );
        })}
        <defs>
          {(Object.keys(COLORS) as Label[]).map((l) => (
            <marker key={l} id={`arrow-${l}`} viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto">
              <path d="M0,0 L10,5 L0,10 z" fill={COLORS[l]} />
            </marker>
          ))}
        </defs>
        {sources.map((s) => (
          <Box key={s.seq} x={10} y={leftY.get(s.seq)!} text={`#${s.seq} ${s.tool}`} tone={s.labels?.includes("secret") ? "secret" : "untrusted"} onClick={() => onSelect(s.seq)} />
        ))}
        {calls.map((c) => (
          <Box key={c.seq} x={WIDTH - BOX_W - 10} y={rightY.get(c.seq)!} text={`#${c.seq} ${c.tool}`} tone={c.decision?.allow ? "allow" : "block"} onClick={() => onSelect(c.seq)} />
        ))}
      </svg>
    </div>
  );
}

function Box({ x, y, text, tone, onClick }: { x: number; y: number; text: string; tone: string; onClick: () => void }) {
  const label = text.length > 32 ? text.slice(0, 31) + "…" : text;
  return (
    <g className={`gbox gbox-${tone}`} onClick={onClick}>
      <rect x={x} y={y} width={BOX_W} height={BOX_H} rx={8} />
      <text x={x + 12} y={y + 21}>
        {label}
      </text>
    </g>
  );
}
