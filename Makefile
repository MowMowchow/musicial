.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/handleLogin ./cmd/routes/handleLogin/ 
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/handleRedirect ./cmd/routes/handleRedirect/ 
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/processLogin ./cmd/routes/processLogin/ 
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/pollLogin ./cmd/routes/pollLogin/ 
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/updateUser ./cmd/updateUser/ 
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/fetchUsers ./cmd/fetchs/fetchUsers/ 
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/fetchTracks ./cmd/fetchs/fetchTracks/ 
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/fetchArtists ./cmd/fetchs/fetchArtists/ 


clean:
	# rm -rf ./bin
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
