build:
    go build -o bin/optique main.go

test:
    cd test && go test -v
