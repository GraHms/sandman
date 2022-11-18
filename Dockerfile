FROM docker.pkg.dev.vm.co.mz/alpine-base:2.0.1

COPY main $APP_DIR/app

EXPOSE 8080

ENV GIN_MODE=release

# to set permissions of added files/artifacts to app user and perform further hardening
RUN $APP_DIR/post-install.sh

RUN chmod +x $APP_DIR/app

# to run as unprivileged user
USER $APP_USER

CMD ["./app"]
