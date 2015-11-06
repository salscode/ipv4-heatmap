package locdata

import (
	"fmt"
    "log"
    "net/http"
	"strconv"
	//"time"
	"database/sql"

	"github.com/golang/protobuf/proto"
    "github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

type Page struct {
    Title string
}

func Init() {
	db, err = sql.Open("mysql", "geodata:hF9yaD5XNTnDXVwf@/geodata?charset=utf8")
    checkErr(err)

	router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", Index)
	router.HandleFunc("/locations", GetAllLocations)
	router.HandleFunc("/locations/{latitude1}/{longitude1}/{latitude2}/{longitude2}", GetLocations)
	
	http.Handle("/", router)
    log.Fatal(http.ListenAndServe(":80", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Use the API\n")
}

func GetAllLocations(w http.ResponseWriter, r *http.Request) {
	//t0 := time.Now()
	
	totalQuery := fmt.Sprintf("SELECT COUNT(tmploc.intensity) count, MAX(tmploc.intensity) maxintensity FROM (SELECT LOG(SUM(ipcount)) intensity FROM locations GROUP BY TRUNCATE(latitude, 0), TRUNCATE(longitude, 0)) tmploc")
	count, maxintensity := processTotalQuery(totalQuery)
	
	rows, err := db.Query("SELECT TRUNCATE(latitude, 0) lattrunc, TRUNCATE(longitude, 0) lontrunc, LOG(SUM(ipcount)) ipcount FROM locations GROUP BY lattrunc, lontrunc")
    checkErr(err)
	
	printLocations(w, rows, count, maxintensity, "[%4.0f, %4.0f, %2.1f],\n")
	rows.Close()
	
	//t1 := time.Now()
    //fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
}

func GetLocations(w http.ResponseWriter, r *http.Request) {
	//t0 := time.Now()
	
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
	coordFormat := "[%4.0f, %4.0f, %2.1f],"
	
	fmt.Printf("Offset = %f ", offset)
	
	if offset < 30 {
		trunc = 4
		coordFormat = "[%8.4f, %8.4f, %2.1f],"
	} else if offset < 200 {
		trunc = 1
		coordFormat = "[%5.1f, %5.1f, %2.1f],"
	}
	
	totalQuery := fmt.Sprintf("SELECT COUNT(tmploc.intensity) count, MAX(tmploc.intensity) maxintensity FROM (SELECT LOG(SUM(ipcount)) intensity FROM locations WHERE latitude >= %s AND latitude <= %s AND longitude >= %s AND longitude <= %s GROUP BY TRUNCATE(latitude, %d), TRUNCATE(longitude, %d)) tmploc", lowLat, highLat, lowLon, highLon, trunc, trunc)
	count, maxintensity := processTotalQuery(totalQuery)
	
	query := fmt.Sprintf("SELECT TRUNCATE(latitude, %d) lattrunc, TRUNCATE(longitude, %d) lontrunc, LOG(SUM(ipcount)) ipcount FROM locations WHERE latitude >= %s AND latitude <= %s AND longitude >= %s AND longitude <= %s GROUP BY lattrunc, lontrunc", trunc, trunc, lowLat, highLat, lowLon, highLon)
	
	rows, err := db.Query(query)
    checkErr(err)
	
	printLocations(w, rows, count, maxintensity, coordFormat)
	rows.Close()
	
	//t1 := time.Now()
    //fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
}

func printLocations(w http.ResponseWriter, rows *sql.Rows, count int32, maxintensity float32, coordFormat string) {
	locations := &LocDisplayData{Count: &count, Coordformat: &coordFormat}
	
	for rows.Next() {
		var latitude float32
        var longitude float32
        var intensity float32
		
        err = rows.Scan(&latitude, &longitude, &intensity)
        checkErr(err)
		
		intensity = intensity/maxintensity
		if (intensity > 1) {
			intensity = 1
		}
		
		locdatall := &LocDisplayData_LatLon{Latitude: &latitude, Longitude: &longitude, Intensity: &intensity}
		locations.Lldata = append(locations.Lldata, locdatall)
	}
	
	data, err := proto.Marshal(locations)
	checkErr(err)
	
	w.Write(data)
}

func processTotalQuery(query string) (count int32, maxintensity float32) {
	totalData, err := db.Query(query)
    checkErr(err)
	
	totalData.Next()
	err = totalData.Scan(&count, &maxintensity)
	checkErr(err)
	
	// Reduce max intensity for better display.
	maxintensity = maxintensity * 0.7
	
	totalData.Close()
	
	return count, maxintensity
}

func checkErr(inerr error) {
    if inerr != nil {
        panic(inerr)
    }
}