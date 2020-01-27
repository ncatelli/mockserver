ARG BASEIMG="alpine:3.11"
ARG BUILDIMG="golang:1.13.6-alpine3.11"
FROM $BUILDIMG as builder

ARG APP_NAME="mockserver"
ENV GOPATH=""

COPY . /go/

RUN cd /go && go build -o /${APP_NAME}

FROM $BASEIMG
LABEL maintainer="Nate Catelli <ncatelli@packetfire.org>"
LABEL description="Container for mockserver"

ARG SERVICE_USER="service"
ARG APP_NAME="mockserver"

RUN addgroup ${SERVICE_USER} && \
    adduser -D -G ${SERVICE_USER} ${SERVICE_USER}

COPY --from=builder /${APP_NAME} /opt/${APP_NAME}/bin/${APP_NAME}

RUN chown -R ${SERVICE_USER}:${SERVICE_USER} /opt/${APP_NAME}/bin/${APP_NAME} && \
    chmod +x /opt/${APP_NAME}/bin/${APP_NAME}

WORKDIR "/opt/${APP_NAME}/"
USER ${SERVICE_USER}

ENTRYPOINT [ "/opt/mockserver/bin/mockserver" ]
CMD [ "-h" ]
