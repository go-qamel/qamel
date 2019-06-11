FROM radhifadlillah/qamel:linux as linux

# ========== END OF LINUX ========== #

FROM ubuntu:16.04 as base

# Install dependencies for Qt5
RUN apt-get -qq update && \
    apt-get -qq -y install curl build-essential libgl1-mesa-dev libfontconfig1-dev libglib2.0-dev libglu1-mesa-dev libxrender1 libdbus-1-dev libx11-dev libx11-xcb-dev

# Download source of Qt5
ENV QT_MAJOR 5.12
ENV QT_VERSION 5.12.0
RUN curl -SL --retry 10 --retry-delay 60 https://download.qt.io/official_releases/qt/$QT_MAJOR/$QT_VERSION/single/qt-everywhere-src-$QT_VERSION.tar.xz | \
    tar -xJC /

# Build Qt5 static
RUN cd qt-everywhere-src-$QT_VERSION && \
    ./configure -prefix "/opt/Qt$QT_VERSION" \
        -confirm-license -opensource -static \
        -release -optimize-size -qt-doubleconversion \
        -no-icu -qt-pcre -qt-zlib -qt-freetype \
        -qt-harfbuzz -qt-xcb -qt-libpng -qt-libjpeg \
        -make libs -nomake tools -nomake examples -nomake tests \
        -skip qtwebengine && \
    make && make install

# ========== END OF BASE ========== #

FROM ubuntu:16.04

ENV HOME /home/user
ENV GOPATH $HOME/go
ENV QT_VERSION 5.12.0
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Copy Go and Qamel from linux
COPY --from=linux /usr/local/go /usr/local/go
COPY --from=linux $GOPATH/bin $GOPATH/bin
COPY --from=linux $GOPATH/src/github.com/RadhiFadlillah/qamel $GOPATH/src/github.com/RadhiFadlillah/qamel

# Copy Qt5 from base
COPY --from=base /opt/Qt$QT_VERSION /opt/Qt$QT_VERSION

# Install dependencies for Qt5
RUN apt-get -qq update && \
    apt-get -qq -y install build-essential libgl1-mesa-dev libfontconfig1-dev libglib2.0-dev libglu1-mesa-dev libxrender1 libdbus-1-dev

# Create profile for Qamel
RUN mkdir -p $HOME/.config/qamel
RUN printf '%s %s %s %s %s %s %s %s %s %s\n' \
    '{"default": {' \
    '"OS":"linux",' \
    '"Arch":"amd64",' \
    '"Static":true,' \
    '"Qmake":"/opt/Qt5.12.0/bin/qmake",' \
    '"Moc":"/opt/Qt5.12.0/bin/moc",' \
    '"Rcc":"/opt/Qt5.12.0/bin/rcc",' \
    '"Gcc":"gcc",' \
    '"Gxx":"g++"' \
    '}}' > $HOME/.config/qamel/config.json

# Build app
ENTRYPOINT ["qamel", "build"]