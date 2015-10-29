package crisis

import (
	"gopkg.in/pg.v3"
	"html/template"
	"log"
	"net/http"
	"time"
)

type servlet func(http.ResponseWriter, *http.Request)

const (
	staticPath = "static/"
	htmlPath   = "webcontent/html/"
)

type pageInfo struct {
	CanEdit   bool
	JSUrls    []string
	CSSUrl    string
	Types     []UnitType
	Factions  []Faction
	Divisions []Division
}

var mainPageTmpl *template.Template

func StartListening() {
	staticServer := http.FileServer(http.Dir(staticPath))
	http.Handle("/static/", http.StripPrefix("/static/", staticServer))

	ajaxHandler := GetAjaxHandlerInstance()
	http.HandleFunc("/ajax/", func(w http.ResponseWriter, r *http.Request) {
		ajaxHandler.HandleRequest(w, r)
	})

	http.HandleFunc("/staff", mainPage)
	http.HandleFunc("/view", mainPage)

	go MoveDivisions()
}

func MoveDivisions() {
	for {
		time.Sleep(10 * time.Second)
		err := GetDatabaseInstance().db.RunInTransaction(func(tx *pg.Tx) error {
			return DoUnitMovement(tx)
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	var err error

	if mainPageTmpl == nil {
		mainPageTmpl, err = template.ParseFiles(htmlPath + "mainpage.gohtml")
		maybePanic(err)
	}

	err = GetDatabaseInstance().db.RunInTransaction(func(tx *pg.Tx) error {
		authInfo, err := AuthInfoOf(tx, req)
		if err != nil {
			return err
		}

		types, err := GetUnitTypesByCrisisId(tx, authInfo.CrisisId)
		if err != nil {
			return err
		}

		facs, err := GetFactionsByCrisisId(tx, authInfo.CrisisId)
		if err != nil {
			return err
		}

		return mainPageTmpl.Execute(res, pageInfo{
			JSUrls: []string{
				"static/jquery.mousewheel.js",
				"static/buckets.min.js",
				"static/compiled.js",
			},
			CSSUrl:   "static/main.css",
			Types:    types,
			Factions: facs,
			CanEdit:  authInfo.CanEdit,
			ViewAs:   *authInfo.ViewAs,
		})
	})
	maybePanic(err)
}
