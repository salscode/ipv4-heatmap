package main

import (
    "os"
    "fmt"
    "encoding/csv"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
    "database/sql"
)

type LocData struct {
    latitude float64
    longitude float64
    ipcount int32
}

func (l *LocData) SetLatLon(lat float64, lon float64) {
    l.latitude = lat;
    l.longitude = lon;
}

func (l *LocData) IncrementIPCount(count int32) {
    if count > 0 {
        l.ipcount += count;
    }
}

func main() {
    var locDataList map[int64]*LocData
    locDataList = make(map[int64]*LocData)
    
    // Read Blocks
    fmt.Println("Read Blocks")
    readBlocks(locDataList)
    
    // Read Locations
    fmt.Println("Read Locations")
	readLocations(locDataList)
    
    // Open DB
    fmt.Println("Insert into Database")
    
    db, err := sql.Open("mysql", "geodata:hF9yaD5XNTnDXVwf@/geodata?charset=utf8")
    checkErr(err)
    
    insertStmt, err := db.Prepare("INSERT INTO locations (locid, latitude, longitude, ipcount) VALUES(?, ?, ?, ?)")
    checkErr(err)
    
    _, err = db.Query("TRUNCATE TABLE locations")
    checkErr(err)
    
    transaction, err := db.Begin()
    checkErr(err)
    
    for locid, loc := range locDataList {
        if loc.ipcount <= 0 {
            continue;
        }
        
        _, err = insertStmt.Exec(locid, loc.latitude, loc.longitude, loc.ipcount)
        checkErr(err)
    }
    
    err = transaction.Commit()
    checkErr(err)
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func readLocations(locDataList map[int64]*LocData) {
    locFile, err := os.Open("GeoLiteCity-Location.csv")
    checkErr(err)
    
    locFileReader := csv.NewReader(locFile)
    locFileReader.FieldsPerRecord = -1
    
    locCSVData, err := locFileReader.ReadAll()
    checkErr(err)
    
    locFile.Close()

    for i, each := range locCSVData {
        if i <= 2 {
            continue
        }
        
        locid, err := strconv.ParseInt(each[0], 10, 64)
        checkErr(err)
        
		latitude, err := strconv.ParseFloat(each[5], 64)
        checkErr(err)
        
		longitude, err := strconv.ParseFloat(each[6], 64)
        checkErr(err)

        _, ok := locDataList[locid]
        if ok && locDataList[locid].ipcount > 0 {
            locDataList[locid].SetLatLon(latitude, longitude)
        }
    }
}

func readBlocks(locDataList map[int64]*LocData) {
    blockFile, err := os.Open("GeoLiteCity-Blocks.csv")
    checkErr(err)
    
    blockFileReader := csv.NewReader(blockFile)
    blockFileReader.FieldsPerRecord = -1
    
    blockCSVData, err := blockFileReader.ReadAll()
    checkErr(err)
    
    blockFile.Close()
    
    for i, each := range blockCSVData {
        if i <= 2 {
            continue
        }
        
        ipstart, err := strconv.ParseInt(each[0], 10, 64)
        checkErr(err)
        
        ipend, err := strconv.ParseInt(each[1], 10, 64)
        checkErr(err)
        
        locid, err := strconv.ParseInt(each[2], 10, 64)
        checkErr(err)
        
		ipcount := int32(ipend - ipstart);
		if ipcount <= 0 {
			continue;
        }
        
        _, ok := locDataList[locid]
        if !ok {
            locDataList[locid] = new(LocData)
        }
        
        locDataList[locid].IncrementIPCount(ipcount)
    }
}