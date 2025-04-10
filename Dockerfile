FROM golang:latest

ARG NAME=${NAME}
ENV NAME=${NAME}
ARG DB_PATH=${DB_PATH}
ENV DB_PATH=${DB_PATH}
ARG DEV_MODE=${DEV_MODE}
ENV DEV_MODE=${DEV_MODE}
ARG INSTANCE_ID=${INSTANCE_ID}
ENV INSTANCE_ID=${INSTANCE_ID}
ARG LOG_PATH=${LOG_PATH}
ENV LOG_PATH=${LOG_PATH}
ARG PORT=${PORT}
ENV PORT=${PORT}
ARG SECRET_KEY=${SECRET_KEY}
ENV SECRET_KEY=${SECRET_KEY}
ARG DISCORD_WEBHOOK=${DISCORD_WEBHOOK}
ENV DISCORD_WEBHOOK=${DISCORD_WEBHOOK}

WORKDIR /app

# Copy go.mod and go.sum files first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/notifapi

COPY create_config.sh /app/create_config.sh
RUN chmod +x /app/create_config.sh
RUN /app/create_config.sh

CMD ["/app/notifapi"]
EXPOSE ${PORT}