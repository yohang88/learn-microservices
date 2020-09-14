#!/bin/bash

protoc proto/review.proto --go_out=plugins=grpc:.