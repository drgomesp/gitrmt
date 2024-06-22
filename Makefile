build:
	go build -o git-remote-myremote examples/my-remote/main.go \
    && sudo mv git-remote-myremote /usr/local/bin/git-remote-myremote
