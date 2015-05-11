// +build !debug
package main

import (
	"github.com/3zcurdia/merithub/webhooks"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"net/http"
	"os"
)

func main() {
	r := render.New()
	mux := httprouter.New()
	mux.GET("/", func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		w.Write([]byte("Go to /languages/:username"))
	})

	mux.GET("/languages/:username", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		res := webhooks.GithubCount(ps.ByName("username"))
		r.JSON(w, http.StatusOK, res)
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":" + os.Getenv("PORT"))
}
