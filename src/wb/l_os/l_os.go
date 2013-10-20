package main

import "fmt"
import "os"
import "path/filepath"

type Visitor struct {
	ret string
}

func walkFn(path string, info os.FileInfo, err error) error {

	fmt.Printf("is dir %v \n", info.IsDir())
	fmt.Printf("file name %v \n", info.Name())

	if info.IsDir() && path != "D:/github/go_code/src" {
		return filepath.SkipDir //如果不想遍历某目录下的文件 则返回 filepath.SkipDir
	}

	fmt.Printf("file mode %v \n", info.Mode())
	fmt.Printf("file size %v \n", info.Size())
	fmt.Println("====================================")

	return err
}

func createFile(path string) bool {
	dir, _ := filepath.Split(path)

	os.MkdirAll(dir, os.ModeDir)

	io_file, _ := os.Create(path)
	if io_file != nil {
		io_file.Close()
	}

	return true
}

func main() {
	// v := &Visitor{}
	//os.MkdirAll("D:/github1/go_code/src", os.ModeDir)
	createFile("D:/github1/go_code/src/gooo.go")

	dir, file := filepath.Split("D:/github/go_code/src")
	fmt.Println("On Unix:", file, dir)
	//filepath.Walk("D:/github/go_code/src", walkFn)
}
