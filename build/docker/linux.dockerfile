FROM ubuntu:18.04 as base

ENV HOME /home/user
ENV GOPATH $HOME/go
ENV GO_VERSION 1.12.6
ENV QT_MAJOR 5.13
ENV QT_VERSION 5.13.0

# Install Go
RUN apt-get -qq update && \
    apt-get -qq -y install ca-certificates curl git
RUN curl -SL --retry 10 --retry-delay 60 https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz | \
    tar -xzC /usr/local

# Install Qamel
RUN /usr/local/go/bin/go get -u github.com/RadhiFadlillah/qamel/cmd/qamel

# Install Qt5
RUN apt-get -qq update && \
    apt-get -qq -y install dbus libfontconfig1 libx11-6 libx11-xcb1
RUN curl -SL --retry 10 --retry-delay 60 -O https://download.qt.io/official_releases/qt/$QT_MAJOR/$QT_VERSION/qt-opensource-linux-x64-$QT_VERSION.run

RUN chmod +x qt-opensource-linux-x64-$QT_VERSION.run && \
    ./qt-opensource-linux-x64-$QT_VERSION.run -v \
        --script $GOPATH/src/github.com/RadhiFadlillah/qamel/build/docker/installer-script.qs \
        --platform minimal

# Clean up after installing Qt5
RUN rm -Rf /opt/Qt$QT_VERSION/Docs \
           /opt/Qt$QT_VERSION/Examples \
           /opt/Qt$QT_VERSION/Tools

# ========== END OF BASE ========== #

FROM ubuntu:18.04

ENV HOME /home/user
ENV GOPATH $HOME/go
ENV QT_VERSION 5.13.0
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Copy Go and Qamel from base
COPY --from=base /usr/local/go /usr/local/go
COPY --from=base $GOPATH/bin $GOPATH/bin
COPY --from=base $GOPATH/src/github.com/RadhiFadlillah/qamel $GOPATH/src/github.com/RadhiFadlillah/qamel

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
    '"Static":false,' \
    '"Qmake":"/opt/Qt5.13.0/5.13.0/gcc_64/bin/qmake",' \
    '"Moc":"/opt/Qt5.13.0/5.13.0/gcc_64/bin/moc",' \
    '"Rcc":"/opt/Qt5.13.0/5.13.0/gcc_64/bin/rcc",' \
    '"Gcc":"gcc",' \
    '"Gxx":"g++"' \
    '}}' > $HOME/.config/qamel/config.json

# Build app
ENTRYPOINT ["qamel", "build"]