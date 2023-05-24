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

package updater

import (
	"log"
	"testing"

	"github.com/rdimitrov/go-tuf-metadata/metadata"
	"github.com/rdimitrov/go-tuf-metadata/metadata/config"
	simulator "github.com/rdimitrov/go-tuf-metadata/testutils/simulators"

	"github.com/stretchr/testify/assert"
)

var (
	metadataDir   string
	rootBytes     []byte
	updaterConfig *config.UpdaterConfig
)

func TestMain(m *testing.M) {
	err := loadTrustedRootMetadata()
	if err != nil {
		simulator.Cleanup(metadataDir)
		log.Fatalf("failed to load TrustedRootMetadata: %v\n", err)
	}

	defer simulator.Cleanup(metadataDir)
	m.Run()
}

func loadTrustedRootMetadata() error {
	var err error

	metadataDir, targetsDir, err := simulator.InitMetadataDir()
	if err != nil {
		log.Printf("failed to initialize metadata dir: %v", err)
		return err
	}

	rootBytes, err = simulator.GetRootBytes(metadataDir)
	if err != nil {
		log.Printf("failed to load root bytes: %v", err)
	}

	updaterConfig, err = config.New(metadataDir, rootBytes)
	updaterConfig.LocalMetadataDir = metadataDir
	updaterConfig.LocalTargetsDir = targetsDir
	return err
}

func TestLoadTrustedRootMetadata(t *testing.T) {
	updater, err := New(updaterConfig)

	assert.Nil(t, err)
	if assert.NotNil(t, updater) {
		assert.Equal(t, metadata.ROOT, updater.trusted.Root.Signed.Type)
		assert.Equal(t, metadata.SPECIFICATION_VERSION, updater.trusted.Root.Signed.SpecVersion)
		assert.True(t, updater.trusted.Root.Signed.ConsistentSnapshot)
		assert.Equal(t, int64(1), updater.trusted.Root.Signed.Version)
		assert.Nil(t, updater.trusted.Snapshot)
		assert.Nil(t, updater.trusted.Timestamp)
		assert.Empty(t, updater.trusted.Targets)
	}
}

func TestRefreshTopLevelMetadata(t *testing.T) {}

func TestGetTargetInfo(t *testing.T) {}

func TestDownloadTarget(t *testing.T) {}

func TestFindCachedTarget(t *testing.T) {}

func TestLoadTimestamp(t *testing.T) {}

func TestLoadSnapshot(t *testing.T) {}

func TestLoadTargets(t *testing.T) {}

func TestLoadRoot(t *testing.T) {}

func TestPreOrderDepthFirstWalk(t *testing.T) {}

func TestPersistMetadata(t *testing.T) {}

func TestDownloadMetadata(t *testing.T) {}

func TestGenerateTargetFilePath(t *testing.T) {}
