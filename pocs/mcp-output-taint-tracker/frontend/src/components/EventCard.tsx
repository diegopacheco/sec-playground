import { argText } from "../lib/highlight";
import type { Flow, TaintEvent } from "../types";
import { LabelBadge, VerdictBadge } from "./Badge";
import { Highlight } from "./Highlight";

interface Props {
  event: TaintEvent;
  onSource?: (seq: number) => void;
}

export function EventCard({ event, onSource }: Props) {
  return (
    <article id={`event-${event.seq}`} className={`card card-${event.kind} ${event.decision && !event.decision.allow ? "card-blocked" : ""}`}>
      <header className="card-head">
        <span className="seq">#{event.seq}</span>
        <span className="kind">{kindTitle(event)}</span>
        {event.tool && <code className="tool">{event.tool}</code>}
        <span className="spacer" />
        {event.labels?.map((l) => <LabelBadge key={l} label={l} />)}
        {event.decision && <VerdictBadge allow={event.decision.allow} />}
      </header>
      {event.kind === "prompt" && <Highlight text={event.text ?? ""} marks={[]} />}
      {event.kind === "result" && <ResultBody event={event} />}
      {event.kind === "call" && <CallBody event={event} onSource={onSource} />}
    </article>
  );
}

function kindTitle(e: TaintEvent): string {
  if (e.kind === "prompt") return "user prompt (trusted)";
  if (e.kind === "result") return e.labels?.length ? "tool result (source)" : "tool result (clean)";
  return "tool call";
}

function ResultBody({ event }: { event: TaintEvent }) {
  const entities = event.entities ?? [];
  const marks = entities.map((e) => ({ text: e.value, label: e.labels.includes("secret") ? ("secret" as const) : ("untrusted" as const) }));
  return (
    <>
      {event.args && <ArgsLine args={event.args} />}
      <Highlight text={event.text ?? ""} marks={marks} />
      {entities.length > 0 && (
        <div className="chips">
          {entities.map((e) => (
            <span key={e.kind + e.value} className={`chip chip-${e.labels.includes("secret") ? "secret" : "untrusted"}`}>
              {e.kind}: {e.value}
            </span>
          ))}
        </div>
      )}
    </>
  );
}

function CallBody({ event, onSource }: Props) {
  const d = event.decision;
  const flows = d?.flows ?? [];
  const args = event.args ?? {};
  return (
    <>
      <div className="args">
        {Object.entries(args).map(([k, v]) => {
          const checked = d?.checked.includes(k);
          const marks = flows.filter((f) => f.arg === k).map((f) => ({ text: f.evidence, label: f.label }));
          return (
            <div key={k} className="arg">
              <div className="arg-name">
                {k} {checked && <span className="sink">sink</span>}
              </div>
              <Highlight text={argText(v)} marks={marks} />
            </div>
          );
        })}
      </div>
      {flows.length > 0 && (
        <ul className="flows">
          {flows.map((f, i) => (
            <FlowLine key={i} flow={f} onSource={onSource} />
          ))}
        </ul>
      )}
    </>
  );
}

function FlowLine({ flow, onSource }: { flow: Flow; onSource?: (seq: number) => void }) {
  return (
    <li className={`flow flow-${flow.label}`}>
      <LabelBadge label={flow.label} /> data from{" "}
      <button className="link" onClick={() => onSource?.(flow.sourceId)}>
        #{flow.sourceId} {flow.sourceTool}
      </button>{" "}
      reached argument <code>{flow.arg}</code> by {flow.kind === "text" ? "copied text" : flow.entityKind} match
      {flow.encoding !== "plain" && (
        <>
          {" "}
          after <strong>{flow.encoding}</strong> decoding: <code>{flow.decoded}</code>
        </>
      )}
      {flow.encoding === "plain" && (
        <>
          : <code>{flow.evidence}</code>
        </>
      )}
    </li>
  );
}

function ArgsLine({ args }: { args: Record<string, unknown> }) {
  return (
    <div className="args-line">
      {Object.entries(args).map(([k, v]) => (
        <span key={k}>
          {k}=<code>{argText(v)}</code>
        </span>
      ))}
    </div>
  );
}
