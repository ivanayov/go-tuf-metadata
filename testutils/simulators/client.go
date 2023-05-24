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

var downloadPath = "/download"
var metadataPath = "/metadata"
var targetsPath = "/targets"

func InitLocalEnv() (string, error) {

	tmp := fs.TempDir()

	tmpDir, err := fs.MkdirTemp(tmp, "0750")
	if err != nil {
		log.Fatal("failed to create temporary directory: ", err)
	}

	// create a destination folder for storing the downloaded target
	fs.Mkdir(tmpDir + downloadPath)
	fs.Mkdir(tmpDir + metadataPath)
	fs.Mkdir(tmpDir + targetsPath)
	return tmpDir, nil
}

func InitMetadataDir() (string, string, error) {
	localDir, err := InitLocalEnv()
	if err != nil {
		log.Fatal("failed to initialize environment: ", err)
	}
	metadataDir := localDir + metadataPath

	sim := NewRepository()

	f, err := os.Create(metadataDir + "/root.json")
	if err != nil {
		log.Fatalf("failed to create root: %v", err)
	}

	f.Write(sim.SignedRoots[0])
	targetsDir := localDir + targetsPath
	return metadataDir, targetsDir, err
}

func GetRootBytes(localMetadataDir string) ([]byte, error) {
	return fs.ReadFile(localMetadataDir + "/root.json")
}

func Cleanup(tmpDir string) {
	log.Printf("Cleaning temporary directory: %s\n", tmpDir)
	fs.RemoveAll(tmpDir)
}
