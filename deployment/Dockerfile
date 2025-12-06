##############################################################
############################ AUTH ############################
##############################################################

ARG auth

# BUILDER
FROM golang:1.25.1-alpine AS builder
WORKDIR /svc

COPY go.mod go.sum ./ 

COPY ./api ./api
COPY ./pkg ./pkg
COPY ./services/auth ./services/auth


RUN go mod download \
    && go mod tidy \
    && go build -o /svc/app ./services/auth/cmd


# RUNNER
FROM alpine:3.22 AS runner
WORKDIR /root/
COPY --from=builder ./svc/app .
EXPOSE ${AUTH_SERVICE_PORT}
CMD ["./app"]