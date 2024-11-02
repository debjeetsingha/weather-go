FILE_PATH := cmd/main.go

build: $(FILE_PATH)
	go build -o weather-go $(FILE_PATH)

execute: build
	./weather-go $(LOCATION)

run:
	go run $(FILE_PATH) $(LOCATION)

clean: weather-go
	rm weather-go
