FROM golang:1.10.3
WORKDIR /go/src/github.com/bgzzz/go-elevate/
RUN go get -d github.com/Masterminds/glide
RUN cd /go/src/github.com/Masterminds/glide && git checkout v0.12.3
RUN cd /go/src/github.com/Masterminds/glide && go install 
COPY . /go/src/github.com/bgzzz/go-elevate/
RUN /go/bin/glide install
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM scratch
COPY --from=0 /go/src/github.com/bgzzz/go-elevate/go-elevate /.
ENTRYPOINT ["/go-elevate"]
