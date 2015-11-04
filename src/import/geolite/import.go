package geolite

import (
    "os"
    "fmt"
    "encoding/csv"
	"strconv"
    "strings"

	_ "github.com/go-sql-driver/mysql"
    "database/sql"
)

// Used for the legacy GeoLite db, http://dev.maxmind.com/geoip/legacy/geolite/

func Import1() error {
    var locDataList map[int64]*LocData
    locDataList = make(map[int64]*LocData)
    
    // Read Blocks
    fmt.Println("Read Blocks")
    readBlocks1(locDataList)
    
    // Read Locations
    fmt.Println("Read Locations")
	readLocations1(locDataList)
    
    // Delete Duplicates
    fmt.Println("Delete Duplicates")
    condensedList := deleteDuplicates1(locDataList)
    
    // Open DB
    fmt.Println("Insert into Database")
    fillDatabase1(condensedList)
    
    return nil;
}

func readLocations1(locDataList map[int64]*LocData) {
    locFile, err := os.Open("data/GeoLiteCity-Location.csv")
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

func readBlocks1(locDataList map[int64]*LocData) {
    blockFile, err := os.Open("data/GeoLiteCity-Blocks.csv")
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

func deleteDuplicates1(locDataList map[int64]*LocData) map[string]int32 {
    condensedList := make(map[string]int32)
    
    for _, each := range locDataList {
        key := fmt.Sprintf("%8.4f,%8.4f", each.latitude, each.longitude)
        
        _, ok := condensedList[key]
        if !ok {
            condensedList[key] = each.ipcount
        } else {
            condensedList[key] += each.ipcount
        }
    }
    
    return condensedList
}

func fillDatabase1(condensedList map[string]int32) {
    db, err := sql.Open("mysql", "geodata:hF9yaD5XNTnDXVwf@/geodata?charset=utf8")
    checkErr(err)
    
    insertStmt, err := db.Prepare("INSERT INTO locations (latitude, longitude, ipcount) VALUES(?, ?, ?)")
    checkErr(err)
    
    _, err = db.Query("TRUNCATE TABLE locations")
    checkErr(err)
    
    transaction, err := db.Begin()
    checkErr(err)
    
    for ipstr, ipcnt := range condensedList {
        if ipcnt <= 0 {
            continue;
        }
        
        s := strings.Split(ipstr, ",")
        if len(s) != 2 {
            continue;
        }
        
        _, err = insertStmt.Exec(s[0], s[1], ipcnt)
        checkErr(err)
    }
    
    err = transaction.Commit()
    checkErr(err)
}