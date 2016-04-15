//usr/bin/env go run $0 $@; exit $?
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexmullins/zip"
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

func unzip(src, dest, passphrase string, t transform.Transformer) error {
	zc, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zc.Close()

	for _, item := range zc.File {
		fname, err := toUtf8(item.Name, t)
		if err != nil {
			fname = item.Name
		}
		if item.IsEncrypted() {
			item.SetPassword(passphrase)
		}
		path := filepath.Join(dest, fname)
		if item.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		} else {
			dir := filepath.Dir(path)
			if _, err := os.Stat(dir); !os.IsExist(err) {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return err
				}
			}
			output, err := os.Create(path)
			if err != nil {
				return err
			}
			defer output.Close()
			fp, err := item.Open()
			if err != nil {
				return err
			}
			defer fp.Close()
			if _, err := io.Copy(output, fp); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	dest := "./"
	passphrase := ""
	flag.StringVar(&dest, "d", dest, "destination folder")
	flag.StringVar(&passphrase, "p", passphrase, "passphrase")
	flag.Parse()
	fmt.Println("dest:", dest)
	err := unzip(flag.Arg(0), dest, passphrase, japanese.ShiftJIS.NewDecoder())
	if err != nil {
		log.Fatal(err)
	}
}
