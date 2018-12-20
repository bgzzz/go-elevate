PREFIX = 
NAME = go-elevate

build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/$(NAME)
clean:
	rm ./bin/$(NAME)

dock-build:
	docker build -t $(PREFIX)$(NAME) .

dock-compose: dock-build
	docker-compose up

#rest is for testing 
