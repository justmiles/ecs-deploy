NAME=ecs-deploy
VERSION=`git tag | tail -1`
DATE=`date +"%Y%m%d_%H%M%S"`
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
	cd build && zip ecs-deploy-$(VERSION).zip ecs-deploy

clean:
	rm -rf build

test: clean build
	docker run $(DOCKER_ARGS) "`cat test/direct-invocation.json`"
	
invoke:
	mkdir -p lambda_output
	aws lambda invoke \
		--function-name "ecs-deploy" \
		--log-type "Tail" \
		--payload '{ "Application": "myapp", "Version": "latest", "Environment": "ops" }' lambda_output/$(DATE).log \
		| jq -r '.LogResult' | base64 -d

testall: install
	for test in `ls test`; do echo "\n\n================\n$$test\n\n"; docker run $(DOCKER_ARGS) "`cat test/$$test`"; done
	
.PHONY: test testall