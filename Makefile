.PHONY: release all

all: release

release:
	mkdir dist
	go build
	mv hamgo dist/hamgo.x86_64
	GOARCH=mipsle go build
	mv hamgo dist/hamgo.mipsle
	GOARCH=arm go build
	mv hamgo dist/hamgo.arm
