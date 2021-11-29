`docker-compose up`  

**Project setup**  
To load server on different port, create `.env` file with example value (HTTP_PORT=8182)


**How to run tests ?**  
To run tests, simply `cd` to root of project & run `go test .`
It will automatically set up `gin` test environment (no need to run api separately, just to cover tests)  