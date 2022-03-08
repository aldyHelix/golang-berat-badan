package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Berat struct {
	Id        int
	Tanggal   string
	Max       int
	Min       int
	Perbedaan int
}

type Average struct {
	ListBerat        []Berat
	AverageMax       float32
	AverageMin       float32
	AveragePerbedaan float32
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "berat_badan"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM berat_badan ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	berat := Berat{}
	avg := Average{}
	res := []Berat{}
	sumMax := 0
	sumMin := 0
	sumPerbedaan := 0

	for selDB.Next() {
		var id, max, min int
		var tanggal string
		err = selDB.Scan(&id, &tanggal, &max, &min)
		if err != nil {
			panic(err.Error())
		}

		berat.Id = id
		berat.Tanggal = tanggal
		berat.Max = max
		berat.Min = min
		berat.Perbedaan = max - min
		sumMax += max
		sumMin += min
		sumPerbedaan += (max - min)
		res = append(res, berat)
	}

	sumResult := len(res)
	avg.ListBerat = res
	avg.AverageMax = (float32(sumMax) / float32(sumResult))
	avg.AverageMin = (float32(sumMin) / float32(sumResult))
	avg.AveragePerbedaan = (float32(sumPerbedaan) / float32(sumResult))

	tmpl.ExecuteTemplate(w, "Index", avg)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM berat_badan WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	berat := Berat{}

	for selDB.Next() {
		var id, min, max int
		var tanggal string
		err = selDB.Scan(&id, &tanggal, &max, &min)
		if err != nil {
			panic(err.Error())
		}

		berat.Id = id
		berat.Tanggal = tanggal
		berat.Max = max
		berat.Min = min
		berat.Perbedaan = max - min
	}
	tmpl.ExecuteTemplate(w, "Show", berat)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM berat_badan WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	berat := Berat{}
	for selDB.Next() {
		var id, min, max int
		var tanggal string
		err = selDB.Scan(&id, &tanggal, &min, &max)
		if err != nil {
			panic(err.Error())
		}

		berat.Id = id
		berat.Tanggal = tanggal
		berat.Max = max
		berat.Min = min
	}
	tmpl.ExecuteTemplate(w, "Edit", berat)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		tanggal := r.FormValue("tanggal")
		min := r.FormValue("min")
		max := r.FormValue("max")
		insForm, err := db.Prepare("INSERT INTO berat_badan(tanggal, max, min) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(tanggal, max, min)
		log.Println("INSERT: Tanggal: " + tanggal + "| MIN: " + min + " | MAX: " + max)
	}

	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		tanggal := r.FormValue("tanggal")
		min := r.FormValue("min")
		max := r.FormValue("max")
		id := r.FormValue("id")
		insForm, err := db.Prepare("UPDATE berat_badan SET tanggal=?, min=?, max=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(tanggal, min, max, id)
		log.Println("UPDATE: Tanggal: "+tanggal+"| MIN: %s | MAX: %s", min, max)
	}

	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	berat := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM berat_badan WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(berat)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
