# ---------------------------------------------------------#
#                     Build web image                      #
# ---------------------------------------------------------#
FROM --platform=$BUILDPLATFORM node:16 as web

WORKDIR /usr/src/app

COPY web/package.json ./
COPY web/yarn.lock ./

# If you are building your code for production
# RUN npm ci --omit=dev

COPY ./web .

RUN yarn && yarn build && yarn cache clean

# ---------------------------------------------------------#
#                   Build SoloDev image                    #
# ---------------------------------------------------------#
FROM --platform=$BUILDPLATFORM golang:1.24.9-alpine3.22 as builder

RUN apk update \
    && apk add --no-cache protoc build-base git

# Setup workig dir
WORKDIR /app
RUN git config --global --add safe.directory '/app'

# Get dependencies - will also be cached if we won't change mod/sum
COPY go.mod .
COPY go.sum .

COPY Makefile .
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
RUN make dep
RUN make tools
# COPY the source code as the last step
COPY . .

COPY --from=web /usr/src/app/dist /app/web/dist

# build
ARG GIT_COMMIT
ARG SOLODEV_VERSION_MAJOR
ARG SOLODEV_VERSION_MINOR
ARG SOLODEV_VERSION_PATCH
ARG GITNESS_VERSION_MAJOR
ARG GITNESS_VERSION_MINOR
ARG GITNESS_VERSION_PATCH
ARG TARGETOS TARGETARCH

RUN if [ "$TARGETARCH" = "arm64" ]; then \
    wget -P ~ https://musl.cc/aarch64-linux-musl-cross.tgz && \
    tar -xvf ~/aarch64-linux-musl-cross.tgz -C ~ ; \
    fi

# set required build flags
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    if [ "$TARGETARCH" = "arm64" ]; then CC=~/aarch64-linux-musl-cross/bin/aarch64-linux-musl-gcc; fi && \
    VERSION_MAJOR=${SOLODEV_VERSION_MAJOR:-$GITNESS_VERSION_MAJOR} && \
    VERSION_MINOR=${SOLODEV_VERSION_MINOR:-$GITNESS_VERSION_MINOR} && \
    VERSION_PATCH=${SOLODEV_VERSION_PATCH:-$GITNESS_VERSION_PATCH} && \
    LDFLAGS="-X github.com/EolaFam1828/SoloDev/version.GitCommit=${GIT_COMMIT} -X github.com/EolaFam1828/SoloDev/version.major=${VERSION_MAJOR} -X github.com/EolaFam1828/SoloDev/version.minor=${VERSION_MINOR} -X github.com/EolaFam1828/SoloDev/version.patch=${VERSION_PATCH} -extldflags '-static'" && \
    CGO_ENABLED=1 \
    GOOS=$TARGETOS GOARCH=$TARGETARCH \
    CC=$CC go build -ldflags="$LDFLAGS" -o ./solodev ./cmd/solodev

### Pull CA Certs
FROM --platform=$BUILDPLATFORM alpine:latest as cert-image

RUN apk --update add ca-certificates

# ---------------------------------------------------------#
#                   Create final image                     #
# ---------------------------------------------------------#
FROM --platform=$TARGETPLATFORM alpine/git:2.49.1 as final

# setup app dir and its content
WORKDIR /app
VOLUME /data

ENV XDG_CACHE_HOME /data
ENV SOLODEV_GIT_ROOT /data
ENV SOLODEV_REGISTRY_FILESYSTEM_ROOT_DIRECTORY /data/registry
ENV SOLODEV_DATABASE_DRIVER sqlite3
ENV SOLODEV_DATABASE_DATASOURCE /data/database.sqlite
ENV SOLODEV_METRIC_ENABLED=false
ENV SOLODEV_TOKEN_COOKIE_NAME=token
ENV SOLODEV_DOCKER_API_VERSION 1.41
ENV SOLODEV_SSH_ENABLE=true
ENV SOLODEV_GITSPACE_ENABLE=true

COPY --from=builder /app/solodev /app/solodev
COPY --from=cert-image /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
RUN ln -sf /app/solodev /app/gitness

EXPOSE 3000
EXPOSE 3022

ENTRYPOINT [ "/app/solodev", "server" ]
