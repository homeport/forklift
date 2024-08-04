FROM scratch
COPY base-layer /boot

COPY update /etc

COPY run-0 /usr/local/bin
ENV FOO=BAR

COPY run-1 /opt/tool
ENV BAR=FOO
