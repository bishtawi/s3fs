FROM localstack/localstack:latest AS builder

ENV DATA_DIR=/opt/localstack-data

COPY init.sh .
COPY data data

RUN ./init.sh


FROM localstack/localstack:latest

ENV DATA_DIR=/opt/localstack-data \
    SERVICES=s3:4572

COPY --from=builder /opt/localstack-data /opt/localstack-data
