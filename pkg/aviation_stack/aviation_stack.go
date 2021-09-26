package aviation_stack

import (
	"os"
   "io/ioutil"
   "net/url"
	"net/http"
   "encoding/json"
   "log"
   "github.com/davecgh/go-spew/spew"
)

const (
   FLIGHTS_URL = "http://api.aviationstack.com/v1/flights"
)


type ResponseJson struct {
   Data []Flight `json:"data"`
}

type Flight struct {
   Airline Airline `json:"airline"`
   Number Number `json:"flight"`
   Departure Departure `json:"departure"`
   Arrival Arrival `json:"arrival"`
}

type Airline struct {
   Name string `json:"name"`
}

type Number struct {
   IATA string `json:"iata"`
   ICAO string `json:"icao"`
}

type Departure struct {
   Airport string `json:"iata"`
   Terminal string `json:"terminal"`
   Gate string `json:"gate"`
}

type Arrival struct {
   Airport string `json:"iata"`
   Terminal string `json:"terminal"`
   Gate string `json:"gate"`
}


func LookupActiveFlightByICAO(flight_icao string) (*Flight, error) {

   payload := url.Values{}
	payload.Add("access_key", os.Getenv("AVIATION_STACK_API_KEY"))
	payload.Add("flight_icao", flight_icao)
	payload.Add("flight_status", "active")

   requestUrl := FLIGHTS_URL + "?" + payload.Encode()
   // spew.Dump(requestUrl)

   resp, err := http.Get(requestUrl)
   if err != nil {
      log.Println("ERROR : LookupActiveFlightByICAO()")
      log.Println(spew.Sdump(err))
      return nil, err
   }
   defer resp.Body.Close()

   body, err := ioutil.ReadAll(resp.Body)
   // spew.Dump(string(body))

   var responseJson ResponseJson
   json.Unmarshal(body, &responseJson)
   // spew.Dump(responseJson)

   if len(responseJson.Data) > 0 {
      return &responseJson.Data[0], nil
   } else {
      return nil, nil
   }

}


/*

   https://aviationstack.com/quickstart
   https://aviationstack.com/documentation

   return from /flights

   {
      "data":[
         {
            "aircraft":"None",
            "airline":{
               "iata":"AA",
               "icao":"AAL",
               "name":"American Airlines"
            },
            "arrival":{
               "actual":"None",
               "actual_runway":"None",
               "airport":"Miami International Airport",
               "baggage":"24",
               "delay":"None",
               "estimated":"2021-09-05T00:37:00+00:00",
               "estimated_runway":"None",
               "gate":"D16",
               "iata":"MIA",
               "icao":"KMIA",
               "scheduled":"2021-09-05T00:37:00+00:00",
               "terminal":"N",
               "timezone":"America/New_York"
            },
            "departure":{
               "actual":"2021-09-04T20:59:00+00:00",
               "actual_runway":"2021-09-04T20:59:00+00:00",
               "airport":"Dallas/Fort Worth International",
               "delay":14,
               "estimated":"2021-09-04T20:46:00+00:00",
               "estimated_runway":"2021-09-04T20:59:00+00:00",
               "gate":"D17",
               "iata":"DFW",
               "icao":"KDFW",
               "scheduled":"2021-09-04T20:46:00+00:00",
               "terminal":"D",
               "timezone":"America/Chicago"
            },
            "flight":{
               "codeshared":"None",
               "iata":"AA2720",
               "icao":"AAL2720",
               "number":"2720"
            },
            "flight_date":"2021-09-04",
            "flight_status":"active",
            "live":"None"
         }
      ],
      "pagination":{
         "count":1,
         "limit":100,
         "offset":0,
         "total":1
      }
   }

*/
