build:
	go build -o bin/reverseproxy 

run:
	go run main.go --config.file=./config/config.yaml

clean:
	rm -f ./bin