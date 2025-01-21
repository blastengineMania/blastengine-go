#!/bin/bash

env $(grep -v '^#' .env | xargs) go test
