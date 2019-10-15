build:
	go build -o bin/reverseproxy 

run:
	go run main.go ./config/config.yaml

clean:
	rm -f ./bin