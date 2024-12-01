package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	dataBlockSize := 5
	sliceLen := len(hashStr) / dataBlockSize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*dataBlockSize, i*dataBlockSize+dataBlockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		pathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

const defaultRootFolderName = "novastore"

type PathTransformFunc func(string) PathKey

type PathKey struct {
	pathName string
	FileName string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.pathName, p.FileName)
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.pathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

type StoreOpts struct {
	// Root is the topmost folder containing all files/folders
	Root              string
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		pathName: key,
		FileName: key,
	}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullFilePathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	_, err := os.Stat(fullFilePathWithRoot)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	firstPathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())

	defer func() {
		log.Printf("Deleted file: [%s] from disk", s.Root+"/"+pathKey.FullPath())
	}()

	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	return os.Open(fullPathWithRoot)
}

func (s *Store) streamWrite(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)

	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.pathName)

	err := os.MkdirAll(pathNameWithRoot, os.ModePerm)
	if err != nil {
		return err
	}

	fullFilePathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	f, err := os.Create(fullFilePathWithRoot)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("Written (%d) bytes to disk: %s", n, fullFilePathWithRoot)

	return nil
}
