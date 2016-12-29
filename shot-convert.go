// Command shot-convert watches specified directory for new png screenshots and
// convert them to jpeg images
package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/artyom/autoflags"
	"github.com/fsnotify/fsnotify"
)

func main() {
	args := struct {
		Dir    string `flag:"dir,directory to watch"`
		Prefix string `flag:"prefix,file prefix"`
		Q      int    `flag:"q,jpeg quality"`
	}{Dir: "$HOME/Downloads", Prefix: "Screen Shot", Q: 90}
	autoflags.Define(&args)
	flag.Parse()
	if err := do(os.ExpandEnv(args.Dir), args.Prefix, args.Q); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func do(dir, prefix string, quality int) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	// defer w.Close() // https://github.com/fsnotify/fsnotify/issues/187
	if err := w.Add(dir); err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nameCh := make(chan string, 10)
	errCh := make(chan error, 1)
	go func() {
		for {
			select {
			case name := <-nameCh:
				time.Sleep(500 * time.Millisecond)
				if err := convert(name, quality); err != nil {
					errCh <- err
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	for {
		select {
		case err := <-errCh:
			return err
		case evt := <-w.Events:
			if evt.Op&fsnotify.Create != fsnotify.Create {
				continue
			}
			if b := filepath.Base(evt.Name); strings.HasPrefix(b, prefix) &&
				strings.HasSuffix(b, suffix) {
				nameCh <- evt.Name
			}
		case err := <-w.Errors:
			return err
		}
	}
}

func convert(name string, quality int) error {
	jpegName := strings.TrimSuffix(name, suffix) + ".jpg"
	if _, err := os.Stat(jpegName); err == nil {
		return nil
	}
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	tf, err := ioutil.TempFile(filepath.Dir(name), ".shot-convert-")
	if err != nil {
		return err
	}
	defer tf.Close()
	defer os.Remove(tf.Name())
	if err := jpeg.Encode(tf, img, &jpeg.Options{Quality: quality}); err != nil {
		return err
	}
	if err := tf.Close(); err != nil {
		return err
	}
	return os.Rename(tf.Name(), jpegName)
}

const suffix = ".png"
