# build
.PHONY : build txstorm bootnode
build :
	go build -o build/lachesis ./cmd/lachesis

txstorm :
	go build -o build/tx-storm ./cmd/tx-storm

bootnode :
	go build -o build/bootnode ./cmd/bootnode

# dist
.PHONY : dist
dist :
	env GOOS=linux GOARCH=amd64 go build -o dist/lachesis ./cmd/lachesis	
	env GOOS=linux GOARCH=amd64 go build -o dist/tx-storm ./cmd/tx-storm
	env GOOS=linux GOARCH=amd64 go build -o dist/bootnode ./cmd/bootnode	

#test
.PHONY : test
test :
	go test ./...

#clean
.PHONY : clean
clean :
	rm ./build/lachesis ./build/tx-storm ./build/bootnode

clean-dist :
	rm -rf ./dist/lachesis ./dist/tx-storm ./dist/bootnode