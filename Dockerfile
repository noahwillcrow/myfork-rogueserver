FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /rogueserver

USER 1000

# Fix DNS issues with alpine
ENV ENABLE_ALPINE_PRIVATE_NETWORKING=true

# Load in arguments
ARG IS_DEBUG=false
ARG CORS_ALLOWED_ORIGINS=*
ARG PROTOCOL=tcp
ARG PORT=8080
ARG MYSQL_DB_USERNAME=pokerogue
ARG MYSQL_DB_PASSWORD=
ARG MYSQL_DB_PROTOCOL=tcp
ARG MYSQL_DB_ADDRESS=localhost:3306
ARG MYSQL_DB_NAME=mysql
ARG STAT_REFRESH_CRON_SPEC=@hourly

# Set ENVs based on ARGs
ENV IS_DEBUG=${IS_DEBUG} \
    CORS_ALLOWED_ORIGINS=${CORS_ALLOWED_ORIGINS} \
    PROTOCOL=${PROTOCOL} \
    ADDRESS=0.0.0.0:${PORT} \
    MYSQL_DB_USERNAME=${MYSQL_DB_USERNAME} \
    MYSQL_DB_PASSWORD=${MYSQL_DB_PASSWORD} \
    MYSQL_DB_PROTOCOL=${MYSQL_DB_PROTOCOL} \
    MYSQL_DB_ADDRESS=${MYSQL_DB_ADDRESS} \
    MYSQL_DB_NAME=${MYSQL_DB_NAME} \
    STAT_REFRESH_CRON_SPEC=${STAT_REFRESH_CRON_SPEC}

# Start the app
CMD sleep 3 && /rogueserver -debug=$IS_DEBUG -proto=$PROTOCOL -addr=$ADDRESS -dbuser=$MYSQL_DB_USERNAME -dbpass=$MYSQL_DB_PASSWORD -dbproto=$MYSQL_DB_PROTOCOL -dbaddr=$MYSQL_DB_ADDRESS -dbname=$MYSQL_DB_NAME -allowedorigins=$CORS_ALLOWED_ORIGINS -statrefreshcronspec=$STAT_REFRESH_CRON_SPEC