# iron/go:dev is the alpine image with the go tools added
FROM iron/go:dev
WORKDIR /app
# Set an env var that matches your github repo name, replace treeder/dockergo here with your repo name
ENV SRC_DIR=/go/src/github.com/samuelagm/ccsi-tf/
# Add the source code:
ADD . $SRC_DIR
# Build it:
RUN cd $SRC_DIR; go get; go build -o main; cp main /app/
ENTRYPOINT ["./main"]