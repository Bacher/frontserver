FROM debian
COPY build/frontserver /bin/frontserver
EXPOSE 7000
ENV DB_HOST "go-postgres"
ENV RPC_ADDR "0.0.0.0:9999"
ENTRYPOINT ["/bin/frontserver"]