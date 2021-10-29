package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

type NotifyFile struct {
	watch *fsnotify.Watcher
}

func NewNotify() *NotifyFile {
	notify := &NotifyFile{}
	notify.watch, _ = fsnotify.NewWatcher()

	return notify
}

func (slf *NotifyFile) Watch(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = slf.watch.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("watching: ", path)
		}
		return nil
	})

	go slf.WatchEvent()
}

func (slf *NotifyFile) WatchEvent() {
	for {
		select {
		case ev := <-slf.watch.Events:
			if ev.Op&fsnotify.Create == fsnotify.Create {
				fmt.Println("创建文件 : ", ev.Name)
				//获取新创建文件的信息，如果是目录，则加入监控中
				file, err := os.Stat(ev.Name)
				if err == nil && file.IsDir() {
					slf.watch.Add(ev.Name)
					fmt.Println("添加监控 : ", ev.Name)
				}
			}
			if ev.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("写入文件 : ", ev.Name)
			}
		case err := <-slf.watch.Errors:
			fmt.Println("error: ", err)
		}
	}
}

func main() {
	watch := NewNotify()
	watch.Watch("D:\\tmp")
	select {}
}
