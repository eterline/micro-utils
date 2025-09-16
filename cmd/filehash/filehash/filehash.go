// Copyright (c) 2025 EterLine (Andrew)
// This file is part of My-Go-Project.
// Licensed under the MIT License. See the LICENSE file for details.


package filehash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type HashBytes []byte

func (hb HashBytes) String() string {
	if hb == nil {
		return "null"
	}
	return hex.EncodeToString(hb)
}

func (hb HashBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hb.String())
}

type FileData struct {
	Filename string    `json:"filename,omitempty"`
	Size     int64     `json:"size"`
	MD5      HashBytes `json:"md5"`
	SHA1     HashBytes `json:"sha1"`
	SHA256   HashBytes `json:"sha256"`
}

func CalcFileHash(name string) (*FileData, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("file opening error: %w", err)
	}
	defer file.Close()

	data := &FileData{
		Filename: file.Name(),
	}

	hMd5 := md5.New()
	hSha1 := sha1.New()
	hSha256 := sha256.New()

	multiWriter := io.MultiWriter(hMd5, hSha1, hSha256)

	written, err := io.Copy(multiWriter, file)
	if err != nil {
		return nil, fmt.Errorf("failed to compute hashes: %w", err)
	}

	data.Size = written
	data.MD5 = hMd5.Sum(nil)
	data.SHA1 = hSha1.Sum(nil)
	data.SHA256 = hSha256.Sum(nil)

	return data, nil
}

func CalcPipelineHash() (*FileData, error) {

	data := &FileData{
		Filename: "stdin",
	}

	hMd5 := md5.New()
	hSha1 := sha1.New()
	hSha256 := sha256.New()

	multiWriter := io.MultiWriter(hMd5, hSha1, hSha256)

	written, err := io.Copy(multiWriter, os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to compute hashes: %w", err)
	}

	data.Size = written
	data.MD5 = hMd5.Sum(nil)
	data.SHA1 = hSha1.Sum(nil)
	data.SHA256 = hSha256.Sum(nil)

	return data, nil
}
