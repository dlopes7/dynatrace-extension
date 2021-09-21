package downloader

import (
	"archive/zip"
	"fmt"
	"github.com/go-logr/logr"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func NewExtensionDownloader(log logr.Logger) ExtensionDownloader {
	return ExtensionDownloader{
		name: os.Getenv("DT_EXTENSION_NAME"),
		link: os.Getenv("DT_EXTENSION_LINK"),
		log:  log,
	}

}

type ExtensionDownloader struct {
	name string
	link string
	log  logr.Logger
}

func (e *ExtensionDownloader) getFileNameFromURL() (string, error) {
	fileURL, err := url.Parse(e.link)
	if err != nil {
		e.log.Error(err, "Could not parse the link", "link", e.link)
		return "", err
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	return segments[len(segments)-1], nil

}

func (e *ExtensionDownloader) Extract() error {

	installPath := "/plugin_deployment"

	fileName, err := e.getFileNameFromURL()
	if err != nil {
		return err
	}
	downloadPath := fmt.Sprintf("/%s/%s", installPath, fileName)
	files, err := e.unzip(downloadPath, installPath)
	if err != nil {
		e.log.Error(err, "could not extract files", "installPath", installPath, "downloadPath", downloadPath)
		return err
	}
	e.log.Info(fmt.Sprintf("Extracted %d files to %s", len(files), installPath))

	return nil

}

func (e *ExtensionDownloader) Download() error {
	e.log.Info("Starting download", "name", e.name, "link", e.link)

	installPath := "/plugin_deployment"
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		e.log.Error(err, "pluginDeploymentPath does not exist", "pluginDeploymentPath", installPath)
		return err
	}

	fileName, err := e.getFileNameFromURL()
	if err != nil {
		return err
	}
	downloadPath := fmt.Sprintf("/%s/%s", installPath, fileName)

	resp, err := http.Get(e.link)
	if err != nil {
		e.log.Error(err, "could not download", "link", e.link)
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(downloadPath)
	if err != nil {
		e.log.Error(err, "could not create pluginDeploymentPath", "pluginDeploymentPath", downloadPath)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		e.log.Error(err, "could not write to file", "pluginDeploymentPath", downloadPath)
		return err
	}
	e.log.Info("downloaded the file", "pluginDeploymentPath", downloadPath)
	return nil
}

func (e *ExtensionDownloader) unzip(src string, dest string) ([]string, error) {
	e.log.Info(fmt.Sprintf("Attempting to extract from '%s' to '%s'", src, dest))

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func (e *ExtensionDownloader) CheckIfDownloaded() bool {
	fileName, err := e.getFileNameFromURL()
	if err != nil {
		return false
	}
	path := fmt.Sprintf("/plugin_deployment/%s", fileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		e.log.Info("path does not exist", "path", path)
		return false
	}
	return true
}
