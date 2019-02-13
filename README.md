# Goblero Demo

Demo repo for [Goblero Job Queue](https://github.com/didil/goblero)

## Usage 

````
# Build
go get -u github.com/dgraph-io/badger
go get -u github.com/didil/goblero/pkg/blero
go build .
# For the first run
./goblero-demo enqueue
# To continue processing if interrupted
./goblero-demo
````