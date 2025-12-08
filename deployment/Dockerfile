##############################################################
############################ AUTH ############################
##############################################################

# BUILDER
FROM golang:1.25.1-alpine AS auth-builder
WORKDIR /svc

COPY go.mod go.sum ./ 

COPY ./api ./api
COPY ./pkg ./pkg
COPY ./services/auth ./services/auth


RUN go mod download \
    && go mod tidy \
    && go build -o /svc/app ./services/auth/cmd


# RUNNER
FROM alpine:3.22 AS auth
WORKDIR /root/
COPY --from=auth-builder ./svc/app .
EXPOSE ${AUTH_SERVICE_PORT}
CMD ["./app"]

##############################################################
########################## GATEWAY ###########################
##############################################################

# BUILDER
FROM golang:1.25.1-alpine AS gateway-builder
WORKDIR /svc

COPY go.mod go.sum ./ 

COPY ./api ./api
COPY ./pkg ./pkg
COPY ./services/gateway ./services/gateway


RUN go mod download \
    && go mod tidy \
    && go build -o /svc/app ./services/gateway/cmd


# RUNNER
FROM alpine:3.22 AS gateway
WORKDIR /root/
COPY --from=gateway-builder ./svc/app .
EXPOSE ${AUTH_SERVICE_PORT}
CMD ["./app"]
