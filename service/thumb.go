package service

import (
	"bytes"
	"errors"
	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func BuildThumbPath(path string, width int, height int) string {
	dir := filepath.Dir(path)
	name := filepath.Base(path)
	ext := filepath.Ext(path)
	if ext != "" {
		name = name[:len(name)-len(ext)]
	}
	return filepath.Join(dir, name+"@w"+strconv.Itoa(width)+"_h"+strconv.Itoa(height)+ext)
}

func IsThumb(path string) (string, bool) {
	n := filepath.Base(path)
	yes := strings.Contains(n, "@w") && strings.Contains(n, "_h")
	if yes {
		base := filepath.Base(path)
		sl := strings.Split(base, "@")
		filename := sl[0] + filepath.Ext(base)
		return filepath.Join(filepath.Dir(path), filename), yes
	} else {
		return path, yes
	}
}

func CreateThumb(path string) ([]byte, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".jpg" && ext != ".jpeg" &&
		ext != ".png" && ext != ".bmp" && ext != ".gif" {
		return nil, errors.New("Unsupport file format " + ext)
	}

	// example path: ../file/o/2016-05-19/91bd48ce-1dae-11e6-9d7f-408d5cdf2c91@w360_h100.jpg
	// take the string after "@" (w360_h100.jpg) to extract the size (360*270)
	sl := strings.Split(path, "@")
	sizeString := sl[1]
	ws := sizeString[strings.Index(sizeString, "@w")+2 : strings.Index(sizeString, "_h")]
	var hs string
	if strings.Index(sizeString, ".") == -1 {
		hs = sizeString[strings.Index(sizeString, "_h")+2:]
	} else {
		hs = sizeString[strings.Index(sizeString, "_h")+2 : strings.Index(sizeString, ".")]
	}
	w, _ := strconv.Atoi(ws)
	h, _ := strconv.Atoi(hs)

	srcPath := sl[0] + ext
	src, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	i, _, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	m := imaging.Fill(i, w, h, imaging.Center, imaging.Lanczos)

	var buffer bytes.Buffer
	if ext == ".jpg" || ext == ".jpeg" {
		if err := jpeg.Encode(&buffer, m, nil); err != nil {
			return nil, err
		}
	} else if ext == ".png" {
		if err := png.Encode(&buffer, m); err != nil {
			return nil, err
		}
	} else if ext == ".bmp" {
		if err := bmp.Encode(&buffer, m); err != nil {
			return nil, err
		}
	} else if ext == ".gif" {
		if err := gif.Encode(&buffer, m, nil); err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}
