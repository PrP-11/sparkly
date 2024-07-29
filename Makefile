network:
	docker network create service_mesh;

redis:
	docker run --name redis --rm --network=service_mesh -p 6379:6379 -d redis:6.0.9

images:
	docker pull golang:1.22.5
	docker pull wurstmeister/zookeeper
	docker pull wurstmeister/kafka

build:
	docker build . -t sparkly-services

dev:
	docker-compose -f docker-compose.yml up -d

run: images build dev

restart:
	docker restart sparkly-services-rest
	docker restart sparkly-services-worker

logs:
	docker logs sparkly-services -f

clean:
	docker kill sparkly-services

tests:
	go test -coverprofile coverage.out -covermode count ./internal/app/...

coverage/report: tests
	go tool cover -html=coverage.out

coverage: tests
	go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+'

mocks-cleanup:
	rm -rf mocks

mock-ports:
	mockery --dir=internal/ports/rest/handlers --name=Handler --filename=handlers.go --output=mocks/ports/rest/handlers --outpkg=handlers --with-expecter

mock-connectors:
	mockery --dir=internal/connectors/repository/mongo/logins --name=Repository --filename=logins.go --output=mocks/connectors/repository/mongo/logins --with-expecter
	mockery --dir=internal/connectors/repository/mongo/posts --name=Repository --filename=posts.go --output=mocks/connectors/repository/mongo/posts --with-expecter
	mockery --dir=internal/connectors/services/cache --name=CacheService --filename=cache.go --output=mocks/connectors/services/cache --with-expecter
	mockery --dir=internal/connectors/services/kafka --name=ProducerService --filename=producer.go --output=mocks/connectors/services/kafka --with-expecter
	mockery --dir=internal/connectors/services/kafka --name=ConsumerService --filename=consumer.go --output=mocks/connectors/services/kafka --with-expecter
	mockery --dir=internal/connectors/services/clock --name=Service --filename=clock.go --output=mocks/connectors/services/clock --with-expecter

mock-services:
	mockery --dir=internal/app/logins --name=Service --filename=logins.go --output=mocks/services/logins --with-expecter
	mockery --dir=internal/app/posts --name=Service --filename=posts.go --output=mocks/services/posts --with-expecter

mocks: mocks-cleanup mock-connectors mock-services mock-ports
	@echo "Done"