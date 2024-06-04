# Args for building
ARG PORT=8080
ARG GO_VERSION=1.22.3
ARG DEFAULT_API_KEY="PROVIDE-APIKEY-ON-DEPLOY"
ARG LOG_FILE=1
ARG VERIFY_SESSION_IDENTITY=1
ARG BACKEND_ENDPOINT="sparkling-boldest-bridge.quiknode.pro"
ARG BACKEND_ENDPOINT_TOKEN="PROVIDE-TOKEN-ON-DEPLOY"
ARG BACKEND_USE_WEBSOCKET=1

FROM golang:${GO_VERSION}
LABEL authors="RuntimeRacer"

# Set required env vars
ENV ETHVAL_DEFAULT_API_KEY=${DEFAULT_API_KEY}
ENV ETHVAL_PORT=${PORT}
ENV ETHVAL_LOG_FILE=${LOG_FILE}
ENV ETHVAL_VERIFY_SESSION_IDENTITY=${VERIFY_SESSION_IDENTITY}
ENV ETHVAL_BACKEND_ENDPOINT=${BACKEND_ENDPOINT}
ENV BACKEND_ENDPOINT_TOKEN=${BACKEND_ENDPOINT_TOKEN}
ENV ETHVAL_BACKEND_USE_WEBSOCKET=${BACKEND_USE_WEBSOCKET}

# vend module; required for proper vendoring of dependencies which may contain non-golang files or modules (e.g. CGO dependencies)
RUN go install github.com/nomad-software/vend

WORKDIR /usr/src/app

# Clone Repo directly from github version to avoid contamination from local changes
# If you're developing, just copy comment this out and use the COPY command below insted
RUN git clone https://github.com/RuntimeRacer/ethereum-validator-go.git /tmp/ethereum-validator-go \
    && mv /tmp/ethereum-validator-go/validator-app/* . \
    && rm -fr /tmp/ethereum-validator-go
# COPY . .

# Ensure all the dependencies are properly set
RUN go mod tidy && go mod vender && vend

# Build the application
RUN go build -v -o /usr/local/bin/ethereum-validator-go ./...

# Open a local Port for incoming connections
EXPOSE ${PORT}

# Run the application
CMD ["ethereum-validator-go"]