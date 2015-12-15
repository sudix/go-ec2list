package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/mitchellh/go-homedir"
)

type Cache struct {
	ExpireMinute        int
	directory, filepath string
}

func NewCache(cachemin int) Cache {
	c := Cache{ExpireMinute: cachemin}
	c.setFilePath()
	return c
}

func (c *Cache) Delete() {
	if !c.exists() {
		return
	}

	if err := os.Remove(c.filepath); err != nil {
		log.Fatal(err)
	}
}

func (c *Cache) Use() bool {
	return c.ExpireMinute > 0
}

func (c *Cache) Available() bool {
	if !c.exists() {
		return false
	}
	if c.expired() {
		return false
	}
	return true
}

func (c *Cache) Save(list EC2List) {
	mkDir(c.directory)
	f, err := os.Create(c.filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	list.Output(f)
}

func (c *Cache) Output(w io.Writer) {
	f, err := os.Open(c.filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Cache) setFilePath() {
	home := getHomeDir()
	c.directory = path.Join(home, ".go-ec2list")
	c.filepath = path.Join(c.directory, "cache")
}

func (c *Cache) exists() bool {
	_, err := os.Stat(c.filepath)
	return err == nil
}

func (c *Cache) expired() bool {
	s, err := os.Stat(c.filepath)
	if err != nil {
		log.Fatal(err)
	}
	duration := time.Now().Sub(s.ModTime())
	return float64(c.ExpireMinute) < duration.Minutes()
}

func mkDir(dirPath string) {
	if err := os.MkdirAll(dirPath, 0777); err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}
}

func getHomeDir() string {
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
