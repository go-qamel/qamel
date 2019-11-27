FROM ubuntu:16.04 as base

# Download source of Qt5
RUN apt-get -qq update && \
    apt-get -qq -y install curl build-essential 

ENV QT_MAJOR 5.13
ENV QT_VERSION 5.13.2
RUN curl -SL --retry 10 --retry-delay 60 https://download.qt.io/official_releases/qt/$QT_MAJOR/$QT_VERSION/single/qt-everywhere-src-$QT_VERSION.tar.xz | \
    tar -xJC /

# Install dependencies for Qt5
RUN apt-get -qq update && \
    apt-get -qq -y install python libgl1-mesa-dev \
    libfontconfig1-dev libfreetype6-dev libx11-dev libxext-dev \
    libxfixes-dev libxi-dev libxrender-dev libxcb1-dev \
    libx11-xcb-dev libxcb-glx0-dev libxkbcommon-dev \
    libxkbcommon-x11-dev '^libxcb.*-dev'

# Build Qt5 static
RUN cd qt-everywhere-src-$QT_VERSION && \
    ./configure -static -prefix "/opt/Qt$QT_VERSION" \
        -opensource -confirm-license -release \
        -optimize-size -strip -fontconfig \
        -qt-zlib -qt-libjpeg -qt-libpng -qt-xcb \
        -qt-pcre -qt-harfbuzz -qt-doubleconversion \
        -nomake tools -nomake examples -nomake tests \
        -no-pch -skip qtwebengine && \
    make -j $(grep -c ^processor /proc/cpuinfo) && \
    make install -j $(grep -c ^processor /proc/cpuinfo)