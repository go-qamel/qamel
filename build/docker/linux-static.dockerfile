FROM radhifadlillah/qamel:linux as linux

# ========== END OF LINUX ========== #

FROM radhifadlillah/qamel:qt-static as base

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