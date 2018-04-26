NAME=ecs-deploy
DOCKER_ARGS=--name $(NAME) \
	--rm \
	-v "`pwd`/build":/var/task \
	-e DEBUG=true \
	-e AWS_ACCESS_KEY_ID \
	-e AWS_SECRET_ACCESS_KEY \
	-e AWS_SESSION_TOKEN \
	lambci/lambda:go1.x $(NAME)

build: 
	go build -o build/$(NAME) ./src
	
zip: build
	zip build/ecs-deploy build/ecs-deploy

clean:
	rm -rf build

test: clean build
	docker run $(DOCKER_ARGS) "`cat test/direct-invocation.json`"
	
testall: install
	for test in `ls test`; do echo "\n\n================\n$$test\n\n"; docker run $(DOCKER_ARGS) "`cat test/$$test`"; done
	
.PHONY: test testall