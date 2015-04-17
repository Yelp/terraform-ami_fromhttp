FROM docker-dev.yelpcorp.com/lucid_yelp
MAINTAINER Keith Mosher <kmosher@yelp.com>

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
    go \
    git \
    build-essential \
    ruby1.8 rubygems \
    libopenssl-ruby \
    ruby-dev \
    --no-install-recommends

ENV GOPATH /go

RUN gem install fpm
RUN ln -s /var/lib/gems/1.8/bin/fpm /usr/local/bin/fpm

RUN git clone --branch 0.4.1-pre2 https://github.com/bobtfish/terraform.git && \
    mkdir -p /go/src/github.com/hashicorp && mv terraform /go/src/github.com/hashicorp

#WORKDIR /go/src/terraform-provider-yelpaws
