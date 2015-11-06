package main

import (
	"fmt"
    "html/template"
    "log"
	"strings"
    "net/http"
	//"strconv"
	"time"
	"io/ioutil"
	
	"main/locdata"

	"github.com/golang/protobuf/proto"
    "github.com/gorilla/mux"
)

var err error

type Page struct {
    Title string
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", Index)
	router.HandleFunc("/locations", GetAllLocations)
	router.HandleFunc("/locations/{latitude1}/{longitude1}/{latitude2}/{longitude2}", GetLocations)
	
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("/home")))
	
	http.Handle("/", router)
    log.Fatal(http.ListenAndServe(":80", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("/home/tpl/index.html")
	p := &Page{Title: "Homepage"}
	t.Execute(w, p)
}

func GetAllLocations(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	
	requestStr := fmt.Sprintf("http://heatmap-endpoint.salscode.com/locations/")
	resp, err := http.Get(requestStr)
	checkErr(err)
	
	locResp, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	
	locations := &locdata.LocDisplayData{}
	
	err = proto.Unmarshal(locResp, locations)
	checkErr(err)
	
	printLocations(w, locations)
	
	t1 := time.Now()
    fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
}

func GetLocations(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	
	vars := mux.Vars(r)
	
	requestStr := fmt.Sprintf("http://heatmap-endpoint.salscode.com/locations/%s/%s/%s/%s/", vars["latitude1"], vars["longitude1"], vars["latitude2"], vars["longitude2"])
	resp, err := http.Get(requestStr)
	checkErr(err)
	
	locResp, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	
	locations := &locdata.LocDisplayData{}
	
	err = proto.Unmarshal(locResp, locations)
	checkErr(err)
	
	printLocations(w, locations)
	
	t1 := time.Now()
    fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
}

func printLocations(w http.ResponseWriter, locations *locdata.LocDisplayData) {
	var locStr string = "";
	for _, each := range locations.Lldata {
		str := strings.Replace(fmt.Sprintf(*locations.Coordformat, *each.Latitude, *each.Longitude, *each.Intensity), " ", "", -1)
        locStr += str
	}
	
	fmt.Fprintf(w, "// %d Locations\n", *locations.Count)
    fmt.Fprintf(w, "var locations = [%s];", locStr)
}

func checkErr(inerr error) {
    if inerr != nil {
        panic(inerr)
    }
}