FROM scratch

ARG COMPONENT_NAME=registry
ARG TARGETOS=linux
ARG TARGETARCH=amd64

COPY --chmod=0755 dist/${COMPONENT_NAME}_${TARGETOS}_${TARGETARCH}*/${COMPONENT_NAME} /service

EXPOSE 50051

ENTRYPOINT ["/service"]
