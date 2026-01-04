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
########################### REPORT ###########################
##############################################################

# BUILDER
FROM golang:1.25.1-alpine AS report-builder
WORKDIR /svc

COPY go.mod go.sum ./ 

COPY ./api ./api
COPY ./pkg ./pkg
COPY ./services/report ./services/report


RUN go mod download \
    && go mod tidy \
    && go build -o /svc/app ./services/report/cmd


# RUNNER
FROM alpine:3.22 AS report
WORKDIR /root/
COPY --from=report-builder ./svc/app .
EXPOSE ${REPORT_SERVICE_PORT}
CMD ["./app"]


##############################################################
########################## GATEWAY ###########################
##############################################################

# BUILDER
FROM golang:1.25.1-alpine AS gateway-builder
WORKDIR /svc

COPY go.mod go.sum ./
RUN go mod download

COPY ./api ./api
COPY ./pkg ./pkg
COPY ./services/gateway ./services/gateway

RUN go install github.com/swaggo/swag/cmd/swag@latest \
 && cd /svc/services/gateway \
 && swag init -g cmd/main.go -d . -o ./docs --parseDependency --parseInternal
 

RUN  go build -o /svc/app ./services/gateway/cmd

# RUNNER
FROM alpine:3.22 AS gateway
WORKDIR /root/
COPY --from=gateway-builder /svc/app .
EXPOSE ${GATEWAY_SERVICE_PORT}

CMD ["./app"]

##############################################################
####################### REPORT-WORKER ########################
##############################################################
# BUILDER
FROM golang:1.25.1-alpine AS report-worker-builder
WORKDIR /svc

COPY go.mod go.sum ./ 

COPY ./api ./api
COPY ./pkg ./pkg
COPY ./services/report-worker ./services/report-worker


RUN go mod download \
    && go mod tidy \
    && go build -o /svc/app ./services/report-worker/cmd


# RUNNER
FROM alpine:3.22 AS report-worker
WORKDIR /root/


COPY --from=report-worker-builder ./svc/app .
COPY --from=report-worker-builder /svc/services/report-worker/internal/convert/fonts ./internal/convert/fonts
COPY --from=report-worker-builder /svc/services/report-worker/internal/convert/pgups_icon.png ./internal/convert/


CMD ["./app"]
