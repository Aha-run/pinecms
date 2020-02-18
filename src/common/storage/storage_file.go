package storage

import (
	"io"
	"os"
	"path/filepath"
)

type FileUploader struct {
	fixDir  string
	baseDir string
}

func NewFileUploader(fixDir, uploadDir string) *FileUploader {
	return &FileUploader{
		fixDir:  fixDir,
		baseDir: uploadDir,
	}
}

// Upload 上传图片
// storageName 云端路径名.
// LocalFile 要上传的文件名
func (s *FileUploader) Upload(storageName string, LocalFile io.Reader) (string, error) {
	//检测是否可以生成目录
	originName := storageName
	storageName = filepath.Join(s.baseDir, storageName)
	uploadDir := filepath.Dir(storageName)
	f, err := os.Open(uploadDir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	f.Close()
	out, err := os.OpenFile(storageName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, LocalFile)
	if err != nil {
		return "", err
	}
	return filepath.Join(s.fixDir, originName), nil
}
