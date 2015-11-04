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

// Used for the newer GeoLite2 db, http://dev.maxmind.com/geoip/geoip2/geolite2/

func Import2() error {
    var locDataList map[string]int32
    locDataList = make(map[string]int32)
    
    // Read Blocks
    fmt.Println("Read Blocks")
    readBlocks2(locDataList)
    
    // Open DB
    fmt.Println("Insert into Database")
    fillDatabase2(locDataList)
    
    return nil;
}

func readBlocks2(locDataList map[string]int32) {
    blockFile, err := os.Open("data/GeoLite2-City-Blocks-IPv4.csv")
    checkErr(err)
    
    blockFileReader := csv.NewReader(blockFile)
    blockFileReader.FieldsPerRecord = -1
    
    blockCSVData, err := blockFileReader.ReadAll()
    checkErr(err)
    
    blockFile.Close()
    
    for i, each := range blockCSVData {
        if i < 9 {
            continue
        }
        
        networkArr := strings.Split(each[0], "/")
        if len(networkArr) < 2 {
            continue;
        }
        
        network, err := strconv.ParseInt(networkArr[len(networkArr) - 1], 10, 64)
        checkErr(err)
        if network > 32 || network < 1 {
            continue;
        }
        
        ipcount := int32(2^(32-network))
        
        latitude, err := strconv.ParseFloat(each[7], 64)
        checkErr(err)
        
		longitude, err := strconv.ParseFloat(each[8], 64)
        checkErr(err)
        
        key := fmt.Sprintf("%8.4f,%8.4f", latitude, longitude)
        
        _, ok := locDataList[key]
        if !ok {
            locDataList[key] = ipcount
        } else {
            locDataList[key] += ipcount
        }
    }
}

func fillDatabase2(condensedList map[string]int32) {
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