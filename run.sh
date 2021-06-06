go mod tidy
go mod vendor
swag init
make build
./pd conf --path .