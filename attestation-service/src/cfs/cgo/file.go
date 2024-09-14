package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/confidential-filesystems/filesystem-ownership/utils"
)

type File struct {
	Resource

	RepoDir string
}

func NewFile(config Config) (Resource, error) {
	repoDir := config.StoreFileRepoDir // defaultRepoDir
	err := os.MkdirAll(repoDir, os.ModePerm)
	fmt.Printf("confilesystem-go - NewFile: os.MkdirAll(%v) -> err = %v\n",
		repoDir, err)
	if err != nil {
		return nil, err
	}
	file := &File{
		RepoDir: repoDir,
	}
	return file, nil
}

func (f *File) SetResource(repoDir, addr, typ, tag string, data []byte) error {
	fmt.Printf("confilesystem-go - File.SetResource(): addr = %v, typ = %v, tag = %v, repoDir = %v, f.RepoDir = %v\n",
		addr, typ, tag, repoDir, f.RepoDir)

	if repoDir == "" {
		repoDir = f.RepoDir
	}
	resourcePath := path.Join(repoDir, addr, typ, tag)
	folder := getDir(resourcePath)
	fmt.Printf("confilesystem-go - File.SetResource(): resourcePath = %v folder = %v\n", resourcePath, folder)
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(resourcePath, data, 0644)
	return err
}

func (f *File) DeleteResource(repoDir, addr, typ, tag string) error {
	fmt.Printf("confilesystem-go - File.DeleteResource(): addr = %v, typ = %v, tag = %v, repoDir = %v, f.RepoDir = %v\n",
		addr, typ, tag, repoDir, f.RepoDir)
	if repoDir == "" {
		repoDir = f.RepoDir
	}
	resourcePath := path.Join(repoDir, addr, typ, tag)
	fmt.Printf("confilesystem-go - File.DeleteResource(): resourcePath = %v\n", resourcePath)
	err := os.Remove(resourcePath)
	return err
}

func (f *File) GetResource(repoDir, addr, typ, tag, extraRequest string) ([]byte, error) {
	fmt.Printf("confilesystem-go - File.GetResource(): addr = %v, typ = %v, tag = %v, extraRequest = %v, repoDir = %v, f.RepoDir = %v\n",
		addr, typ, tag, extraRequest, repoDir, f.RepoDir)
	if repoDir == "" {
		repoDir = f.RepoDir
	}

	return utils.ToGetResource(context.Background(), repoDir, addr, typ, tag, extraRequest)
}

// utils api
func getDir(path string) string {
	return path[:len(path)-len(filepath.Base(path))]
}
