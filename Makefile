.PHONY: test
test: testdata/fixtures/.unpacked
	go test ./... -race -v

.PHONY: install-mockgen
install-mockgen:
	go install github.com/golang/mock/mockgen@v1.6.0

.PHONY: generate
generate:
	go generate ./...

.PHONY: get-ttar
bundle-ttar:
	wget https://raw.githubusercontent.com/ideaship/ttar/master/ttar -O ttar
	chmod a+x ttar

.PHONY: get-fixtures-archive
bundle-fixtures-archive:
	wget https://raw.githubusercontent.com/prometheus/procfs/5be723405f68d8381a9accdc0bd7bdbe7621e2b6/testdata/fixtures.ttar -O testdata/fixtures.ttar

%/.unpacked: %.ttar
	@echo ">> extracting fixtures $*"
	./ttar -C $(dir $*) -x -f $*.ttar
	touch $@

fixtures: testdata/fixtures/.unpacked

update_fixtures:
	rm -vf testdata/fixtures/.unpacked
	./ttar -c -f testdata/fixtures.ttar -C testdata/ fixtures/