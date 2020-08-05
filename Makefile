#############################
# test 
############################
test:
	GO111MODULE=on go vet ./...
	GO111MODULE=on go test -i .
	GO111MODULE=on go test -v -short -failfast -race -count=1 \
		    -coverpkg=./... \
		    -coverprofile coverage.txt \
		    -covermode=atomic \
		    ./...
	GO111MODULE=on go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage: file://$(PWD)/coverage.html"
integrations:
	GO111MODULE=on go test -i .
	GO111MODULE=on go test -v ./...

#############################
# build, migrateand run
############################
build:
	go build -mod=vendor -o bin/cart github.com/cubny/cart/cmd/...
  
migrate:
	./bin/cart -data ./data/cart.db -migrate

run:
	./bin/cart -data ./data/cart.db

firstrun: build migrate run

#############################
# build, migrateand run with docker
############################
docker-build:
	docker build -t cubny/cart .

docker-migrate:
	docker run --name cart -v `pwd`/data:/app/data -it --rm  cubny/cart app/bin/cart -data /app/data/cart.db -migrate

docker-run:
	docker run --name cart -v `pwd`/data:/app/data -it --rm  cubny/cart app/bin/cart -data /app/data/cart.db

docker-firstrun: docker-build docker-migrate docker-run
