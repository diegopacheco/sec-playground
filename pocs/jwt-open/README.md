### Run 

```bash
./run.sh
```

```
❯ ./run.sh
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

### Recomended Reading

* JWT: https://jwt.io/
* JWT Best Pratices: https://curity.io/resources/learn/jwt-best-practices/
* Stop using JWT for Sessions: http://cryto.net/~joepie91/blog/2016/06/13/stop-using-jwt-for-sessions/
* Why not? Go ahead and use JWTs for Authentication… https://medium.com/@jbyj/why-not-go-ahead-and-use-jwts-for-authentication-31810a4ce605