#!/usr/bin/env bash
protoc -I proto/ proto/supervisor.proto --go_out=plugins=grpc:generated