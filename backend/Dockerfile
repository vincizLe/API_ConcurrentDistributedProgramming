FROM golang:alpine
#
WORKDIR /app
# renombrar
COPY ./main.go ./api-svc01.go
COPY go.mod ./
COPY go.sum ./

COPY ./api/ /usr/local/go/src/enigma/api/
COPY ./helpers/ /usr/local/go/src/enigma/helpers/
COPY ./models/ /usr/local/go/src/enigma/models/
# ADD api enigma/api
# ADD helpers /app/
#
RUN go mod download
#port
EXPOSE 8080
# run api-svc01.go
CMD ["go", "run", "api-svc01.go"]

