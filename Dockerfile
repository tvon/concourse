# syntax=docker/dockerfile:experimental

# NOTE: this Dockerfile is purely for local development! it is *not* used for
# the official 'concourse/concourse' image.

FROM concourse/dev

# build Concourse without using 'packr' and set up a volume so the web assets
# live-update
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
      go build -gcflags=all="-N -l" -o /usr/local/concourse/bin/concourse \
      ./cmd/concourse
VOLUME /src

# generate keys (with 1024 bits just so they generate faster)
RUN mkdir -p /concourse-keys
RUN concourse generate-key -t rsa -b 1024 -f /concourse-keys/session_signing_key
RUN concourse generate-key -t ssh -b 1024 -f /concourse-keys/tsa_host_key
RUN concourse generate-key -t ssh -b 1024 -f /concourse-keys/worker_key
RUN cp /concourse-keys/worker_key.pub /concourse-keys/authorized_worker_keys
