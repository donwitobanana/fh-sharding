package storage

import (
	"os"
	"path"
)

type Service interface {
	Write(dest string, data []byte) error
}

type service struct{
	workingDir string
}

func NewService(workingDir string) Service {
	return &service{
		workingDir: workingDir,
	}
}

func (s *service) Write(dest string, data []byte) error {
	file, err := os.Create(path.Join(s.workingDir, dest))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
