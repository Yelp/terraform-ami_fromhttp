.PHONY: all fmt clean package test itest_%

all: fmt test

fmt:
	go fmt terraform-provider-yelpaws/yelpaws
	go fmt terraform-provider-yelpaws

clean:
	make -C yelppack clean

itest_%:
	make -C yelppack $@

package: itest_lucid

test:
	go test terraform-provider-yelpaws/yelpaws
	go test terraform-provider-yelpaws
