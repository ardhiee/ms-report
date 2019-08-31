# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc

ADD . /go/src/ms-report
WORKDIR /go/src/ms-report
RUN go get ms-report
RUN go install
RUN cd /go/src/ms-report && go build -o ms-report

# final stage
FROM alpine
WORKDIR /go/src/ms-report
COPY --from=build-env /go/src/ms-report /go/src/ms-report
ENTRYPOINT ./ms-report
EXPOSE 7777
