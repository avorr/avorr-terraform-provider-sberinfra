package imports

import (
	"bufio"
	"os"
	"path/filepath"
)

type FileWriter struct {
	FileName    string
	AbsFileName string
	Dir         string
	Buffer      *bufio.Writer
	File        *os.File
}

func NewFileWriter(name string) *FileWriter {
	fw := &FileWriter{FileName: name}

	cdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fw.AbsFileName = filepath.Join(cdir, name)
	fw.Dir = filepath.Dir(fw.AbsFileName)
	err = os.MkdirAll(fw.Dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	newFile, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	fw.File = newFile
	fw.Buffer = bufio.NewWriter(fw.File)
	return fw
}

func (o *FileWriter) Close() error {
	err := o.Buffer.Flush()
	if err != nil {
		return err
	}
	err = o.File.Close()
	if err != nil {
		return err
	}
	return nil
}
