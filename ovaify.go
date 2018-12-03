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
	// OvfFilePath is the OVF file to include in the OVA.
	OvfFilePath string

	// FilePathsToInclude are the paths to the files to included
	// in the OVA.
	FilePathsToInclude []string

	// OutputFileMode specifies the mode for the resulting
	// OVA file. If it is not specified, then the OVF's mode
	// is used.
	OutputFileMode os.FileMode

	// OutputFilePath is the file path where the OVA will be saved.
	OutputFilePath string
}

// Validate returns a non-nil error if the configuration is invalid.
func (o *OvaConfig) Validate() error {
	if len(o.OvfFilePath) == 0 {
		return errors.New("Please specify an OVF file")
	}

	ovfInfo, err := os.Stat(o.OvfFilePath)
	if err != nil {
		return err
	}

	if o.OutputFileMode == 0 {
		o.OutputFileMode = ovfInfo.Mode()
	}

	if len(o.FilePathsToInclude) == 0 {
		return errors.New("Please specify files to include in the .ova")
	}

	if len(o.OutputFilePath) == 0 {
		return errors.New("Please specify an output file path for the .ova")
	}

	return nil
}

// CreateOvaFile creates an OVA file using the specified OvaConfig.
//
// Based on work by Daniel "thebsdbox" Finneran:
// https://github.com/thebsdbox/gova
func CreateOvaFile(config OvaConfig) error {
	err := config.Validate()
	if err != nil {
		return err
	}

	ova, err := os.OpenFile(config.OutputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, config.OutputFileMode)
	if err != nil {
		return err
	}
	defer ova.Close()

	tarw := tar.NewWriter(ova)

	err = CopyFileIntoOva(config.OvfFilePath, tarw)
	if err != nil {
		return err
	}

	for _, filePath := range config.FilePathsToInclude {
		err = CopyFileIntoOva(filePath, tarw)
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyFileIntoOva copies the specified file into a .ova.
//
// Based on work by Daniel "thebsdbox" Finneran:
// https://github.com/thebsdbox/gova
func CopyFileIntoOva(filePath string, ova *tar.Writer) error {
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
