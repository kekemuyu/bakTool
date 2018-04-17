package sync

import (
	//  "flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileIO struct {
}

func (f *FileIO) CopyDir(src string, dest string) {
	src_original := src
	err := filepath.Walk(src, func(src string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			err := os.Mkdir(dest, os.ModePerm) //在当前目录下生成md目录
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(src)
			fmt.Println(src_original)
			fmt.Println(dest)

			dest_new := strings.Replace(src, src_original, dest, -1)
			fmt.Println(dest_new)
			fmt.Println("CopyFile:" + src + " to " + dest_new)
			CopyFile(src, dest_new)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func (f *FileIO) Remove(file string) {
	err := os.Remove(file)
	if err != nil {
		fmt.Println(err)
	}
}

//egodic directories
func getFilelist(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		println(path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//copy file
func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()
	fmt.Println("dst:" + dst)
	dst_slices := strings.Split(dst, "\\")
	dst_slices_len := len(dst_slices)
	dest_dir := ""
	for i := 0; i < dst_slices_len-1; i++ {
		dest_dir = dest_dir + dst_slices[i] + "\\"
	}
	fmt.Println("dest_dir:" + dest_dir)
	b, err := PathExists(dest_dir)
	if b == false {
		err := os.Mkdir(dest_dir, os.ModePerm) //在当前目录下生成md目录
		if err != nil {
			fmt.Println(err)
		}
	}
	dstFile, err := os.Create(dst)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}
