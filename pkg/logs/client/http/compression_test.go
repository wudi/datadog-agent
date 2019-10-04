// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package http

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoCompression(t *testing.T) {
	payload := []byte("my payload")

	compressedPayload, err := NoCompression.compress(payload)
	assert.Nil(t, err)

	assert.Equal(t, payload, compressedPayload)
}

func TestNoCompressionHeader(t *testing.T) {
	header := make(http.Header)

	NoCompression.setHeader(&header)

	assert.Empty(t, header.Get("Content-Encoding"))
}

func TestGzipCompression(t *testing.T) {
	payload := []byte("my payload")

	compressedPayload, err := NewGzipCompression(gzip.BestCompression).compress(payload)
	assert.Nil(t, err)

	decompressedPayload, err := decompress(compressedPayload)
	assert.Nil(t, err)

	assert.Equal(t, payload, decompressedPayload)
}

func TestGzipCompressionHeader(t *testing.T) {
	header := make(http.Header)

	NewGzipCompression(gzip.BestCompression).setHeader(&header)

	assert.Equal(t, header.Get("Content-Encoding"), "gzip")
}

func decompress(payload []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
