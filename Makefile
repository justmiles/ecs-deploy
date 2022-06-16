NAME=ecs-deploy
VERSION=v0.4.1
DATE=`date +"%Y%m%d_%H%M%S"`
TEST_JSON='{ "Application": "bender", "Version": "latest", "Environment": "ops" }'
DOCKER_ARGS=--name $(NAME) \
	--rm \
	-v "`pwd`/build":/var/task \
	-e DEBUG=true \
	-e AWS_ACCESS_KEY_ID \
	-e AWS_SECRET_ACCESS_KEY \
	-e AWS_SESSION_TOKEN \
	lambci/lambda:go1.x $(NAME)

build: clean
	go build -o build/$(NAME) ./src
	
zip: build
	cd build && zip lambda-ecs-deploy-$(VERSION).zip ecs-deploy && rm ecs-deploy

test-cli: clean
	go build -o build/$(NAME) ./src
	./build/ecs-deploy ship -a asdf -v latest -e ops --max-attempts 90 --debug

test: clean
	go build -o build/$(NAME) ./src
	docker run $(DOCKER_ARGS) $(TEST_JSON)

invoke:
	mkdir -p lambda_output
	aws lambda invoke \
		--function-name "ecs-deploy" \
		--log-type "Tail" \
		--payload $(TEST_JSON) lambda_output/$(DATE).log \
		| jq -r '.LogResult' | base64 -d

clean:
	rm -rf build

release:
	git tag -a $(VERSION) -m "release version $(VERSION)" && git push origin $(VERSION)
	goreleaser release --rm-dist

.PHONY: test testall