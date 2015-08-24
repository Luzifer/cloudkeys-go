FROM scratch

VOLUME /data
EXPOSE 3000
ENTRYPOINT ["/cloudkeys-go"]
CMD ["--storage=local:////data", "--password-salt=changeme", "--username-salt=changeme"]

ADD https://rootcastore.hub.luzifer.io/v1/store/latest /etc/ssl/ca-bundle.pem
ADD ./cloudkeys-go /cloudkeys-go
ADD ./templates /templates
