// Copyright (c) 2018, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package types

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sylabs/singularity/internal/pkg/sylog"
)

// Bundle is the temporary build environment used during the image
// building process. A Bundle is the programmatic representation of
// the directory structure which will constitute this environmenb.
// /tmp/...:
//     fs/ - A chroot filesystem
//     .singularity.d/ - Container metadata (from 2.x image format)
//     config.json (optional) - Contain information for OCI image bundle
//     etc... - The Bundle dir can theoretically contain arbitrary directories,
//              files, etc... which can be interpreted by the Chef
type Bundle struct {
	// FSObjects is a map of the filesystem objects contained in the Bundle. An object
	// will be built as one section of a SIF file.
	//
	// Known FSObjects labels:
	//   * rootfs -> root file system
	//   * .singularity.d -> .singularity.d directory (includes image exec scripts)
	//   * data -> directory containing data files
	FSObjects   map[string]string `json:"fsObjects"`
	JSONObjects map[string][]byte `json:"jsonObjects"`
	Recipe      Definition        `json:"rawDeffile"`
	BindPath    []string          `json:"bindPath"`
	Path        string            `json:"bundlePath"`
	Force       bool              `json:"force"`
	Update      bool              `json:"update"`
	NoTest      bool              `json:"noTest"`
	Sections    []string          `json:"sections"`
}

// NewBundle creates a Bundle environment
func NewBundle(directoryPrefix string) (b *Bundle, err error) {
	b = &Bundle{}

	if directoryPrefix == "" {
		directoryPrefix = "sbuild-"
	}

	b.Path, err = ioutil.TempDir("", directoryPrefix+"-")
	if err != nil {
		return nil, err
	}
	sylog.Debugf("Created temporary directory for bundle %v\n", b.Path)

	b.FSObjects = map[string]string{
		"rootfs": "fs",
	}

	for _, fso := range b.FSObjects {
		if err = os.MkdirAll(filepath.Join(b.Path, fso), 0755); err != nil {
			return
		}
	}

	return b, nil
}

// Rootfs give the path to the root filesystem in the Bundle
func (b *Bundle) Rootfs() string {
	return filepath.Join(b.Path, b.FSObjects["rootfs"])
}

// RunSection iterates through the sections specified in a bundle
// and returns true if the given string, s, is a section of the
// definition that should be executed during the build process
func (b Bundle) RunSection(s string) bool {
	for _, section := range b.Sections {
		if section == "none" {
			return false
		}
		if section == "all" || section == s {
			return true
		}
	}
	return false
}