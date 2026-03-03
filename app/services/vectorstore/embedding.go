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

package vectorstore

import (
	"math"
	"regexp"
	"strings"
	"unicode"
)

// EmbeddingDimension is the fixed size of our TF-IDF embedding vectors.
// Using a hash-based approach (feature hashing) to avoid needing a global vocabulary.
const EmbeddingDimension = 512

var (
	camelRe    = regexp.MustCompile(`([a-z])([A-Z])`)
	nonAlphaRe = regexp.MustCompile(`[^a-zA-Z0-9]+`)
)

// Embed produces a fixed-size vector from text using feature hashing (hashing trick).
// This is a self-contained approach that requires no external API or trained model.
func Embed(text string) []float64 {
	tokens := tokenize(text)
	if len(tokens) == 0 {
		return make([]float64, EmbeddingDimension)
	}

	// Count term frequencies.
	tf := make(map[string]float64)
	for _, t := range tokens {
		tf[t]++
	}
	// Normalize by total tokens.
	total := float64(len(tokens))
	for k := range tf {
		tf[k] /= total
	}

	// Feature hashing: map each token to a dimension via FNV-like hash.
	vec := make([]float64, EmbeddingDimension)
	for token, freq := range tf {
		idx := hashToken(token) % uint32(EmbeddingDimension)
		// Use sign hashing to reduce collisions.
		if hashToken(token+"_sign")%2 == 0 {
			vec[idx] += freq
		} else {
			vec[idx] -= freq
		}
	}

	// L2-normalize the vector.
	norm := l2Norm(vec)
	if norm > 0 {
		for i := range vec {
			vec[i] /= norm
		}
	}

	return vec
}

// CosineSimilarity computes the cosine similarity between two vectors.
func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	denom := math.Sqrt(normA) * math.Sqrt(normB)
	if denom == 0 {
		return 0
	}
	return dot / denom
}

// tokenize splits text into normalized tokens suitable for embedding.
func tokenize(text string) []string {
	// Split camelCase: "getConfig" → "get Config"
	text = camelRe.ReplaceAllString(text, "${1} ${2}")
	// Replace non-alphanumeric with spaces.
	text = nonAlphaRe.ReplaceAllString(text, " ")
	text = strings.ToLower(text)

	parts := strings.Fields(text)
	var tokens []string
	for _, p := range parts {
		p = strings.TrimFunc(p, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		})
		if len(p) >= 2 {
			tokens = append(tokens, p)
		}
	}
	return tokens
}

// hashToken is a simple FNV-1a hash for deterministic feature hashing.
func hashToken(s string) uint32 {
	const (
		offset32 = 2166136261
		prime32  = 16777619
	)
	h := uint32(offset32)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= prime32
	}
	return h
}

func l2Norm(v []float64) float64 {
	var sum float64
	for _, x := range v {
		sum += x * x
	}
	return math.Sqrt(sum)
}
