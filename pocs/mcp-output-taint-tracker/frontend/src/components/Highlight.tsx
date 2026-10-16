import { segments, type Mark } from "../lib/highlight";

export function Highlight({ text, marks }: { text: string; marks: Mark[] }) {
  return (
    <pre className="code">
      {segments(text, marks).map((s, i) =>
        s.label ? (
          <mark key={i} className={`mark-${s.label}`}>
            {s.text}
          </mark>
        ) : (
          <span key={i}>{s.text}</span>
        ),
      )}
    </pre>
  );
}
