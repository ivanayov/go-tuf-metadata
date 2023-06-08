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
	"log"
	"os"
)

var (
	MetadataURL = "https://jku.github.io/tuf-demo/metadata"
	TargetsURL  = "https://jku.github.io/tuf-demo/targets"
)

var (
	// downloadPath = "/download"
	metadataPath = "/metadata"
	targetsPath  = "/targets"
	LocalDir     string
	DumpDir      string
)

func InitLocalEnv() error {

	tmp := fs.TempDir()

	tmpDir, err := fs.MkdirTemp(tmp, "0750")
	if err != nil {
		log.Fatal("failed to create temporary directory: ", err)
	}

	fs.Mkdir(tmpDir + metadataPath)
	fs.Mkdir(tmpDir + targetsPath)
	LocalDir = tmpDir
	return nil
}

func InitMetadataDir() (*RepositorySimulator, string, string, error) {
	err := InitLocalEnv()
	if err != nil {
		log.Fatal("failed to initialize environment: ", err)
	}
	metadataDir := LocalDir + metadataPath

	sim := NewRepository()

	f, err := os.Create(metadataDir + "/root.json")
	if err != nil {
		log.Fatalf("failed to create root: %v", err)
	}

	f.Write(sim.SignedRoots[0])

	// if len(DumpDir) > 0 {

	// }
	targetsDir := LocalDir + targetsPath
	return sim, metadataDir, targetsDir, err
}

func GetRootBytes(localMetadataDir string) ([]byte, error) {
	return fs.ReadFile(localMetadataDir + "/root.json")
}

func Cleanup(tmpDir string) {
	log.Printf("Cleaning temporary directory: %s\n", tmpDir)
	fs.RemoveAll(tmpDir)
}
