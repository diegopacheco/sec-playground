import type { Label } from "../types";

export interface Mark {
  text: string;
  label: Label;
}

export interface Segment {
  text: string;
  label?: Label;
}

export function segments(text: string, marks: Mark[]): Segment[] {
  const lower = text.toLowerCase();
  const ranges: { start: number; end: number; label: Label }[] = [];
  for (const m of marks) {
    if (!m.text) continue;
    const needle = m.text.toLowerCase();
    let from = 0;
    for (;;) {
      const i = lower.indexOf(needle, from);
      if (i < 0) break;
      ranges.push({ start: i, end: i + needle.length, label: m.label });
      from = i + needle.length;
    }
  }
  ranges.sort((a, b) => a.start - b.start || b.end - a.end);
  const out: Segment[] = [];
  let pos = 0;
  for (const r of ranges) {
    if (r.start < pos) continue;
    if (r.start > pos) out.push({ text: text.slice(pos, r.start) });
    out.push({ text: text.slice(r.start, r.end), label: r.label });
    pos = r.end;
  }
  if (pos < text.length) out.push({ text: text.slice(pos) });
  return out;
}

export function argText(value: unknown): string {
  return typeof value === "string" ? value : JSON.stringify(value);
}
