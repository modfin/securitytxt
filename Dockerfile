FROM golang:1.20-alpine3.17 as builder

COPY . /
RUN cd / && go build -o securitytxt securitytxt.go

FROM alpine:3.17
COPY --from=builder /securitytxt /
CMD /securitytxt



