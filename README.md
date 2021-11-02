# BEFE the BackEnd For Everything
BEFE is a scriptable reverse proxy. It simplifies checking, rewriting and transforming incoming request through the help of a simple DSL.

The service is intended to run along side your backend service/api and act as the entry point to that service.

Some things this service can do
- conditionally reroute incoming requests to your backend services
- authorize request by checking JWT/JWK sets
- rewrite request to create new endpoints
- transform, add, filter fields from the response body
- enrich data by doing external lookups  
- dynamically reload your scripts if they change (without downtime)
- and much more

# Background and state of the project
I created this project to protect some of my backend api's with a JWT. 
There was a need to protect some endpoints to have additional permissions that I get from the scopes claim.

I also had the need to create virtual endpoints for my user where they can fetch data that is only intended for them.
For example i have an API that returns all the customers with their profile data. This is only intended to for the administrator.
By rewriting the request, filtering it with the user id extracted from the JWT and transforming the response to only export data 
that can only be exposed to the loggedin user. I created a new endpoint by just using my general purpose customers overview endpoint on my customers api.

This project is still very limited to the usecases I needed it for. It will probably grow as more come along. 

### Generate scripting engine bindings
Use the extract generator provided by yagni before you build
```
go generate ./...
```
In the root of the project would be enough

### Build and run in docker
To build a local docker image of the project, issue the following command
```
docker build -t befe:latest .
```

Example running with internal test server, with minimal resources
```
docker run -m=128m --memory-swap=0 --cpus="0.5" --rm -it -p 8083:8080 -v /Users/michael.boke/Projects/befe/examples/transform-filter-fields-in-response/:/script befe:latest --service http://localhost:8081 --path /script
```


## Run one of the examples
If you are in one of the example directories, execute the following command to load that program and start the instance of befe.
```
docker run -v ${PWD}:/script --rm -it docker.mbict.nl/mbict/befe:latest -addr=:8080 -path=/script
```

## Cross compile for docker
```
docker buildx build --platform linux/amd64,linux/arm64 -t mbict/befe:latest --push .
```