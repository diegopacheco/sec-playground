export type GeneratorOptions = { length: number; lower: boolean; upper: boolean; digits: boolean; symbols: boolean };

const SETS = {
  lower: "abcdefghijkmnopqrstuvwxyz",
  upper: "ABCDEFGHJKLMNPQRSTUVWXYZ",
  digits: "23456789",
  symbols: "!@#$%^&*()-_=+[]{};:,.?",
};

export const DEFAULT_OPTIONS: GeneratorOptions = { length: 20, lower: true, upper: true, digits: true, symbols: true };

function randomIndex(max: number): number {
  const limit = Math.floor(0x100000000 / max) * max;
  const buf = new Uint32Array(1);
  for (;;) {
    crypto.getRandomValues(buf);
    if (buf[0]! < limit) return buf[0]! % max;
  }
}

export function generatePassword(opts: GeneratorOptions): string {
  const sets = (Object.keys(SETS) as (keyof typeof SETS)[]).filter((k) => opts[k]).map((k) => SETS[k]);
  if (sets.length === 0) throw new Error("pick at least one character set");
  if (opts.length < sets.length || opts.length > 128) throw new Error("length out of range");
  const all = sets.join("");
  const chars = sets.map((s) => s[randomIndex(s.length)]!);
  while (chars.length < opts.length) chars.push(all[randomIndex(all.length)]!);
  for (let i = chars.length - 1; i > 0; i--) {
    const j = randomIndex(i + 1);
    [chars[i], chars[j]] = [chars[j]!, chars[i]!];
  }
  return chars.join("");
}

export function entropyBits(opts: GeneratorOptions): number {
  const size = (Object.keys(SETS) as (keyof typeof SETS)[]).filter((k) => opts[k]).reduce((n, k) => n + SETS[k].length, 0);
  return size === 0 ? 0 : Math.round(opts.length * Math.log2(size));
}
