package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/goLangCoin/blockchain"
)

var templates *template.Template

const templateDir string = "explorer/templates/"

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	// template.Must => err를 체크해줌
	// tmpl := template.Must(template.ParseFiles("templates/pages/home.gohtml"))
	data := homeData{"Home", nil}
	templates.ExecuteTemplate(rw, "home.gohtml", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		// Form 값을 채워줌
		blockchain.Blockchain().AddBlock()
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(aPort int) {
	handler := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	// 이 핸들러가 동작하기 전에 gohtml은 rendering이 되어있어야 한다.
	handler.HandleFunc("/add", add)
	handler.HandleFunc("/", home)
	fmt.Printf("Listening on http://localhost:%d\n", aPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", aPort), handler))
}
