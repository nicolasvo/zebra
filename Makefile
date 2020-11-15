.PHONY: build clean

build:
	export GO111MODULE=on
	@if [ "$(target)" == "arm" ]; then\
		env GOOS=linux GOARCH=arm GOARM=7 go build -v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o zebra ./main.go;\
	else\
		env GOOS=linux GOARCH=amd64 go build -v -o zebra ./main.go;\
	fi

clean:
	rm -f zebra
