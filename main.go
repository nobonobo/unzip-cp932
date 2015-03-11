//usr/bin/env go run $0 $@; exit $?
package main

import (
	"archive/zip"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func toUtf8(str string, t transform.Transformer) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), t))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func unzip(src, dest string, t transform.Transformer) error {
	zc, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zc.Close()

	for _, item := range zc.File {
		fname, err := toUtf8(item.Name, t)
		if err != nil {
			return err
		}
		path := filepath.Join(dest, fname)
		if item.FileInfo().IsDir() {
			if err := os.MkdirAll(path, item.Mode()); err != nil {
				return err
			}
		} else {
			output, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, item.Mode())
			if err != nil {
				return err
			}
			defer output.Close()
			fp, err := item.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			if _, err := io.Copy(output, fp); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	output := "./"
	flag.StringVar(&output, "o", output, "destination folder")
	flag.Parse()
	err := unzip(flag.Arg(0), output, japanese.ShiftJIS.NewDecoder())
	if err != nil {
		log.Fatal(err)
	}
}
