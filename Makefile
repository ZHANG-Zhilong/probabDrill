PROJECT:=probabDrill

.PHONY: build
build:
	CGO_ENABLED=0 go build -o pd main.go
.PHONY: clean
clean:
	rm pd