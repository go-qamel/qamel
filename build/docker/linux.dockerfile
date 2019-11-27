FROM ubuntu:16.04 as base

ENV HOME /home/user
ENV GOPATH $HOME/go
ENV GO_VERSION 1.13.4
ENV QT_MAJOR 5.13
ENV QT_VERSION 5.13.2

# Install Go
RUN apt-get -qq update && \
    apt-get -qq -y install ca-certificates curl git
RUN curl -SL --retry 10 --retry-delay 60 https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz | \
    tar -xzC /usr/local

# Install Qt5 dependencies
RUN apt-get -qq update && \
    apt-get -qq -y install dbus libfontconfig1 libx11-6 libx11-xcb1

# Download Qt5
RUN curl -SL --retry 10 --retry-delay 60 -O \
    https://download.qt.io/official_releases/qt/$QT_MAJOR/$QT_VERSION/qt-opensource-linux-x64-$QT_VERSION.run

# Download Qt5 installation script
RUN curl -SL --retry 10 --retry-delay 60 -O \
    https://raw.githubusercontent.com/RadhiFadlillah/qamel/master/build/docker/installer-script.qs #c35715

# Install Qt5
RUN chmod +x qt-opensource-linux-x64-$QT_VERSION.run && \
    ./qt-opensource-linux-x64-$QT_VERSION.run -v \
        --script installer-script.qs \
        --platform minimal

# Clean up after installing Qt5
RUN rm -Rf /opt/Qt$QT_VERSION/Docs \
            /opt/Qt$QT_VERSION/Examples \
            /opt/Qt$QT_VERSION/Tools

# Install Qamel
RUN /usr/local/go/bin/go get -u github.com/go-qamel/qamel/cmd/qamel #9729573

# ========== END OF BASE ========== #

FROM ubuntu:16.04

ENV HOME /home/user
ENV GOPATH $HOME/go
ENV QT_VERSION 5.13.2
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Install dependencies for Qt5
RUN apt-get -qq update && \
    apt-get -qq -y install build-essential libgl1-mesa-dev \
        libfontconfig1-dev libfreetype6-dev libx11-dev libxext-dev \
        libxfixes-dev libxi-dev libxrender-dev libxcb1-dev \
        libx11-xcb-dev libxcb-glx0-dev libxkbcommon-x11-dev && \
    apt-get -qq clean

# Install ccache for faster build
RUN apt-get -qq update && \
    apt-get -qq -y install ccache && \
    apt-get -qq clean
ENV PATH "/usr/lib/ccache:$PATH"

# Copy Go and Qamel from base
COPY --from=base /usr/local/go /usr/local/go
COPY --from=base $GOPATH/bin $GOPATH/bin
COPY --from=base $GOPATH/src/github.com/go-qamel/qamel $GOPATH/src/github.com/go-qamel/qamel

# Copy Qt5 from base
COPY --from=base /opt/Qt$QT_VERSION /opt/Qt$QT_VERSION

# Create profile for Qamel
RUN mkdir -p $HOME/.config/qamel
RUN printf '%s %s %s %s %s %s %s %s %s %s\n' \
    '{"default": {' \
    '"OS":"linux",' \
    '"Arch":"amd64",' \
    '"Static":false,' \
    '"Qmake":"/opt/Qt5.13.0/5.13.0/gcc_64/bin/qmake",' \
    '"Moc":"/opt/Qt5.13.0/5.13.0/gcc_64/bin/moc",' \
    '"Rcc":"/opt/Qt5.13.0/5.13.0/gcc_64/bin/rcc",' \
    '"Gcc":"gcc",' \
    '"Gxx":"g++"' \
    '}}' > $HOME/.config/qamel/config.json

# Build app
ENTRYPOINT ["qamel", "build"]