FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /rogueserver

USER 1000

# Load in arguments
ARG IS_DEBUG=false
ARG PROTOCOL=tcp
ARG ADDRESS=:8080
ARG MYSQL_DB_USERNAME=pokerogue
ARG MYSQL_DB_PASSWORD=
ARG MYSQL_DB_PROTOCOL=tcp
ARG MYSQL_DB_ADDRESS=localhost:3306
ARG MYSQL_DB_NAME=mysql

# Set ENVs based on ARGs
ENV IS_DEBUG=${IS_DEBUG} \
    PROTOCOL=${PROTOCOL} \
    ADDRESS=${ADDRESS} \
    MYSQL_DB_USERNAME=${MYSQL_DB_USERNAME} \
    MYSQL_DB_PASSWORD=${MYSQL_DB_PASSWORD} \
    MYSQL_DB_PROTOCOL=${MYSQL_DB_PROTOCOL} \
    MYSQL_DB_ADDRESS=${MYSQL_DB_ADDRESS} \
    MYSQL_DB_NAME=${MYSQL_DB_NAME}

# Start the app
CMD /rogueserver \
    -debug $IS_DEBUG \
    -proto $PROTOCOL \
    -addr $ADDRESS \
    -dbuser $MYSQL_DB_USERNAME \
    -dbpass $MYSQL_DB_PASSWORD \
    -dbproto $MYSQL_DB_PROTOCOL \
    -dbaddr $MYSQL_DB_ADDRESS \
    -dbname $MYSQL_DB_NAME