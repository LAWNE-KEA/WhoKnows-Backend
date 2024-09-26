FROM golang:latest
 
WORKDIR /app
 
# Effectively tracks changes within your go.mod file
COPY go.mod .
 
RUN go mod download
 
# Copies your source code into the app directory
COPY go.sum .
COPY main.go .

 
RUN go build -o /whoknows
 
EXPOSE 8080
 
CMD [ "/whoknows" ]