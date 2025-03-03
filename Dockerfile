ARG GOLANG_VERSION=1.17
FROM golang:${GOLANG_VERSION} as builder

#ARGS are changed during build - dev is just a placeholder
ARG RELEASE=dev
ARG COMMIT=dev
ARG LIBVIPS_VERSION=8.11.2
ARG GOLANGCILINT_VERSION=1.29.0

# Installs libvips + required libraries
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get upgrade && \
  apt-get update && \ 
  apt-get install --no-install-recommends -y \
  ca-certificates \
  automake build-essential curl \
  gobject-introspection gtk-doc-tools libglib2.0-dev libjpeg62-turbo-dev libpng-dev \
  libwebp-dev libtiff5-dev libgif-dev libexif-dev libxml2-dev libpoppler-glib-dev \
  swig libmagickwand-dev libpango1.0-dev libmatio-dev libopenslide-dev libcfitsio-dev \
<<<<<<< HEAD
  libgsf-1-dev fftw3-dev liborc-0.4-dev librsvg2-dev libopenjp2-7-dev libheif-dev \
  libimagequant-dev 

RUN cd /tmp && \
  ldconfig && \
=======
  libgsf-1-dev fftw3-dev liborc-0.4-dev librsvg2-dev libimagequant-dev libheif-dev && \
  cd /tmp && \
>>>>>>> 826d2f4958de2a7862fe92b8e4afde1f206c56f6
  curl -fsSLO https://github.com/libvips/libvips/releases/download/v${LIBVIPS_VERSION}/vips-${LIBVIPS_VERSION}.tar.gz && \
  tar zvxf vips-${LIBVIPS_VERSION}.tar.gz && \
  cd /tmp/vips-${LIBVIPS_VERSION} && \
	CFLAGS="-g -O3" CXXFLAGS="-D_GLIBCXX_USE_CXX11_ABI=0 -g -O3" \
    ./configure \
    --disable-debug \
    --disable-dependency-tracking \
    --disable-introspection \
    --disable-static \
    --enable-gtk-doc-html=no \
    --enable-gtk-doc=no \
    --with-openslide && \
  make && \
  make install && \
  ldconfig

<<<<<<< HEAD
=======
# Installing golangci-lint
WORKDIR /tmp
RUN curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH}/bin" v${GOLANGCILINT_VERSION}
>>>>>>> 826d2f4958de2a7862fe92b8e4afde1f206c56f6

WORKDIR ${GOPATH}/src/github.com/h2non/imaginary

# Cache go modules
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
<<<<<<< HEAD

RUN go mod download

# Copy imaginary sources
COPY . .
=======

RUN go mod download

# Copy imaginary sources
COPY . .

# Run quality control
RUN go test -test.v -test.race -test.covermode=atomic .
RUN golangci-lint run .
>>>>>>> 826d2f4958de2a7862fe92b8e4afde1f206c56f6

# Compile imaginary
RUN go build -a \
    -o ${GOPATH}/bin/imaginary \
    -ldflags="-s -w -h -X main.Version=${RELEASE}-${COMMIT}" \
    github.com/h2non/imaginary

<<<<<<< HEAD
FROM debian:stable-slim
=======
FROM debian:buster-slim
>>>>>>> 826d2f4958de2a7862fe92b8e4afde1f206c56f6

ARG IMAGINARY_VERSION

LABEL maintainer="tomas@aparicio.me" \
      org.label-schema.description="Fast, simple, scalable HTTP microservice for high-level image processing with first-class Docker support" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://github.com/h2non/imaginary" \
      org.label-schema.vcs-url="https://github.com/h2non/imaginary" \
      org.label-schema.version="${RELEASE}-${COMMIT}"

COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /go/bin/imaginary /usr/local/bin/imaginary
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /usr/local/bin/vips /usr/local/bin/vips

# Install runtime dependencies
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get upgrade && \ 
  apt-get update && \
  apt-get install --no-install-recommends -y \
<<<<<<< HEAD
  libglib2.0-0 libjpeg62-turbo libpng16-16 libopenexr25 \
  libwebp6 libwebpmux3 libwebpdemux2 libtiff5 libgif7 libexif12 libxml2 libpoppler-glib8 \
  libmagickwand-6.q16-6 libpango1.0-0 libmatio11 libopenslide0 \
  libgsf-1-114 fftw3 liborc-0.4-0 librsvg2-2 libcfitsio9 libopenjp2-7 libheif1 \
  libimagequant0 && \
=======
  procps libglib2.0-0 libjpeg62-turbo libpng16-16 libopenexr23 \
  libwebp6 libwebpmux3 libwebpdemux2 libtiff5 libgif7 libexif12 libxml2 libpoppler-glib8 \
  libmagickwand-6.q16-6 libpango1.0-0 libmatio4 libopenslide0 \
  libgsf-1-114 fftw3 liborc-0.4-0 librsvg2-2 libcfitsio7 libimagequant0 libheif1 && \
>>>>>>> 826d2f4958de2a7862fe92b8e4afde1f206c56f6
  apt-get autoremove -y && \
  apt-get autoclean && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Server port to listen
ENV PORT 9000

# Drop privileges for non-UID mapped environments
USER nobody

# Run the entrypoint command by default when the container starts.
ENTRYPOINT ["/usr/local/bin/imaginary"]

# Expose the server TCP port
EXPOSE ${PORT}
