build: cmd/main.go
	go build -o weather-go cmd/main.go

run: build
	./weather-go

clean: weather-go
	rm weather-go
