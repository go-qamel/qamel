FROM radhifadlillah/qamel:linux as linux

# ========== END OF LINUX ========== #

FROM radhifadlillah/qamel:qt-static as base

# ========== END OF BASE ========== #

FROM ubuntu:16.04

ENV HOME /home/user
ENV GOPATH $HOME/go
ENV QT_VERSION 5.13.2
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Install ca-certificates which might be needed by Go proxy
RUN apt-get -qq update && \
    apt-get -qq -y install ca-certificates git

# Install dependencies for Qt5
RUN apt-get -qq update && \
    apt-get -qq -y install python libgl1-mesa-dev \
    libfontconfig1-dev libfreetype6-dev libx11-dev libxext-dev \
    libxfixes-dev libxi-dev libxrender-dev libxcb1-dev \
    libx11-xcb-dev libxcb-glx0-dev libxkbcommon-dev \
    libxkbcommon-x11-dev '^libxcb.*-dev' && \
    apt-get -qq clean

# Install ccache for faster build
RUN apt-get -qq update && \
    apt-get -qq -y install ccache && \
    apt-get -qq clean
ENV PATH "/usr/lib/ccache:$PATH"

# Copy Qt5 from base
COPY --from=base /opt/Qt$QT_VERSION /opt/Qt$QT_VERSION

# Copy Go and Qamel from linux
COPY --from=linux /usr/local/go /usr/local/go
COPY --from=linux $GOPATH/bin $GOPATH/bin
COPY --from=linux $GOPATH/src/github.com/go-qamel/qamel $GOPATH/src/github.com/go-qamel/qamel

# Create profile for Qamel
RUN mkdir -p $HOME/.config/qamel
RUN printf '%s %s %s %s %s %s %s %s %s %s\n' \
    '{"default": {' \
    '"OS":"linux",' \
    '"Arch":"amd64",' \
    '"Static":true,' \
    '"Qmake":"/opt/Qt5.13.2/bin/qmake",' \
    '"Moc":"/opt/Qt5.13.2/bin/moc",' \
    '"Rcc":"/opt/Qt5.13.2/bin/rcc",' \
    '"Gcc":"gcc",' \
    '"Gxx":"g++"' \
    '}}' > $HOME/.config/qamel/config.json

# Build app
ENTRYPOINT ["qamel", "build"]