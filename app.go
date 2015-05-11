// +build !debug
package main

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/3zcurdia/badger/webhooks"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"net/http"
	"os"
)

func sha1String(s string) string {
	crypo_hash := sha1.New()
	crypo_hash.Write([]byte(s))
	return hex.EncodeToString(crypo_hash.Sum(nil))
}

func main() {
	r := render.New()
	mux := httprouter.New()
	mux.GET("/", func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		w.Write([]byte("Go to /github/:username/languages/"))
	})

	mux.GET("/github/:username/languages/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")        // Enable CORS
		w.Header().Set("Cache-Control", "public, max-age=86400")  // Exprie headers after 24 hrs
		w.Header().Set("Etag", sha1String(ps.ByName("username"))) // Calculate Etag

		res := webhooks.GithubCount(ps.ByName("username"))
		r.JSON(w, http.StatusOK, res)
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":" + os.Getenv("PORT"))
}
