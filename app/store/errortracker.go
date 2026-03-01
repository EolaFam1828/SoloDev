// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"crypto/sha256"
	"encoding/hex"
)

// Fingerprint generates a fingerprint hash from title and stack trace.
func Fingerprint(title, stackTrace string) string {
	hash := sha256.Sum256([]byte(title + stackTrace))
	return hex.EncodeToString(hash[:])
}
