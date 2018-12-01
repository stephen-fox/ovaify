package ovaify

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path"
)

// OvaConfig provides a configuration for creating a .ova file.
type OvaConfig struct {
	// OvfFilePath is the path to the original .ovf file.
	OvfFilePath string

	// FilePathsToInclude are the paths to the files to included in the .ova.
	FilePathsToInclude []string

	// OutputFilePath is the file path where the .ova should be saved.
	OutputFilePath string
}

// Validate returns a non-nil error if the configuration is invalid.
func (o *OvaConfig) Validate() error {
	if len(o.OvfFilePath) == 0 {
		return errors.New("Please specify a .ovf file path")
	}

	if len(o.FilePathsToInclude) == 0 {
		return errors.New("Please specify files to include in the .ova")
	}

	if len(o.OutputFilePath) == 0 {
		return errors.New("Please specify an output file path for the .ova")
	}

	return nil
}

// ConvertOvfToOva converts the specified .ovf and its associated files into
// a single .ova.
//
// Based on work by Daniel "thebsdbox" Finneran:
// https://github.com/thebsdbox/gova
func ConvertOvfToOva(config OvaConfig) error {
	err := config.Validate()
	if err != nil {
		return err
	}

	ovfInfo, err := os.Stat(config.OvfFilePath)
	if err != nil {
		return err
	}

	ova, err := os.OpenFile(config.OutputFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, ovfInfo.Mode())
	if err != nil {
		return err
	}
	defer ova.Close()

	tarw := tar.NewWriter(ova)

	err = CopyFileToOva(config.OvfFilePath, tarw)
	if err != nil {
		return err
	}

	for _, filePath := range config.FilePathsToInclude {
		err = CopyFileToOva(filePath, tarw)
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyFileToOva copies the specified file into a .ova.
//
// Based on work by Daniel "thebsdbox" Finneran:
// https://github.com/thebsdbox/gova
func CopyFileToOva(filePath string, ova *tar.Writer) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	header := &tar.Header{}
	header.Name = path.Base(filePath)
	header.Size = info.Size()
	header.Mode = int64(info.Mode())
	header.ModTime = info.ModTime()

	err = ova.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(ova, f)
	if err != nil {
		return err
	}

	return nil
}

