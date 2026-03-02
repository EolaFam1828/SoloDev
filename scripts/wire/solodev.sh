#!/usr/bin/env sh
# Copyright 2026 EolaFam1828. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

echo "Updating cmd/solodev/wire_gen.go"
go run github.com/google/wire/cmd/wire gen github.com/EolaFam1828/SoloDev/cmd/solodev
# format generated file as we can't exclude it from being formatted easily.
goimports -w ./cmd/solodev/wire_gen.go
