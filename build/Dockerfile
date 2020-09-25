FROM golang:1.15
MAINTAINER farir1881@gmail.com

RUN mkdir /app

COPY /bin/calendar /app

ADD /configs /app/configs
RUN chmod +x /app/calendar

WORKDIR /app
EXPOSE 8091 8092
ENTRYPOINT ["./calendar"]