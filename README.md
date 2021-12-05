**How to run code ?**  
- directly using `go run .`
- using `docker compose up`, it will by default listen on ` :8081`

`.env` file configures value of gin (debug/test/release), it’s set to `release` mode, so http routes should be logged to `logs/gin.log` in `release` mode, otherwise, they are displayed in console in `debug` mode

**Example routes:**  
- http://localhost:8081/api/v1/block/BTC/000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf
- http://localhost:8081/api/v1/tx/BTC/dbaf14e1c476e76ea05a8b71921a46d6b06f0a950f17c5f9f1a03b8fae467f10

**Documentation:** 
- Swagger docs could be viewed at http://localhost:8081/api/v1/swagger/index.html

**How to run tests ?**  
`go test `  
All tests can be run without actually launching the API on separate port. It’ll autoconfigure `gin` to `test` mode and will handle routes. 
Alternatively, we can set `gin` framework's `EnvGinMode` to `test` 

**Caching:**  
Currently, caching is implemented through https://github.com/patrickmn/go-cache, so repeated requests are returning much faster (without actually generating full response). Other options are `redis / mecached / Varnish`

**Logging:**  
There is a single instance of logger being used throughout app & it will log to `logs/app.log`. Currently unexpected API responses are logged.
logrotate sholud be considered for large file. On production environment, we should use a more proper solution e.g. Sentry

**Http Client:**  
`Timeout ` is configured, so if external APi doesn’t respond in time, we should be able to log the error and not wait for it forever.

**Deployment:**  
- Docker supports `Automated docker builds`. Other options is to use `Jenkins` 

**Security:**  
- Currently, validation is applied on input values. In all other cases, all inputs should be filtered before performing any operation.
- API should be protected against DDOS

**Scalability:**
- “API Rate limiting” should be applied to avoid unwanted load on backend
- if the amount of network traffic increases, we can run multiple instances of service behind a load balancer.