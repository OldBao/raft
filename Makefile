all: iface main

iface: proto/vote.pb.go

.PHONY: all iface

proto/vote.pb.go : interfaces/vote.proto
	cd interfaces &&\
	protoc --go_out=plugins=grpc:../proto vote.proto

main:
	go build
