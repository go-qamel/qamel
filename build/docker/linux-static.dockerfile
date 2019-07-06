FROM radhifadlillah/qamel:linux as linux

# ========== END OF LINUX ========== #

FROM ubuntu:18.04 as base

# Install dependencies for Qt5
RUN apt-get -qq update && \
    apt-get -qq -y install curl build-essential python \
    libglib2.0-dev libglu1-mesa-dev libpulse-dev fontconfig \
    libasound2 libegl1-mesa libnss3 libpci3 libxcomposite1 \
    libxcursor1 libxi6 libxrandr2 libxtst6 libdbus-1-dev libssl-dev \
    libxkbcommon-dev libfontconfig1-dev libfreetype6-dev libx11-dev \
    libxext-dev libxfixes-dev libxi-dev libxrender-dev libxcb1-dev \
    libx11-xcb-dev libxcb-glx0-dev libxcb-keysyms1-dev \
    libxcb-image0-dev libxcb-shm0-dev libxcb-icccm4-dev \
    libxcb-sync0-dev libxcb-xfixes0-dev libxcb-shape0-dev \
    libxcb-randr0-dev libxcb-render-util0-dev

# Download source of Qt5
ENV QT_MAJOR 5.13
ENV QT_VERSION 5.13.0
RUN curl -SL --retry 10 --retry-delay 60 https://download.qt.io/official_releases/qt/$QT_MAJOR/$QT_VERSION/single/qt-everywhere-src-$QT_VERSION.tar.xz | \
    tar -xJC /

# Build Qt5 static
RUN cd qt-everywhere-src-$QT_VERSION && \
    ./configure -prefix "/opt/Qt$QT_VERSION" \
        -confirm-license -opensource -static \
        -qt-zlib -qt-libpng -qt-libjpeg -qt-xcb \
        -sysconfdir /etc/xdg -dbus-runtime -openssl-runtime \
        -opengl -optimize-size -skip qtwebengine -skip qtfeedback \
        -skip qtpim -skip qtdocgallery -skip qtwebengine \
        -nomake tests -nomake examples && \
    make -j $(grep -c ^processor /proc/cpuinfo) && \
    make install -j $(grep -c ^processor /proc/cpuinfo)

# ========== END OF BASE ========== #

FROM ubuntu:18.04

ENV HOME /home/user
ENV GOPATH $HOME/go
ENV QT_VERSION 5.13.0
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Copy Go and Qamel from linux
COPY --from=linux /usr/local/go /usr/local/go
COPY --from=linux $GOPATH/bin $GOPATH/bin
COPY --from=linux $GOPATH/src/github.com/RadhiFadlillah/qamel $GOPATH/src/github.com/RadhiFadlillah/qamel

# Copy Qt5 from base
COPY --from=base /opt/Qt$QT_VERSION /opt/Qt$QT_VERSION

# Install dependencies for Qt5
RUN apt-get -qq update && \
    apt-get -qq -y install build-essential libglib2.0-dev libglu1-mesa-dev libpulse-dev \
        fontconfig libasound2 libegl1-mesa libnss3 libpci3 libxcomposite1 libxcursor1 \
        libxi6 libxrandr2 libxtst6 libfontconfig1-dev libfreetype6-dev libxrender-dev \
        libxkbcommon-dev && \
    apt-get -qq clean

# Create profile for Qamel
RUN mkdir -p $HOME/.config/qamel
RUN printf '%s %s %s %s %s %s %s %s %s %s\n' \
    '{"default": {' \
    '"OS":"linux",' \
    '"Arch":"amd64",' \
    '"Static":true,' \
    '"Qmake":"/opt/Qt5.13.0/bin/qmake",' \
    '"Moc":"/opt/Qt5.13.0/bin/moc",' \
    '"Rcc":"/opt/Qt5.13.0/bin/rcc",' \
    '"Gcc":"gcc",' \
    '"Gxx":"g++"' \
    '}}' > $HOME/.config/qamel/config.json

# Build app
ENTRYPOINT ["qamel", "build"]