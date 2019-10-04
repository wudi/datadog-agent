// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package http

import (
	"bytes"
	"compress/gzip"
	"net/http"
)

// Compression compresses the payload
type Compression interface {
	compress(payload []byte) ([]byte, error)
	setHeader(header *http.Header)
}

// NoCompression does not compress the payload
var NoCompression Compression = &noCompression{}

type noCompression struct{}

func (c *noCompression) compress(payload []byte) ([]byte, error) {
	return payload, nil
}

func (c *noCompression) setHeader(header *http.Header) {
}

// GzipCompression compresses the payload using Gzip algorithm
type GzipCompression struct {
	level int
}

// NewGzipCompression creates a new Gzip compression
func NewGzipCompression(level int) *GzipCompression {
	if level < gzip.NoCompression {
		level = gzip.NoCompression
	} else if level > gzip.BestCompression {
		level = gzip.BestCompression
	}

	return &GzipCompression{
		level,
	}
}

func (c *GzipCompression) compress(payload []byte) ([]byte, error) {
	var compressedPayload bytes.Buffer
	gzipWriter, err := gzip.NewWriterLevel(&compressedPayload, c.level)
	if err != nil {
		return nil, err
	}
	_, err = gzipWriter.Write(payload)
	if err != nil {
		return nil, err
	}
	err = gzipWriter.Flush()
	if err != nil {
		return nil, err
	}
	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}
	return compressedPayload.Bytes(), nil
}

func (c *GzipCompression) setHeader(header *http.Header) {
	header.Set("Content-Encoding", "gzip")
}
