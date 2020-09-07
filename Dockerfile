# SERVICE BUILDER IMAGE
FROM golang:alpine AS binbuilder

# Build dependencies
RUN apk --no-cache --no-progress add gcc musl-dev make

RUN go version
COPY ./go.mod ./go.sum /uploader/
WORKDIR /uploader

# download deps before bringing in the sources
RUN go mod download

COPY . /uploader/
RUN go build -o /build/uploader /uploader/cmd/uploader

### ============================ ###

# RUNNER IMAGE
FROM alpine:latest

WORKDIR /uploader

# Copy binary and resources into runner image
COPY --from=binbuilder /build/uploader /bin/uploader
COPY ./assets /uploader/assets
COPY ./userlist.json /uploader/userlist.json

ENTRYPOINT /bin/uploader
EXPOSE 3000
