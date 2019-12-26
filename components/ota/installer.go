package ota

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/kardianos/osext"
)

type Installer struct {
	shutdown func() error
}

func NewInstaller(shutdown func() error) *Installer {
	return &Installer{
		shutdown: shutdown,
	}
}

// очистка старых не используемых релизов при запуске
func (i *Installer) AutoClean() error {
	return nil
}

func (i *Installer) InstallTo(release Release, path string) error {
	if release.Architecture() != runtime.GOARCH {
		return errors.New("not valid architecture")
	}

	releaseFile, err := release.FileBinary()
	if err != nil {
		return err
	}
	defer releaseFile.Close()

	stat, err := os.Lstat(path)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return errors.New(path + " is directory, not executable")
	}

	//  раскрываем симлинк
	if stat.Mode()&os.ModeSymlink != 0 {
		path, err = filepath.EvalSymlinks(path)
		if err != nil {
			return err
		}
	}

	// 1. Проверяем чексумму
	if err := release.Validate(); err != nil {
		return err
	}

	// 2. создаем файл path.new в него копируем новый релиз
	newPath := path + ".new"
	_ = os.Remove(newPath)

	newFile, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, stat.Mode())
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, releaseFile)
	if err != nil {
		return err
	}

	// windows
	// newFile.Close()

	// 3. переименовываем текущий файл в path.old
	oldPath := path + ".old"
	_ = os.Remove(oldPath)

	err = os.Rename(path, oldPath)
	if err != nil {
		return err
	}

	// 4. копируем path.new в path
	err = os.Rename(newPath, path)
	if err != nil {
		// rollback
		return os.Rename(oldPath, path)
	}

	// 5. удалить старые релизы
	_ = os.Remove(newPath)
	_ = os.Remove(oldPath)

	return err
}

func (i *Installer) Install(release Release) error {
	execName, err := osext.Executable()
	if err != nil {
		return err
	}

	return i.InstallTo(release, execName)
}

func (i *Installer) Restart() error {
	execName, err := osext.Executable()
	if err != nil {
		return err
	}

	err = i.shutdown()
	if err != nil {
		return err
	}

	execDir := filepath.Dir(execName)

	files := make([]*os.File, 3)
	files[syscall.Stdin] = os.Stdin
	files[syscall.Stdout] = os.Stdout
	files[syscall.Stderr] = os.Stderr

	_, err = os.StartProcess(execName, []string{execName}, &os.ProcAttr{
		Dir:   execDir,
		Env:   os.Environ(),
		Files: files,
		Sys:   &syscall.SysProcAttr{},
	})

	return err
}
