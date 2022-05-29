FROM ubuntu:focal

LABEL maintainer="9tmark, github.com/9tmark"

ARG DEBIAN_FRONTEND=noninteractive

ENV PATH=/usr/local/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/opt/avly-trader/bin \
    DISPLAY=:1 \
    VNC_PORT=5900 \
    USER=root \
    AVL_LOGS=/var/log/avly-trader \
    THIRD_PARTY=/opt/third-party \
    TERMINAL=xterm

RUN set -ex; \
    dpkg --add-architecture i386; \
    apt-get update -yq; \
    apt-get install -yq \
        apt-transport-https \
        apt-utils \
        bzip2 \
        binutils \
        cabextract \
        curl \
        dbus-x11 \
        git \
        gnupg2 \
        i3 \
        locales \
        net-tools \
        procps \
        psmisc \
        samba \
        software-properties-common \
        ssh \
        supervisor \
        tar \
        vim \
        wget \
        winbind \
        x11vnc \
        xauth \
        xdotool \
        xfonts-base \
        xinit \
        xorg \
        xterm \
        xvfb \
        xz-utils; \
    apt-get purge -yq \
        dunst \
        i3lock \
        pm-utils \
        suckless-tools \
        xscreensaver*; \
    apt-get clean -yq; \
    apt-get autoremove -yq

RUN set -ex; \
    mkdir -p /opt/avly-trader/bin; \
    mkdir -p /root/.config/i3; \
    mkdir -p ${AVL_LOGS}; \
    mkdir -p ${THIRD_PARTY}

ADD build/avly /opt/avly-trader/bin
ADD resources/01-build/config /etc/i3

RUN set -ex; \
    chmod -R +x /opt/avly-trader/bin; \
    chmod -R +rw ${AVL_LOGS}; \
    chmod -R +rx ${THIRD_PARTY}

EXPOSE ${VNC_PORT}

CMD [ "/opt/avly-trader/bin/avly", "-enter" ]
