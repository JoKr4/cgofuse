/*
 * hellofs.go
 *
 * Copyright 2017-2022 Bill Zissimopoulos
 */
/*
 * This file is part of Cgofuse.
 *
 * It is licensed under the MIT license. The full license text can be found
 * in the License.txt file at the root of this project.
 */

package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/winfsp/cgofuse/fuse"
)

type Hellofs struct {
	fuse.FileSystemBase
	files []*os.File
	names []string
}

func (self *Hellofs) Open(path string, flags int) (errc int, fh uint64) {

	for _, n := range self.names {
		if n == filepath.Base(path) {
			return 0, 0
		}
	}

	return -fuse.ENOENT, ^uint64(0)
}

func (self *Hellofs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {

	if path == "/" {
		stat.Mode = fuse.S_IFDIR | 0555
		stat.Gid = 0
		stat.Uid = 0
		stat.Size = 4096
		stat.Nlink = 2
		return 0
	}

	for i, n := range self.names {
		if n == filepath.Base(path) {
			// https://stackoverflow.com/questions/46558824/how-do-i-get-a-block-devices-size-correctly-in-go
			pos, err := self.files[i].Seek(0, io.SeekEnd)
			if err != nil {
				log.Printf("error seeking to end: %s\n", err)
				continue
			}
			stat.Size = pos
			stat.Mode = fuse.S_IFREG | 0444
			stat.Gid = 0
			stat.Uid = 0
			stat.Nlink = 1
			return 0
		}
	}

	return -fuse.ENOENT
}

func (self *Hellofs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	for i, name := range self.names {
		if name == filepath.Base(path) {
			n, _ = self.files[i].ReadAt(buff, ofst)
			return n
		}
	}
	return 0
}

func (self *Hellofs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	fill(".", nil, 0)
	fill("..", nil, 0)
	for _, n := range self.names {
		fill(n, nil, 0)
	}
	return 0
}

func main() {

	mountpoint := "/mnt/optical-drives"

	drives := []string{
		"/dev/sr0",
		"/dev/sr1",
	}

	hellofs := &Hellofs{
		files: make([]*os.File, 0, len(drives)),
		names: make([]string, 0, len(drives)),
	}

	for _, d := range drives {
		f, err := os.Open(d)
		if err != nil {
			log.Println(err)
			continue
		}
		mountAs := filepath.Base(f.Name()) + ".iso"
		hellofs.files = append(hellofs.files, f)
		hellofs.names = append(hellofs.names, mountAs)
		log.Printf("opened '%s' to mount as '%s'", f.Name(), mountAs)
	}

	host := fuse.NewFileSystemHost(hellofs)
	host.Mount(mountpoint, []string{"-o", "ro,allow_other"})
}
