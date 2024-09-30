FROM golang:latest
 
WORKDIR /app

ENV ENV_MYSQL_USER=root
ENV ENV_MYSQL_PASSWORD=root

# Effectively tracks changes within your go.mod file
COPY go.mod .
 
RUN go mod download
 
# Copies your source code into the app directory
COPY go.sum .
COPY schema.sql .
COPY main.go .

RUN go build -o /whoknows
 
EXPOSE 8080
 
CMD [ "/whoknows" ]