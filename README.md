# rabbitmq_demostration
a demostration for rabbit usage

# Build and run
## Install dependencies
it will generate the go.sum, and we dont need to generate the go.sum munually.
```
cd rabbitmq_demostration
go mod tidy
```

## Compile
```
go build -o bin/producer ./cmd/producer
go build -o bin/consumer ./cmd/consumer
```

## Start RabbitMQ
You should make sure that the rabbitmq service are running correctly.

## Run
```
./bin/consumer   # start consumer first
./bin/producer   # then start producer
```