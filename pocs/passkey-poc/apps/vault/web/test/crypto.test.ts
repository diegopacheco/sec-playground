import { expect, test } from "bun:test";
import { checkCanary, deriveVaultKey, open, seal, sealCanary } from "../src/crypto.ts";
import { DEFAULT_OPTIONS, entropyBits, generatePassword } from "../src/generator.ts";

const prfA = crypto.getRandomValues(new Uint8Array(32)).buffer;
const prfB = crypto.getRandomValues(new Uint8Array(32)).buffer;

test("the same passkey PRF output always derives the same key, so the vault reopens", async () => {
  const sealed = await seal(await deriveVaultKey(prfA, "alice"), { password: "s3cret" }, "item-1");
  expect(await open(await deriveVaultKey(prfA, "alice"), sealed, "item-1")).toEqual({ password: "s3cret" });
});

test("a different passkey cannot decrypt the vault and the canary reports it", async () => {
  const canary = await sealCanary(await deriveVaultKey(prfA, "alice"));
  expect(await checkCanary(await deriveVaultKey(prfA, "alice"), canary)).toBe(true);
  expect(await checkCanary(await deriveVaultKey(prfB, "alice"), canary)).toBe(false);
});

test("ciphertext is bound to its item id, so the server cannot swap two entries", async () => {
  const key = await deriveVaultKey(prfA, "alice");
  const sealed = await seal(key, { password: "bank" }, "item-1");
  await expect(open(key, sealed, "item-2")).rejects.toThrow();
});

test("the key is not extractable, so XSS cannot export it as raw bytes", async () => {
  const key = await deriveVaultKey(prfA, "alice");
  expect(key.extractable).toBe(false);
  await expect(crypto.subtle.exportKey("raw", key)).rejects.toThrow();
});

test("each seal uses a fresh IV so equal passwords do not produce equal ciphertext", async () => {
  const key = await deriveVaultKey(prfA, "alice");
  const a = await seal(key, "same", "x");
  const b = await seal(key, "same", "x");
  expect(a.iv).not.toBe(b.iv);
  expect(a.ct).not.toBe(b.ct);
});

test("generated passwords honour length and contain every selected class", () => {
  for (let i = 0; i < 200; i++) {
    const p = generatePassword(DEFAULT_OPTIONS);
    expect(p).toHaveLength(20);
    expect(p).toMatch(/[a-z]/);
    expect(p).toMatch(/[A-Z]/);
    expect(p).toMatch(/[2-9]/);
    expect(p).toMatch(/[^a-zA-Z0-9]/);
  }
  expect(generatePassword({ ...DEFAULT_OPTIONS, symbols: false, upper: false, lower: false, length: 12 })).toMatch(/^[2-9]{12}$/);
  expect(entropyBits(DEFAULT_OPTIONS)).toBeGreaterThan(100);
  expect(() => generatePassword({ length: 10, lower: false, upper: false, digits: false, symbols: false })).toThrow();
});
