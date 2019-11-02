# iron/go:dev is the alpine image with the go tools added
FROM iron/go:dev
WORKDIR /app
# Set an env var that matches your github repo name, replace treeder/dockergo here with your repo name
ENV SRC_DIR=/go/src/github.com/samuelagm/ccsi-tf/
ENV CONSUMER_KEY=vizBLoVyy7jCO6YodnZfPQ9uw
ENV CONSUMER_SECRET=cGqyg85zFBJNsQzSNPq1gRKGWoiF0tswk7cZVIYcx0QCK8hw6v
ENV ACCESS_TOKEN=145848213-IC11OPslrRVXRXQCE9v0YLJW1Yle5oLzlFANA5T6
ENV ACCESS_SECRET=ztMup6vGkUg7rqtwOejl9zKq2uKwMJ9QabwqMheTDsz5j
# Add the source code:
ADD . $SRC_DIR
# Build it:
RUN cd $SRC_DIR; go get; go build -o main; cp main /app/
ENTRYPOINT ["./main"]