default: rayz

rayz: deps
	go build

.PHONY: clean deps

clean:
	rm rayz

deps:
	govendor sync
