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
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/rdimitrov/go-tuf-metadata/metadata"
	"github.com/rdimitrov/go-tuf-metadata/metadata/config"
	simulator "github.com/rdimitrov/go-tuf-metadata/testutils/simulators"
	"github.com/sigstore/sigstore/pkg/signature"

	"github.com/stretchr/testify/assert"
)

var (
	metadataDir   string
	targetsDir    string
	rootBytes     []byte
	pastDateTime  time.Time
	updaterConfig *config.UpdaterConfig
	// TODO: rename to repository sim
	sim *simulator.RepositorySimulator
)

func TestMain(m *testing.M) {
	err := loadTrustedRootMetadata()
	pastDateTime = time.Now().UTC().Truncate(24 * time.Hour).Add(-5 * 24 * time.Hour)

	if err != nil {
		simulator.Cleanup(metadataDir)
		log.Fatalf("failed to load TrustedRootMetadata: %v\n", err)
	}

	defer simulator.Cleanup(metadataDir)
	m.Run()
}

func loadTrustedRootMetadata() error {
	var err error

	sim, metadataDir, targetsDir, err = simulator.InitMetadataDir()
	if err != nil {
		log.Printf("failed to initialize metadata dir: %v", err)
		return err
	}

	rootBytes, err = simulator.GetRootBytes(metadataDir)
	if err != nil {
		log.Printf("failed to load root bytes: %v", err)
	}

	updaterConfig, err = config.New(metadataDir, rootBytes)
	updaterConfig.Fetcher = sim
	updaterConfig.LocalMetadataDir = metadataDir
	updaterConfig.LocalTargetsDir = targetsDir
	return err
}

func initDumpDir(t *testing.T) {
	if len(sim.DumpDir) == 0 {
		// create test specific dump directory
		sim.DumpDir = t.TempDir()
	}
}

// runRefresh creates new Updater instance and
// runs Refresh
func runRefresh() Updater {
	if len(sim.DumpDir) > 0 {
		sim.Write()
	}

	updater, err := New(updaterConfig)
	if err != nil {
		//TODO
	}
	updater.Refresh()
	return *updater
}

func initUpdater() Updater {
	if len(sim.DumpDir) > 0 {
		sim.Write()
	}

	updater, err := New(updaterConfig)
	if err != nil {
		//TODO
	}
	return *updater
}

// Asserts that local metadata files exist for 'roles'
func assertFilesExist(t *testing.T, roles []string) {
	expectedFiles := []string{}

	for _, role := range roles {
		expectedFiles = append(expectedFiles, fmt.Sprintf("%s.json", role))
	}
	localMetadataFiles, err := os.ReadDir(metadataDir)
	if err != nil {
		// TODO
	}
	actual := []string{}
	for _, file := range localMetadataFiles {
		actual = append(actual, file.Name())
	}

	for _, file := range expectedFiles {
		assert.Contains(t, actual, file)
	}

}

// Asserts that local file content is the expected
func assertContentEquals(t *testing.T, role string, version *int) {
	expectedContent, err := sim.FetchMetadata(role, version)
	if err != nil {
		// TODO
	}

	content, err := os.ReadFile(fmt.Sprintf("%s/%s.json", metadataDir, role))
	if err != nil {
		// TODO
	}
	assert.Equal(t, expectedContent, content)
}

func assertVersionEquals(t *testing.T, role string, expectedVersion int64) {
	path := fmt.Sprintf("%s/%s.json", metadataDir, role)
	md, err := sim.MDRoot.FromFile(path)
	if err != nil {
		// TODO
	}
	assert.Equal(t, md.Signed.Version, expectedVersion)
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

func TestFirstTimeRefresh(t *testing.T) {
	assertFilesExist(t, []string{metadata.ROOT})
	sim.MDRoot.Signed.Version += 1
	sim.PublishRoot()

	runRefresh()
	assertFilesExist(t, metadata.TOP_LEVEL_ROLE_NAMES[:])

	for _, role := range metadata.TOP_LEVEL_ROLE_NAMES {
		version := -1
		if role == metadata.ROOT {
			version = 2
		}
		// TODO: add ignoreVersion flag instead of -1
		assertContentEquals(t, role, &version)
	}
}

func TestTrustedRootMissing(t *testing.T) {
	// TODO: reset config
	// TODO: check for bugs with global metadata and targets dir - should start clean
	os.Remove(fmt.Sprintf("%s/%s.json", metadataDir, metadata.ROOT))
	tmp := metadataDir

	// TODO: should have failed initializing updater
	// for missing local root
	runRefresh()

	res, err := os.ReadDir(metadataDir)
	fmt.Println(tmp)
	assert.Nil(t, err)
	// assert.Equal(t, 1, 2)
	for _, f := range res {
		fmt.Printf("res: %s\n", f.Name())
	}
}

func TestTrustedRootExpired(t *testing.T) {
	sim.MDRoot.Signed.Expires = pastDateTime
	sim.MDRoot.Signed.Version += 1
	sim.PublishRoot()

	updater := initUpdater()
	// TODO: should fail with expired metadata error
	updater.Refresh()

	assertFilesExist(t, []string{metadata.ROOT})
	version := 2
	assertContentEquals(t, metadata.ROOT, &version)

	updater = initUpdater()

	sim.MDRoot.Signed.Expires = sim.SafeExpiry
	sim.MDRoot.Signed.Version += 1
	sim.PublishRoot()
	updater.Refresh()

	assertFilesExist(t, metadata.TOP_LEVEL_ROLE_NAMES[:])
	version = 3
	assertContentEquals(t, metadata.ROOT, &version)
}

func TestTrustedRootUnsigned(t *testing.T) {
	rootPath := fmt.Sprintf("%s/%s.json", metadataDir, metadata.ROOT)
	mdRoot, err := sim.MDRoot.FromFile(rootPath)
	if err != nil {
		// TODO
	}
	mdRoot.ClearSignatures()
	mdRoot.ToFile(rootPath, true)
	// TODO: should fail with unsigned metadata error
	runRefresh()

	assertFilesExist(t, []string{metadata.ROOT})
	mdRootAfter, err := sim.MDRoot.FromFile(rootPath)
	if err != nil {
		// TODO
	}
	expected, err := mdRoot.ToBytes(false)
	if err != nil {
		// TODO
	}
	actual, err := mdRootAfter.ToBytes(false)
	if err != nil {
		// TODO
	}

	assert.Equal(t, expected, actual)
}

func TestMaxRootRotations(t *testing.T) {
	updater := initUpdater()
	updater.cfg.MaxRootRotations = 3

	for sim.MDRoot.Signed.Version < updater.cfg.MaxRootRotations+3 {
		sim.MDRoot.Signed.Version += 1
		sim.PublishRoot()
	}

	rootPath := fmt.Sprintf("%s/%s.json", metadataDir, metadata.ROOT)
	mdRoot, err := sim.MDRoot.FromFile(rootPath)
	if err != nil {
		//TODO
	}
	initialRootVersion := mdRoot.Signed.Version

	updater.Refresh()

	assertVersionEquals(t, metadata.ROOT, initialRootVersion+updaterConfig.MaxRootRotations)
}

func TestIntermediateRootInclorrectlySigned(t *testing.T) {
	sim.MDRoot.Signed.Version += 1
	rootSigners := make(map[string]*signature.Signer)
	for k, v := range sim.Signers[metadata.ROOT] {
		rootSigners[k] = v
	}
	for k := range sim.Signers[metadata.ROOT] {
		delete(sim.Signers[metadata.ROOT], k)
	}
	sim.PublishRoot()

	//TODO: should fail for unsigned metadata
	runRefresh()

	assertFilesExist(t, []string{metadata.ROOT})
	version := 1
	assertContentEquals(t, metadata.ROOT, &version)
}

// func TestGetTargetInfo(t *testing.T) {

// 	sim.AddTarget("targets", []byte("content"), "targetpath")
// 	sim.MDTargets.Signed.Version += 1
// 	// sim.UpdateSnapshot()

// 	updater, err := New(updaterConfig)
// 	assert.Nil(t, err)
// 	targetInfo, err := updater.GetTargetInfo("targetpath")
// 	assert.Nil(t, err)
// 	assert.Nil(t, targetInfo)
// }

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
