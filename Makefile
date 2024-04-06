
init_postgres:
	docker run --name postgres_avito -e POSTGRES_PASSWORD='test' -e POSTGRES_DB='avito' -e POSTGRES_USER='postgres'  -d -p 54325:5432 postgres

init_redis:
	docker run  --name redis_avito -d -p 6333:6379 redis:latest

init_jaeger:
	docker run -d --name jaeger_avito -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 14250:14250 -p 9411:9411 jaegertracing/all-in-one:1.22

first_build:
	make init_postgres
	make init_redis
	make init_jaeger