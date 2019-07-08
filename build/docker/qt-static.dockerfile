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
        -release -optimize-size -qt-doubleconversion \
        -qt-zlib -qt-libpng -qt-libjpeg -qt-xcb \
        -qt-pcre -qt-freetype -qt-harfbuzz \
        -dbus-runtime -openssl-runtime -opengl \
        -sysconfdir /etc/xdg \
        -skip qtwebengine -skip qtfeedback \
        -skip qtpim -skip qtdocgallery \
        -make libs -nomake tools -nomake examples -nomake tests && \
    make -j $(grep -c ^processor /proc/cpuinfo) && \
    make install -j $(grep -c ^processor /proc/cpuinfo)