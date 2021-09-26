package dump1090

import (
	"os"
	"io/ioutil"
	"strings"
    "encoding/json"
	"log"
    "github.com/davecgh/go-spew/spew"
)


type AircraftJson struct {
	Now float32 `json:"now"`
	Aircraft []Aircraft `json:"aircraft"`
}

type Aircraft struct {
	Transponder string `json:"hex"`
	Flight string `json:"flight"`
}


func GetAircraft() ([]Aircraft, error) {
	jsonFilePath := os.Getenv("AIRCRAFT_JSON_PATH")
	jsonFile, err := os.Open(jsonFilePath)
	jsonBytes, _ := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println("ERROR : dump1090.GetAircraft()")
		log.Println(spew.Sdump(err))
		return nil, err
	}

	var aircraftJson AircraftJson
	json.Unmarshal([]byte(jsonBytes), &aircraftJson)
	for i := range aircraftJson.Aircraft {
		aircraftJson.Aircraft[i].Transponder = strings.ToUpper(aircraftJson.Aircraft[i].Transponder)
	}

	return aircraftJson.Aircraft, nil
}


/*

	https://github.com/SDRplay/dump1090/blob/master/README-json.md

	return aircraft as JSON element

	{
		"now" : 1632526795.6,
		"messages" : 18112630,
		"aircraft" : [
			{
				"hex":"a3d10f",
				"flight":"DAL2691 ",
				"alt_baro":14175,
				"alt_geom":14875,
				"gs":310.5,
				"ias":248,
				"tas":310,
				"mach":0.484,
				"track":284.9,
				"track_rate":0.56,
				"roll":8.4,
				"mag_heading":289.0,
				"baro_rate":1216,
				"geom_rate":1344,
				"squawk":"3765",
				"emergency":"none",
				"category":"A3",
				"nav_qnh":1014.4,
				"nav_altitude_mcp":16000,
				"nav_heading":289.7,
				"lat":26.034432,
				"lon":-80.356275,
				"nic":8,
				"rc":186,
				"seen_pos":1.8,
				"version":2,
				"nic_baro":1,
				"nac_p":9,
				"nac_v":1,
				"sil":3,
				"sil_type":"perhour",
				"gva":2,
				"sda":2,
				"mlat":[],
				"tisb":[],
				"messages":1953,
				"seen":1.5,
				"rssi":-19.1
			},
			...
		]
	}

*/
