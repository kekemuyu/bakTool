package sync

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Sync struct {
	fi  FileIO
	src string
	des []string
}

func New(src string, des []string) *Sync {
	return &Sync{
		src: src,
		des: des,
	}
}

func (s *Sync) Run() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						fmt.Println("创建文件 : ", ev.Name)
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							watcher.Add(ev.Name)
							fmt.Println("添加监控 : ", ev.Name)
						}
						for _, v := range s.des {
							s.fi.CopyDir(ev.Name, v+ev.Name[strings.Index(ev.Name, s.src)+len(s.src):])
						}
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						fmt.Println("写入文件 : ", ev.Name)
						for _, v := range s.des {
							s.fi.CopyDir(ev.Name, v+ev.Name[strings.Index(ev.Name, s.src)+len(s.src):])
						}
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						fmt.Println("删除文件 : ", ev.Name)
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							watcher.Remove(ev.Name)
							fmt.Println("删除监控 : ", ev.Name)
							for _, v := range s.des {
								fmt.Println("删除目录：", v+ev.Name[strings.Index(ev.Name, s.src)+len(s.src):])
								s.fi.Remove(v + ev.Name[strings.Index(ev.Name, s.src)+len(s.src):])
							}
						} else {
							for _, v := range s.des {
								fmt.Println("删除文件：", v+ev.Name[strings.Index(ev.Name, s.src)+len(s.src):])
								s.fi.Remove(v + ev.Name[strings.Index(ev.Name, s.src)+len(s.src):])
							}
						}
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						fmt.Println("重命名文件 : ", ev.Name)
						//						fi, err := os.Stat(ev.Name)
						//						if err == nil && fi.IsDir() {
						//							fmt.Println("文件夹")
						//						} else {
						fmt.Println("文件")
						for _, v := range s.des {
							s.fi.Remove(v + ev.Name[strings.Index(ev.Name, s.src)+len(s.src):])
						}
						//						}
						watcher.Remove(ev.Name)
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						fmt.Println("修改权限 : ", ev.Name)
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(s.src)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
