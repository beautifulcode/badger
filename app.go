// +build !debug
package main

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/3zcurdia/merithub/webhooks"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"net/http"
	"os"
)

func getSha1(s []byte) []byte {
	h := sha1.New()
	h.Write(s)
	return h.Sum(nil)
}

func main() {
	r := render.New()
	mux := httprouter.New()
	mux.GET("/", func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		w.Write([]byte("Go to /github/:username/languages/"))
	})

	mux.GET("/github/:username/languages/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		etag := getSha1([]byte(ps.ByName("username")))
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Etag", hex.EncodeToString(etag))
		res := webhooks.GithubCount(ps.ByName("username"))
		r.JSON(w, http.StatusOK, res)
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":" + os.Getenv("PORT"))
}
