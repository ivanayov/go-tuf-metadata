// Copyright 2023 VMware, Inc.
//
// This product is licensed to you under the BSD-2 license (the "License").
// You may not use this product except in compliance with the BSD-2 License.
// This product may include a number of subcomponents with separate copyright
// notices and license terms. Your use of these subcomponents is subject to
// the terms and conditions of the subcomponent's license, as noted in the
// LICENSE file.
//
// SPDX-License-Identifier: BSD-2-Clause

package simulator

import (
	"os"
)

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (*osFS) MkdirTemp(dir string, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}

func (*osFS) TempDir() string {
	return os.TempDir()
}

func (*osFS) Mkdir(name string) error {
	return os.Mkdir(name, 0750)
}

func (*osFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (*osFS) WriteFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0644)
}

func (*osFS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

var fs fileSystem = &osFS{}

type fileSystem interface {
	TempDir() string
	Mkdir(name string) error
	MkdirTemp(dir string, pattern string) (string, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte) error
	RemoveAll(path string) error
}
