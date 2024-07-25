// Copyright 2022 ROC. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

func PrintHelloBanner(text string) {
	fmt.Println(text)
}

func SidStr(name string, version *semver.Version, size int) string {
	return fmt.Sprintf(fmt.Sprintf("%%s@%%-%ds", size-len(name+version.String())+4), name, version)
}
