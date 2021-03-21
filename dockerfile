FROM golang:alpine AS build

RUN apk update && apk add g++ gcc libxml2 libxslt-dev
COPY go.mod go.sum /go/src/app/

WORKDIR /go/src/app/
RUN go mod download

COPY . . 
RUN go install 

FROM alpine

RUN apk update && apk add libxml2 libxslt-dev
COPY --from=build /go/bin/spa-submissions-api /bin
COPY app.env /
RUN mkdir -p /schemas
COPY schemas/pain.001.001.03.xsd /schemas/pain.001.001.03.xsd
CMD ["spa-submissions-api"]