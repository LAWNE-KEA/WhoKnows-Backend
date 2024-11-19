FROM golang:latest
 
WORKDIR /app


COPY go.mod go.sum ./
 
RUN go mod download
 

COPY go.sum .
COPY schema.sql .
COPY main.go .

RUN go build -o /whoknows
 
EXPOSE 8080
 
CMD [ "/whoknows" ]