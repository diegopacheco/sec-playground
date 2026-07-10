import * as jwt from 'jsonwebtoken';

export function printJwtInfo(token: string): void {
  const decoded = jwt.decode(token, { complete: true });
  console.log('JWT Info:', decoded);
}

printJwtInfo("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJJRCI6IjEyMzQtYWxsLWNhbi1zZWUifQ.QwGAj-LdE-LbMo6iHPNwuMQhPT5vyurbSWLgeQgosME");