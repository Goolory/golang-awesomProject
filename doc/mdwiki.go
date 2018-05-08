package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

type wikiHandler struct {
	users       map[string]string
	fileHandler http.Handler
}

func (h *wikiHandler) authenticate(s string) bool {
	result := false
	if s != "" {
		ss := strings.Split(s, " ")
		if len(ss) == 2 {
			var coder = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
			info, err := coder.DecodeString(ss[1])
			if err == nil {
				user := strings.Split(string(info), ":")
				if len(user) == 2 {
					if h.users[user[0]] != "" && h.users[user[0]] == user[1] {
						result = true
					}
				}
			}
		}
	}

	return result
}

func (h *wikiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authenticate(r.Header.Get("Authorization")) {
		w.Header().Set("WWW-Authenticate", "Basic realm="+r.URL.String())
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	} else {
		h.fileHandler.ServeHTTP(w, r)
	}
}

func main() {
	h := new(wikiHandler)
	h.users = make(map[string]string)
	h.users["admin"] = "admin"
	h.fileHandler = http.FileServer(http.Dir("."))
	http.Handle("/", h)
	err := http.ListenAndServe(":5001", nil) //设置监听的端口
	if err != nil {
		fmt.Println("ListenAndServe: " + err.Error())
	}
}
