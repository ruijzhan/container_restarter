FROM alpine:3.11
COPY restarter /restarter
ENTRYPOINT ["/restarter"]
