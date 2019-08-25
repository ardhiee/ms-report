# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc

ADD . /go/src/ms-report
WORKDIR /go/src/ms-report
RUN go get ms-report
RUN go install
RUN cd /go/src/ms-report && go build -o ms-report

#tambahan http server
# WORKDIR /go/src/ms-report/http
# RUN go get http
# RUN go install
# RUN cd /go/src/ms-report/http && go build -o http-server

# final stage
FROM alpine
WORKDIR /go/src/ms-report
COPY --from=build-env /go/src/ms-report /go/src/ms-report
#COPY --from=build-env /go/src/ms-report/http /go/src/http-server
ENTRYPOINT ./ms-report
EXPOSE 7000
