// Copyright 2013 Unknown
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package gopha is an implementation of Perceptual Hash Algorithm in Go programming language.
package gopha

import (
	"bytes"
	"image"

	//"github.com/nfnt/resize"
)

func PHA(m image.Image) string {
	// Step 1: resize picture to 8*8.
	//m = resize.Resize(8, 8, m, resize.NearestNeighbor)
	m = Resize(m, m.Bounds(), 8, 8)

	// Step 2: grayscale picture.
	gray := grayscaleImg(m)

	// Step 3: calculate average value.
	avg := calAvgValue(gray)

	// Step 4: get fingerprint.
	fg := getFingerprint(avg, gray)

	return string(fg)
}

// grayscaleImg converts picture to grayscale.
func grayscaleImg(src image.Image) []byte {
	// Create a new grayscale image
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	gray := make([]byte, w*h)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := src.At(x, y).RGBA()
			gray[x+y*8] = byte((r*30 + g*59 + b*11) / 100)
		}
	}
	return gray
}

// calAvgValue returns average value of color of picture.
func calAvgValue(gray []byte) byte {
	sum := 0
	for _, v := range gray {
		sum += int(v)
	}
	return byte(sum / len(gray))
}

// getFingerprint returns fingerprint of a picture.
func getFingerprint(avg byte, gray []byte) string {
	var buf bytes.Buffer
	for _, v := range gray {
		if avg >= v {
			buf.WriteByte('1')
		} else {
			buf.WriteByte('0')
		}
	}
	return buf.String()
}
