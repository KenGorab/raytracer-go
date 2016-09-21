default: rayz

rayz: deps
	go build -o rayz

.PHONY: clean deps

clean:
	rm rayz

deps:
	govendor sync
