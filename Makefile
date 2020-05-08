.PHONY: build run openapi redoc

build:
	docker build -t excommerce .

run:
	docker run --rm -it -p 127.0.0.1:8080:8080/tcp excommerce:latest

openapi: api/openapi.yaml
	docker run --rm -v ${PWD}:/local \
		openapitools/openapi-generator-cli:v4.3.0 generate \
			--git-user-id Teelevision \
			--git-repo-id excommerce \
			-g go-server \
			-i /local/$< \
			-o /local

redoc:
	docker run -it --rm -p 8081:80 \
		-v ${PWD}/api/openapi.yaml:/usr/share/nginx/html/openapi.yaml:ro \
		-e SPEC_URL=openapi.yaml \
		redocly/redoc:v2.0.0-rc.28
