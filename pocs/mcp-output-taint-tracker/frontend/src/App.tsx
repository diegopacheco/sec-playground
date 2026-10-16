import { useCallback, useEffect, useState } from "react";
import { api } from "./api/client";
import { PlaygroundTab } from "./components/PlaygroundTab";
import { PolicyTab } from "./components/PolicyTab";
import { SessionsTab } from "./components/SessionsTab";
import type { Summary } from "./types";

type Tab = "sessions" | "playground" | "policy";

const TABS: { id: Tab; title: string }[] = [
  { id: "sessions", title: "Sessions" },
  { id: "playground", title: "Playground" },
  { id: "policy", title: "Policy" },
];

export function App() {
  const [tab, setTab] = useState<Tab>("sessions");
  const [sessions, setSessions] = useState<Summary[]>([]);
  const [selected, setSelected] = useState<string | null>(null);
  const [version, setVersion] = useState(0);
  const [online, setOnline] = useState(true);

  const refresh = useCallback(() => {
    api
      .sessions()
      .then((s) => {
        setSessions(s);
        setOnline(true);
        setVersion((v) => v + 1);
        setSelected((cur) => (cur && s.some((x) => x.id === cur) ? cur : (s[0]?.id ?? null)));
      })
      .catch(() => setOnline(false));
  }, []);

  useEffect(() => {
    refresh();
    const t = setInterval(refresh, 3000);
    return () => clearInterval(t);
  }, [refresh]);

  const loadScenarios = async () => {
    const loaded = await api.loadScenarios();
    refresh();
    if (loaded[0]) setSelected(loaded[0].id);
    setTab("sessions");
  };

  return (
    <div className="app">
      <header className="top">
        <img src="/logo.svg" alt="" width={34} height={34} />
        <div>
          <h1>MCP Output Taint Tracker</h1>
          <p className="sub">Tracks data from tool results and blocks it when it reaches a sensitive tool argument</p>
        </div>
        <span className="spacer" />
        <span className={`status ${online ? "up" : "down"}`}>{online ? "tracker online" : "tracker offline"}</span>
        <button className="primary" onClick={loadScenarios}>
          Load attack scenarios
        </button>
      </header>
      <nav className="tabs">
        {TABS.map((t) => (
          <button key={t.id} className={t.id === tab ? "on" : ""} onClick={() => setTab(t.id)}>
            {t.title}
          </button>
        ))}
      </nav>
      <main>
        {tab === "sessions" && <SessionsTab sessions={sessions} selected={selected} onSelect={setSelected} onChanged={refresh} version={version} />}
        {tab === "playground" && <PlaygroundTab sessions={sessions} onChanged={refresh} />}
        {tab === "policy" && <PolicyTab />}
      </main>
    </div>
  );
}
