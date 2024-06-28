FROM scratch
COPY wgrep /wgrep
ENTRYPOINT ["/wgrep"]

