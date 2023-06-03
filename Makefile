.PHONY: test
test:
	go test -race -cover .

.PHONY: bench
bench:
	go test -run=^$$ -cover -bench .

.PHONY: fuzz
fuzz:
	go test -run=^$$ -race -cover -fuzztime 1h -fuzz "FuzzStringParser" .

.PHONY: security
security:
	@echo "> Racing testing..."
	@go test -race -cover .
	@echo -e "\n> Racing benchmarks..."
	@go test -run=^$$ -race -cover -bench .
	@echo -e "\n> Running gosec..."
	@gosec .

.PHONY: clean
clean:
	rm -f *.out
	go clean
	go fmt
	go vet

.PHONY: profileEncoder
profileEncoder:
	go test -cpuprofile cpu.out -memprofile mem.out -bench "BenchmarkEncoder" .
	@echo "CPU Profile"
	go tool pprof --http localhost:18888 bencode.test cpu.out

.PHONY: profileStringParser
profileStringParser:
	go test -cpuprofile cpu.out -memprofile mem.out -bench "BenchmarkStringParser" .
	@echo "CPU Profile"
	go tool pprof --http localhost:18888 bencode.test cpu.out