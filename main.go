package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/schollz/progressbar/v3"
)

const (
	name    = "DiskIO - Disk performance metering"
	site    = "https://github.com/dassump/diskio"
	author  = "Daniel Dias de Assumpção <dassump@gmail.com>"
	version = "v1.0.0"
)

var (
	dir         string
	dir_key     = "dir"
	dir_default = "."
	dir_info    = "Destination directory"

	size         int64
	size_key     = "size"
	size_default = int64(100)
	size_info    = "Temporary file size in MiB"

	buffer         int64
	buffer_key     = "buffer"
	buffer_default = int64(1024)
	buffer_info    = "Buffer size in KiB"

	output  = "\nSize: %s\nBuffer: %s\nElapsed time: %s\nAverage rate: %.2f MB/s\n"
	pattern = "tmp-diskio-"
	quit    = make(chan os.Signal, 1)
)

func init() {
	flag.StringVar(&dir, dir_key, dir_default, dir_info)
	flag.Int64Var(&size, size_key, size_default, size_info)
	flag.Int64Var(&buffer, buffer_key, buffer_default, buffer_info)

	flag.Usage = func() {
		fmt.Printf(
			"%s\n%s\n\nAuthor: %s\nVersion: %s\n\n",
			name, site, author, version,
		)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
}

func main() {
	file, err := os.CreateTemp(dir, pattern)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		_ = file.Close()

		if err := os.Remove(file.Name()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	total := size * (1024 * 1024)
	bar := progressbar.DefaultBytes(total)
	buf := make([]byte, buffer*1024)
	writer := io.MultiWriter(file, bar)
	start := time.Now()

	go func() {
		for i := int64(0); i < total; i += int64(len(buf)) {
			_, _ = rand.Read(buf)
			_, err := writer.Write(buf)

			if err != nil && !errors.Is(err, os.ErrClosed) {
				fmt.Println(err)

				return
			}
		}

		_ = file.Sync()

		quit <- os.Interrupt
	}()

	<-quit

	elapsed := time.Since(start)
	fmt.Printf(
		output,
		humanize.IBytes(uint64(bar.State().CurrentBytes)),
		humanize.IBytes(uint64(buffer)*1024),
		elapsed.Round(time.Millisecond).String(),
		(bar.State().CurrentBytes/1000000)/(float64(elapsed)/float64(time.Second)),
	)
}
