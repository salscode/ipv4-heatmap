package main

import (
    "os"
	"fmt"
    "html/template"
    "log"
    "net/http"
	"strconv"
	"time"
	"database/sql"
	
	"import/geolite"

    "github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

type Page struct {
    Title string
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "import" {
			fmt.Println("Running import 1.")
			geolite.Import1()
			os.Exit(0)
		}
		
		if os.Args[1] == "import2" {
			fmt.Println("Running import 2.")
			geolite.Import2()
			os.Exit(0)
		}
	}
	
	db, err = sql.Open("mysql", "geodata:hF9yaD5XNTnDXVwf@/geodata?charset=utf8")
    checkErr(err)

	router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", Index)
	router.HandleFunc("/locations", GetAllLocations)
	router.HandleFunc("/locations/{latitude1}/{longitude1}/{latitude2}/{longitude2}", GetLocations)
	
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
	
	http.Handle("/", router)
    log.Fatal(http.ListenAndServe(":80", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tpl/index.html")
	p := &Page{Title: "Homepage"}
	t.Execute(w, p)
}

func GetAllLocations(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	
	rows, err := db.Query("SELECT TRUNCATE(latitude, 0) lattrunc, TRUNCATE(longitude, 0) lontrunc, LOG(SUM(ipcount)) ipcount FROM locations GROUP BY lattrunc, lontrunc")
    checkErr(err)
	
	PrintLocations(w, rows, "[%4.0f, %4.0f, %3.1f],\n")
	
	t1 := time.Now()
    fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
}

func GetLocations(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	
	vars := mux.Vars(r)
	latitude1, err := strconv.ParseFloat(vars["latitude1"], 32)
	longitude1, err := strconv.ParseFloat(vars["longitude1"], 32)
	latitude2, err := strconv.ParseFloat(vars["latitude2"], 32)
	longitude2, err := strconv.ParseFloat(vars["longitude2"], 32)
	
	var lowLat, highLat, lowLon, highLon string = "", "", "", ""
	var offset float64 = 0
	
	if latitude1 > latitude2 {
		lowLat = vars["latitude2"]
		highLat = vars["latitude1"]
		offset += (latitude1 - latitude2)
	} else {
		lowLat = vars["latitude1"]
		highLat = vars["latitude2"]
		offset += (latitude2 - latitude1)
	}
	
	if longitude1 > longitude2 {
		lowLon = vars["longitude2"]
		highLon = vars["longitude1"]
		offset += (longitude1 - longitude2)
	} else {
		lowLon = vars["longitude1"]
		highLon = vars["longitude2"]
		offset += (longitude2 - longitude1)
	}
	
	// Calculate truncation
	// No substantial reduction in result count between 4 decimal places and 3 or 2.
	trunc := 0
	coordFormat := "[%4.0f, %4.0f, %3.1f],\n"
	
	fmt.Printf("Offset = %f ", offset)
	
	if offset < 30 {
		trunc = 4
		coordFormat = "[%8.4f, %8.4f, %3.1f],\n"
	} else if offset < 200 {
		trunc = 1
		coordFormat = "[%5.1f, %5.1f, %3.1f],\n"
	}
	
	query := fmt.Sprintf("SELECT TRUNCATE(latitude, %d) lattrunc, TRUNCATE(longitude, %d) lontrunc, LOG(SUM(ipcount)) ipcount FROM locations WHERE latitude >= %s AND latitude <= %s AND longitude >= %s AND longitude <= %s GROUP BY lattrunc, lontrunc", trunc, trunc, lowLat, highLat, lowLon, highLon)
	
	rows, err := db.Query(query)
    checkErr(err)
	
	PrintLocations(w, rows, coordFormat)
	
	t1 := time.Now()
    fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
}

func PrintLocations(w http.ResponseWriter, rows *sql.Rows, coordFormat string) {
	count := 0
	var locations string = "";
	for rows.Next() {
        var latitude float32
        var longitude float32
        var ipcount float64
		
        err = rows.Scan(&latitude, &longitude, &ipcount)
        checkErr(err)
		
        locations += fmt.Sprintf(coordFormat, latitude, longitude, ipcount)
		
		count++
	}
	
	fmt.Fprintf(w, "// %d Locations\n", count)
    fmt.Fprintf(w, "var locations = [\n%s];", locations)
}

func checkErr(inerr error) {
    if inerr != nil {
        panic(inerr)
    }
}