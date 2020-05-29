package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	AzureProvider string = "azure"
)

func init() {
	// TODO: since we are using rand, we need this seed somewhere better
	rand.Seed(time.Now().UnixNano())
}

type ImageDownUploader interface {
	DownloadImage(container, imageKey, zone string) ([]byte, error)
	UploadImage(data []byte, fileKey, container, zone string) error
}

func initDownloadUploader(dzConf DZFilesConfig) (ImageDownUploader, error) {
	switch dzConf.Provider {
	case "azure":
		source := NewAzureImageSource(nil).(ImageDownUploader)
		return source, nil
	case "s3":

	}

	return nil, fmt.Errorf("dzfiles: unknown provider")
}

type DZFilesConfig struct {
	Provider string // azure ||  s3

	ImageKey      string `json:"imageKey"`
	Container     string `json:"container"`
	TempContainer string `json:"tempContainer"`

	ContainerZone string `json:"containerZone"` // container zone (s3 region)

	// TODO: sas
}

func UploadDZFiles(dzConf DZFilesConfig) error {
	downUploader, err := initDownloadUploader(dzConf)
	if err != nil {
		return fmt.Errorf("dzfiles: error getting source: %w", err)
	}

	// TODO: this is just an ugly hack which is terrible, this needs to be solved.
	go func() (err error) {
		keyDir, imageName := filepath.Split(dzConf.ImageKey)
		imageName = imageName[:len(imageName)-len(filepath.Ext(imageName))]

		// TODO: hack to defer error
		defer func() {
			if err != nil {
				downUploader.UploadImage(
					[]byte(err.Error()),
					filepath.Join(keyDir, imageName+".txt"),
					dzConf.TempContainer,
					dzConf.ContainerZone,
				)
				fmt.Printf("dzfiles: error: %s", err)
			}
		}()

		localDirPath := fmt.Sprintf("/tmp/dzFiles-%d", rand.Uint64())
		if err := os.Mkdir(localDirPath, 0777); err != nil {
			return fmt.Errorf("dzfiles: error creating tmp dir: %w", err)
		}
		defer os.RemoveAll(localDirPath)

		if err := downUploader.UploadImage(
			[]byte("pending"),
			filepath.Join(keyDir, imageName+".txt"),
			dzConf.TempContainer,
			dzConf.ContainerZone,
		); err != nil {
			return fmt.Errorf("dzfiles: error creating txt file: %w", err)
		}

		data, err := downUploader.DownloadImage(dzConf.Container, dzConf.ImageKey, dzConf.ContainerZone)
		if err != nil {
			return fmt.Errorf("dzfiles: error downloading image: %w", err)
		}

		if err := generateDZFiles(localDirPath, data, imageName); err != nil {
			return fmt.Errorf("dzfiles: error generating dz files: %w", err)
		}

		var g errgroup.Group
		g.Go(func() error {
			data, err := ioutil.ReadFile(filepath.Join(localDirPath, imageName+".dzi"))
			if err != nil {
				return fmt.Errorf("dzfiles: error reading index file: %w", err)
			}

			if err := downUploader.UploadImage(
				data,
				filepath.Join(keyDir, imageName+".dzi"),
				dzConf.TempContainer,
				dzConf.ContainerZone,
			); err != nil {
				return fmt.Errorf("dzfiles: error uploading index file: %w", err)
			}

			return nil
		})

		if err := filepath.Walk(
			filepath.Join(localDirPath, imageName+"_files"),
			filepath.WalkFunc(
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					if info.IsDir() {
						return nil
					}

					g.Go(func() error {
						data, err := ioutil.ReadFile(path)
						if err != nil {
							return fmt.Errorf("dzfiles: error reading file: %s: %w", path, err)
						}

						if err := downUploader.UploadImage(
							data,
							keyDir+path[len(localDirPath)+1:], // +1 -> for slash "/",
							dzConf.TempContainer,
							dzConf.ContainerZone,
						); err != nil {
							return fmt.Errorf("dzfiles: error uploading file: %s: %w", path, err)
						}

						return nil
					})

					return nil
				},
			),
		); err != nil {
			return fmt.Errorf("dzfiles: error walking dir path: %w", err)
		}

		if err := g.Wait(); err != nil {
			return err
		}

		if err := downUploader.UploadImage(
			[]byte("ok"),
			filepath.Join(keyDir, imageName+".txt"),
			dzConf.TempContainer,
			dzConf.ContainerZone,
		); err != nil {
			return fmt.Errorf("dzfiles: error creating ok txt file: %s", err)
		}

		return nil
	}()

	return nil
}

// generateDZFiles generates dz files in given dir with same image name, just changed
// extension.
func generateDZFiles(dirPath string, data []byte, imageName string) error {
	imagePath := fmt.Sprintf("%s/%s", dirPath, imageName)

	tiffImagePath := imagePath + ".tiff"
	if err := ioutil.WriteFile(tiffImagePath, data, 0777); err != nil {
		return fmt.Errorf("dzfiles: error saving tiff image to disk: %w", err)
	}

	if err := exec.Command("vips", "dzsave", tiffImagePath, imagePath).Run(); err != nil {
		return fmt.Errorf("dzfiles: error creating dz files: %w", err)
	}

	return nil
}
