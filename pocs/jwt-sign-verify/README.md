### Run 

```bash
./run.sh
```

```
❯ ./run.sh
Signed JWT: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJJRCI6IjEyMzQtYWxsLWNhbi1zZWUifQ.QwGAj-LdE-LbMo6iHPNwuMQhPT5vyurbSWLgeQgosME
Verified JWT Payload: {
  sub: "1234567890",
  name: "John Doe",
  iat: 1516239022,
  ID: "1234-all-can-see",
}
JWT Info: {
  header: {
    alg: "HS256",
    typ: "JWT",
  },
  payload: {
    sub: "1234567890",
    name: "John Doe",
    iat: 1516239022,
    ID: "1234-all-can-see",
  },
  signature: "QwGAj-LdE-LbMo6iHPNwuMQhPT5vyurbSWLgeQgosME",
}
```

### Decoding a token

on jwt.io

<img src="jwt-decode.png" />

### Recomended Reading

* JWT: https://jwt.io/
* What is JWT: https://www.youtube.com/watch?v=7Q17ubqLfaM&ab_channel=WebDevSimplified
* JWT Best Pratices: https://curity.io/resources/learn/jwt-best-practices/
* Stop using JWT for Sessions: http://cryto.net/~joepie91/blog/2016/06/13/stop-using-jwt-for-sessions/
* Why not? Go ahead and use JWTs for Authentication… https://medium.com/@jbyj/why-not-go-ahead-and-use-jwts-for-authentication-31810a4ce605