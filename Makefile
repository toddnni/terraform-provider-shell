TEST?=$$(go list ./...)

test: fmtcheck
	TF_ACC=1 go test -v $(TEST)

install:
	go install

fmtcheck:
	l=`gofmt -l .`; if [ -n "$$l" ]; then echo "Following needs formatting (gofmt):"; echo "$$l"; exit 1; fi
