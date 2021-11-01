# CNS

Emeris configuration service.  
Allows admins to add and configure supported chains and tokens.

## Actions

* `make`  
Build and generate a binary.

* `make generate-swagger`  
Generate `swagger.yaml` under `cns/docs`.
Alternatively, you can get a generated copy as a [Github action artifact](https://github.com/allinbits/emeris-cns-server/actions/workflows/test.yml?query=workflow%3A%22Generate+Swagger%22).  

## Dependencies & Licenses

The list of non-{Cosmos, AiB, Tendermint} dependencies and their licenses are:

|Module   	                  |License          |
|---	                      |---  	        |
|go-playground/validator   	  |MIT   	        |
|sigs.k8s.io/controller-runtime |MIT            |
|go.uber.org/zap   	          |MIT           	|
|stretchr/testify   	      |MIT           	|
|go-redis/redis   	          |BSD-2 Simple    	|
|gin-contrib/zap   	          |MIT    	        |
|lib/pq                       |Open use         |

