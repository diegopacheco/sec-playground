import { expect, test } from "bun:test";
import { fromBase64url, toBase64url } from "../src/base64url.ts";

test("encodes without padding and with url safe alphabet, as WebAuthn and the Go server expect", () => {
  expect(toBase64url(new Uint8Array([251, 255, 191]))).toBe("-_-_");
  expect(toBase64url(new Uint8Array([1]))).toBe("AQ");
});

test("round trips every byte value so challenges and credential ids are never corrupted", () => {
  const all = new Uint8Array(256).map((_, i) => i);
  expect(Array.from(fromBase64url(toBase64url(all)))).toEqual(Array.from(all));
});
