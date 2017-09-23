FROM alpine:3.4

ENV CHEZ_VERSION 9.4

RUN apk add --no-cache ncurses libx11 \
 && apk add --no-cache --virtual .build-deps build-base openssl curl ncurses-dev libx11-dev \
 && wget https://github.com/cisco/ChezScheme/archive/v$CHEZ_VERSION.tar.gz  \
 && tar -xf v$CHEZ_VERSION.tar.gz \
 && cd /ChezScheme-$CHEZ_VERSION \
 && ln -s /usr/include/locale.h /usr/include/xlocale.h \
 && ./configure \
 && make install \
 && cd / \
 && rm -rf /ChezScheme-$CHEZ_VERSION \
 && rm v$CHEZ_VERSION.tar.gz \
 && apk del .build-deps

WORKDIR /data
VOLUME ["/data"]

CMD [ "scheme" ]