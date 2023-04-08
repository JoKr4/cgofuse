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
	"fmt"
	"log"
	"os"

	"github.com/winfsp/cgofuse/fuse"
)

const (
	filename = "hello"
	contents = "hello, world\n"
)

type Hellofs struct {
	fuse.FileSystemBase
	content *os.File
	len     int64
}

func (self *Hellofs) Open(path string, flags int) (errc int, fh uint64) {
	switch path {
	case "/" + filename:
		return 0, 0
	default:
		return -fuse.ENOENT, ^uint64(0)
	}
}

func (self *Hellofs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	fmt.Println("Getattr", path)
	switch path {
	case "/":
		stat.Mode = fuse.S_IFDIR | 0555
		//stat.Mode = fuse.S_IFDIR | 0777
		stat.Gid = 0
		stat.Uid = 0
		stat.Size = 4096
		stat.Nlink = 2
		return 0
	case "/" + filename:
		stat.Mode = fuse.S_IFREG | 0444
		//stat.Mode = fuse.S_IFREG | 0777
		//stat.Size = int64(len(contents))
		stat.Size = self.len
		stat.Gid = 0
		stat.Uid = 0
		stat.Nlink = 1
		return 0
	default:
		return -fuse.ENOENT
	}
}

func (self *Hellofs) Chown(path string, uid uint32, gid uint32) int {
	// fmt.Println(uid, gid)
	// if uid == 1000 && gid == 1000 {
	// 	return 0
	// }
	return -fuse.ENOSYS
}

func (self *Hellofs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	// endofst := ofst + int64(len(buff))
	// if endofst > self.len {
	// 	endofst = self.len
	// }
	// if endofst < ofst {
	// 	return 0
	// }
	n, _ = self.content.ReadAt(buff, ofst)
	return
}

func (self *Hellofs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	// s := fuse.Stat_t{
	// 	Uid:  1000,
	// 	Gid:  1000,
	// 	Mode: 755,
	// }
	fill(".", nil, 0)
	fill("..", nil, 0)
	fill(filename, nil, 0)
	return 0
}

func main() {
	f, err := os.Open("/home/johannes/filme.txt")
	if err != nil {
		log.Fatal(err)
	}
	s, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	hellofs := &Hellofs{
		content: f,
		len:     s.Size(),
	}
	host := fuse.NewFileSystemHost(hellofs)
	host.Mount(os.Args[1], []string{"-o", "ro,allow_other"})
}
