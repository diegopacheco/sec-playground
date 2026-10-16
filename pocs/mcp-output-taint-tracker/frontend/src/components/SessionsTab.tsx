import { useEffect, useState } from "react";
import { api } from "../api/client";
import type { Session, Summary } from "../types";
import { EventCard } from "./EventCard";
import { FlowGraph } from "./FlowGraph";

interface Props {
  sessions: Summary[];
  selected: string | null;
  onSelect: (id: string) => void;
  onChanged: () => void;
  version: number;
}

export function SessionsTab({ sessions, selected, onSelect, onChanged, version }: Props) {
  const [session, setSession] = useState<Session | null>(null);

  useEffect(() => {
    if (!selected) {
      setSession(null);
      return;
    }
    api.session(selected).then(setSession).catch(() => setSession(null));
  }, [selected, version]);

  const jump = (seq: number) => {
    const el = document.getElementById(`event-${seq}`);
    el?.scrollIntoView({ behavior: "smooth", block: "center" });
    el?.classList.add("flash");
    setTimeout(() => el?.classList.remove("flash"), 1200);
  };

  const remove = async (id: string) => {
    await api.remove(id);
    onChanged();
  };

  return (
    <div className="split">
      <aside className="list">
        {sessions.length === 0 && <p className="muted">No sessions yet. Load the attack scenarios or wire the hook into Claude Code.</p>}
        {sessions.map((s) => (
          <div key={s.id} className={`list-item ${s.id === selected ? "active" : ""}`} onClick={() => onSelect(s.id)}>
            <div className="list-name">{s.name}</div>
            <div className="list-meta">
              <span>{s.events} events</span>
              {s.blocked > 0 ? <span className="badge badge-block">{s.blocked} blocked</span> : <span className="badge badge-allow">clean</span>}
              <button
                className="link danger"
                onClick={(e) => {
                  e.stopPropagation();
                  remove(s.id);
                }}
              >
                delete
              </button>
            </div>
          </div>
        ))}
      </aside>
      <section className="detail">
        {!session && <p className="muted">Select a session to see its timeline and taint flows.</p>}
        {session && (
          <>
            <h2>{session.name}</h2>
            <FlowGraph events={session.events} onSelect={jump} />
            <div className="timeline">
              {session.events.map((e) => (
                <EventCard key={e.seq} event={e} onSource={jump} />
              ))}
            </div>
          </>
        )}
      </section>
    </div>
  );
}
