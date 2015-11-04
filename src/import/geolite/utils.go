package geolite

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

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
