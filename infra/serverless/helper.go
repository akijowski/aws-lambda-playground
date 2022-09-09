package serverless

import (
	"encoding/hex"
	"hash"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func mustCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

func hashDirectory(root string, h hash.Hash) (string, error) {
	hashes := make([]byte, 0)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			buf := make([]byte, 30*1024)
			if _, err := io.CopyBuffer(h, file, buf); err != nil {
				return err
			}
			sum := h.Sum(nil)
			// fmt.Printf("sum: %x\n", sum)
			h.Reset()
			hashes = append(hashes, sum...)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	h.Reset()
	h.Write(hashes)
	sum := h.Sum(nil)
	return hex.EncodeToString(sum), nil
}
