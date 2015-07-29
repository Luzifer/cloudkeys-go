FROM scratch

VOLUME /data
EXPOSE 3000
ENTRYPOINT ["/cloudkeys"]
CMD ["--storage=local:////data", "--password-salt=changeme", "--username-salt=changeme"]

ADD ./ca-certificates.pem /etc/ssl/ca-bundle.pem
ADD ./cloudkeys /cloudkeys
ADD ./templates /templates
