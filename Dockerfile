FROM golang 

COPY main $APP_DIR/app

EXPOSE 8080

ENV GIN_MODE=release

# to set permissions of added files/artifacts to app user and perform further hardening
RUN $APP_DIR/post-install.sh

RUN chmod +x $APP_DIR/app

COPY etc/config.env $APP_DIR/etc/config.env
# to run as unprivileged user
USER $APP_USER

CMD ["./app"]
