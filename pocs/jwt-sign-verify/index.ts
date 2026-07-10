import * as jwt from 'jsonwebtoken';

const secretKey = 'your-256-bit-secret';

export function printJwtInfo(token: string): void {
  const decoded = jwt.decode(token, { complete: true });
  console.log('JWT Info:', decoded);
}

export function signJwt(payload: object): string {
  return jwt.sign(payload, secretKey, { algorithm: 'HS256' });
}

export function verifyJwt(token: string): object | string {
  try {
    return jwt.verify(token, secretKey);
  } catch (err) {
    console.error('Invalid token:', err);
    throw err;
  }
}

// Example usage
const payload = { sub: '1234567890', name: 'John Doe', iat: 1516239022, ID: '1234-all-can-see' };
const token = signJwt(payload);
console.log('Signed JWT:', token);

const verifiedPayload = verifyJwt(token);
console.log('Verified JWT Payload:', verifiedPayload);

printJwtInfo(token);