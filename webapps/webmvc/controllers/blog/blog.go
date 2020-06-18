package blogController

import (
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter,r *http.Request,tmplsPath []string){

	files := append(tmplsPath, "views/blog/blog.html")
	tmpls, _ := template.ParseFiles(files...)
	tmpls.ExecuteTemplate(w,"masterpage",nil)
}