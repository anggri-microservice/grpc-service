FROM golang:1.13-alpine
ARG SRV_UID=1000
ARG SRV_GID=1000
ENV SERVICE_NAME grpc-service

RUN apk add --no-cache sudo git openssh-client make postgresql-client curl tzdata

COPY build/assets/sudoers /etc/sudoers

RUN chmod 0400 /etc/sudoers && \
    addgroup anggri -g ${SRV_GID} && \
    addgroup sudo && \
    adduser -h /home/anggri -s /bin/sh -u ${SRV_UID} -G users -D anggri && \
    adduser anggri anggri && \
    adduser anggri sudo

USER anggri
WORKDIR "/go/src/gitlab.com/anggri-microservice/${SERVICE_NAME}"

COPY id_rsa /tmp/id_rsa

COPY build/dev/evans_linux_amd64.tar.gz /tmp

RUN sudo chown -R anggri:anggri /go && \
    sudo chown anggri:anggri /tmp/id_rsa && \
    chmod 600 /tmp/id_rsa && \
    mkdir ~/.ssh && \
    mv /tmp/id_rsa ~/.ssh/id_rsa && \
    eval $(ssh-agent) && \
    ssh-add ~/.ssh/id_rsa && \
    echo "StrictHostKeyChecking no" | sudo tee -a /etc/ssh/ssh_config > /dev/null && \
    echo "Installing additional tools, please wait..." && \
    go get -u golang.org/x/lint/golint && echo "- Success installing golint" && \
    go get github.com/githubnemo/CompileDaemon && echo "- Success installing CompileDaemon" && \
    go get github.com/fullstorydev/grpcurl && \
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl && echo "- Success installing grpcurl" && \
    git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/" && \
    export GOPRIVATE=gitlab.com/anggri-microservice && \
    sudo tar -C /usr/bin -xzf /tmp/evans_linux_amd64.tar.gz && \
    sudo rm /tmp/evans_linux_amd64.tar.gz && \
    export GO111MODULE=off && \
    go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate && echo "- Success installing golang-migrate" && \
    go get github.com/bayucandra/go-seed-pg && echo "- Success installing go-seed-pg" && \
    unset GO111MODULE


ENTRYPOINT ["sh", "./build/dev/entrypoint.sh"]
