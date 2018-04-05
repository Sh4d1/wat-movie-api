build:
	protoc -I/usr/local/include -I. \
		--go_out=plugins=micro:. \
		proto/movieapi/movieapi.proto
	docker build -t sh4d1/wat-movie-api .

run:
	docker run --net="host" \
		-p 50053 \
		-e MICRO_ADDRESS=":50053" \
		-e MICRO_REGISTRY="mdns" \
		-e OMDB_API_KEY="" \
		-e DEV=true \
		sh4d1/wat-movie-api
