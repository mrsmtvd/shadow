package ota

import (
	"crypto/md5"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/kardianos/osext"
)

type Installer struct {
}

func NewInstaller() *Installer {
	return &Installer{}
}

// очистка старых не используемых релизов при запуске
func (u *Installer) AutoClean() error {
	return nil
}

func (u *Installer) InstallTo(release Release, path string) error {
	if release.Architecture() != runtime.GOARCH {
		return errors.New("not valid architecture")
	}

	releaseFile, err := release.FileBinary()
	if err != nil {
		return err
	}
	defer releaseFile.Close()

	// TODO: проверка подписи к файлу
	// TODO: проверка что текущий файл не является релизным

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

	hasher := md5.New()
	reader := io.TeeReader(releaseFile, hasher)

	// 1. создаем файл path.new в него копируем новый релиз
	newPath := path + ".new"
	newFile, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, stat.Mode())
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, reader)
	if err != nil {
		return err
	}

	// windows
	// newFile.Close()

	// 2. Проверяем чексумму TODO: с учетом того, что брали из архива
	//if cs := hasher.Sum(nil); bytes.Compare(release.Checksum(), hasher.Sum(nil)) != 0 {
	//	return fmt.Errorf("invalid checksum want %x have %x", release.Checksum(), cs)
	//}

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
	_ = os.Remove(oldPath)

	return err
}

func (u *Installer) Install(release Release) error {
	execName, err := osext.Executable()
	if err != nil {
		return err
	}

	return u.InstallTo(release, execName)
}

func (u *Installer) Restart() error {
	execName, err := osext.Executable()
	if err != nil {
		return err
	}

	execDir := filepath.Dir(execName)

	files := []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	}

	_, err = os.StartProcess(execName, []string{execName}, &os.ProcAttr{
		Dir:   execDir,
		Env:   os.Environ(),
		Files: files,
		Sys:   &syscall.SysProcAttr{},
	})

	return err
}
