ARG  GO_BUILDER_IMAGE=golang:1.17-bullseye
ARG  NODE_BUILDER_IMAGE=node:16-bullseye
# ARG  DISTROLESS_IMAGE=gcr.io/distroless/static:nonroot
ARG  DISTROLESS_IMAGE=gcr.io/distroless/static:nonroot

############################
# STEP 1 Create infinityd dependency
############################
FROM ${GO_BUILDER_IMAGE} as go-base
RUN update-ca-certificates
WORKDIR /opt/mozaik
COPY golang/go.mod .
RUN go mod download
RUN go mod verify

############################
# STEP 2 Create infinity ui dependecy
############################
FROM ${NODE_BUILDER_IMAGE} as node-base
WORKDIR /opt/mozaik
COPY ui/package.json ./
COPY ui/yarn.lock ./
RUN yarn install

############################
# STEP 3 build infinityd
############################
FROM go-base as go-builder
COPY golang .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o mozaik ./cmd/mozaik/mozaik.go

############################
# STEP 4 Build infinity ui
############################
FROM node-base as node-builder
COPY ui .
RUN yarn run build

############################
# STEP 5 build infinity
############################
# using static nonroot image
# user:group is nobody:nobody, uid:gid = 65534:65534
FROM ${DISTROLESS_IMAGE}
WORKDIR /opt/mozaik
COPY --from=go-builder /opt/mozaik/mozaik .
COPY --from=node-builder /opt/mozaik/dist ./static
COPY assets /opt/mozaik/assets
CMD ["/opt/mozaik/mozaik"]

