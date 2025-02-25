FROM golang:latest

WORKDIR /app

COPY . .

RUN go build

COPY create_config.sh /app/create_config.sh
RUN chmod +x /app/create_config.sh
RUN /app/create_config.sh

CMD ["notifapi"]
EXPOSE ${PORT}