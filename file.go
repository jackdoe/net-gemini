package gemini

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

//--------------------------------------
type byDirs []os.FileInfo

func (s byDirs) Len() int {
    return len(s)
}

func (s byDirs) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func (s byDirs) Less(i, j int) bool {
	if s[i].IsDir() {
		if s[j].IsDir() {
			return s[i].Name() < s[j].Name()
		}
		return true
	}
	if s[j].IsDir() {
		if s[i].IsDir() {
			return s[i].Name() < s[j].Name()
		}
		return false
	}
	return s[i].Name() < s[j].Name()
}
//--------------------------------------

func brief(s string) (r string) {
	if len(s) < 58 {
		return s
	} else {
		return s[0:58] + " ..."
	}
}


type fileHandler struct {
	root string
}

func FileServer(root string) Handler {
	return &fileHandler{root: filepath.Clean(root)}
}

func ServeFilePath(p string, w *Response, r *Request) {
	s, err := os.Stat(p)
	if err != nil {
		w.SetStatus(StatusNotFound, "File Not Found!")
		return
	}
	if !allowed(s) {
		w.SetStatus(StatusGone, "Forbidden!")
		return
	}
	if s.IsDir() {
		files, err := ioutil.ReadDir(p)
		if err != nil {
			w.SetStatus(StatusTemporaryFailure, "Error reading directory!")
			return
		}

		for _, f := range files {
			if f.Name() == "index.gmi" {
				p = path.Join(p, "index.gmi")
				goto FILE
			}
		}
		w.SetStatus(StatusSuccess, "text/gemini")
		w.Write([]byte(fmt.Sprintf("# Listing %s\n\n", p)))

		sort.Sort(byDirs(files))
		//sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

		rpath := r.URL.String()
		w.Write([]byte(fmt.Sprintf("=> %s ..\n", rpath[0:strings.LastIndex(rpath,"/")])))
		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}

			if !allowed(file) {
				continue
			}

			if file.IsDir() {
				w.Write([]byte(fmt.Sprintf("=> %s % 80s %*v\n", filepath.Clean(path.Join(r.URL.Path, file.Name())), brief("⍠: "+file.Name()), 82-len(brief("⍠: "+file.Name())), "["+file.ModTime().Format(time.Stamp)+"]")))
			} else {
				w.Write([]byte(fmt.Sprintf("=> %s % 80s %*v\n", filepath.Clean(path.Join(r.URL.Path, file.Name())), brief(file.Name()), 80-len(brief(file.Name())), "["+file.ModTime().Format(time.Stamp)+"]")))
			}
		}
		return
	}
FILE:
	ext := filepath.Ext(p)
	var mimeType string
	if ext == ".gmi" {
		mimeType = "text/gemini"
	} else {
		mimeType = mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "octet/stream"
		}
	}

	f, err := os.OpenFile(p, os.O_RDONLY, 0600)
	if err != nil {
		w.SetStatus(StatusTemporaryFailure, "Error reading file!")
		return
	}
	defer f.Close()

	w.SetStatus(StatusSuccess, mimeType)
	_, err = io.Copy(w, f)
	if err != nil {
		// .. remote closed the connection, nothing we can do besides log
		// or io error, but status is already sent, everything is broken!
		w.SetStatus(StatusTemporaryFailure, "IO error!")
	}
}
func (fh *fileHandler) ServeGemini(w *Response, r *Request) {
	p := filepath.Clean(path.Join(fh.root, r.URL.Path))
	if !strings.HasPrefix(p, fh.root) {
		w.SetStatus(StatusTemporaryFailure, "Path not in scope!")
		return
	}

	ServeFilePath(p, w, r)
}

func allowed(fi os.FileInfo) bool {
	return uint64(fi.Mode().Perm())&0444 == 0444
}

// almost copy pasta form https://tildegit.org/solderpunk/molly-brown
func generateDirectoryListing(path string, out io.Writer) error {

	return nil
}
