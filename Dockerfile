FROM alpine:3.15
COPY bin/docker-compose-deployer /usr/local/bin/
RUN apk add --no-cache curl docker-cli-compose
ENV DEPLOYER_COMPOSE_V2=true
ENTRYPOINT ["/usr/local/bin/docker-compose-deployer"]
CMD ["server"]