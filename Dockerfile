FROM scratch

COPY swarmkit-client /swarmkit-client

EXPOSE 8888
VOLUME /var/run/docker/cluster/docker-swarmd.sock

ENTRYPOINT ["/swarmkit-client"]