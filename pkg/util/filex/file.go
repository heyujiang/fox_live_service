package filex

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var End = errors.New("end") // 正常提前结束，不再继续读取文件

type ReadLineFunc func(string) (interface{}, error)
type DealLineFunc func(int, []byte) error

func ReadLine(filename string, fn ReadLineFunc) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
		if err != nil {

		}
	}(file)
	reader := bufio.NewReader(file)
	var resp []interface{}
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		i, err := fn(string(line))
		if err != nil {
			return nil, err
		}
		resp = append(resp, i)
	}
	return json.Marshal(resp)
}

func DealLine(filename string, fn DealLineFunc) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	reader := bufio.NewReader(file)
	var lineNum = 1
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		err = fn(lineNum, line)
		if err != nil {
			if err == End {
				break
			}
			return err
		}
		lineNum++
	}
	return nil
}

func ReadLineByNumber(filename string, lineNum int) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	reader := bufio.NewReader(file)
	var i = 1
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if i == lineNum {
			return line, nil
		}
		i++
	}
	return nil, fmt.Errorf("invalid lineNum: %d", lineNum)
}

func Copy(dst, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(srcFile *os.File) {
		_ = srcFile.Close()
	}(srcFile)
	st, err := os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return nil
	}
	if st != nil {
		return nil
	}
	if err = MkdirIfNotExists(path.Dir(dst)); err != nil {
		return err
	}
	writeFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(writeFile *os.File) {
		_ = writeFile.Close()
	}(writeFile)
	_, err = io.Copy(writeFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

func CheckFiles(file string) (bool, error) {
	_, err := os.Stat(file)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func MkdirIfNotExists(dir string) error {
	stat, err := os.Stat(dir)
	if err != nil && !os.IsNotExist(err) {
		return nil
	}
	if stat != nil {
		return nil
	}
	return os.MkdirAll(dir, os.ModePerm)
}

func WriteFile(filename string, data []byte) error {
	if err := MkdirIfNotExists(filepath.Dir(filename)); err != nil {
		return err
	}
	if err := os.WriteFile(filename, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// WriteFileWithLine 逐行写入文件
func WriteFileWithLine(filename string, data []string) error {
	if err := MkdirIfNotExists(filepath.Dir(filename)); err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer file.Close()
	write := bufio.NewWriter(file)

	for _, v := range data {
		_, err := write.WriteString(fmt.Sprintf("%s\n", v))
		if err != nil {
			return err
		}
	}

	return write.Flush()
}

func ReadFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// GetAllDirsAndFiles 获取文件夹下所有的文件夹和文件
func GetAllDirsAndFiles(dir string) ([]string, []string, error) {
	var dirs []string
	var files []string

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}

	for _, f := range fileInfos {
		if f.IsDir() {
			dirs = append(dirs, path.Join(dir, f.Name()))
			dir1, dirFiles, err := GetAllDirsAndFiles(path.Join(dir, f.Name()))
			if err == nil {
				files = append(files, dirFiles...)
				dirs = append(dirs, dir1...)
			}
			continue
		}
		files = append(files, path.Join(dir, f.Name()))
	}
	return dirs, files, err
}

func GetAllFiles(dir string) ([]string, error) {
	var files []string

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range fileInfos {
		if !f.IsDir() {
			files = append(files, path.Join(dir, f.Name()))
		}
	}
	return files, err
}
