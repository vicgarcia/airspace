package faa_registry

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "encoding/json"
    "log"
    "github.com/davecgh/go-spew/spew"
)

const (
    DB_FILE_PATH = "./faa_registry.db"
)


type Registration struct {
    RegistrationNumber string `json:"registration_number"`
    Transponder string `json:"transponder_code_hex"`
    Aircraft Aircraft `json:"aircraft"`
    Engine Engine `json:"engine"`
    Owner Owner `json:"registrant"`
}

type Aircraft struct {
    Manufacturer string `json:"manufacturer"`
    Model string `json:"model"`
    Type string `json:"type"`
    Engines int `json:"number_of_engines"`
    Seats int `json:"number_of_seats"`
}

type Engine struct {
    Manufacturer string `json:"manufacturer"`
    Model string `json:"model"`
    Type string `json:"type"`
}

type Owner struct {
    Name string `json:"name"`
    City string `json:"city"`
    State string `json:"state"`
}


type Registry struct {
    db *sql.DB
}

func (r *Registry) LookupByTransponder(transponder string) (*Registration, error) {
    var data string;
    row := r.db.QueryRow("select data from aircrafts where transponder = ?", transponder)
    switch err := row.Scan(&data); err {
    case sql.ErrNoRows:
        return nil, err
    case nil:
        var registration Registration
        json.Unmarshal([]byte(data), &registration)
        return &registration, nil
    default:
        log.Println("ERROR : LookupByTransponder()")
        log.Println(spew.Sdump(err))
        return nil, err
    }

}

func (r *Registry) LookupByRegistration(registration string) (*Registration, error) {
    var data string;
    row := r.db.QueryRow("select data from aircrafts where registration = ?", registration)
    switch err := row.Scan(&data); err {
    case nil:
        var registration Registration
        json.Unmarshal([]byte(data), &registration)
        return &registration, nil
    case sql.ErrNoRows:
        return nil, err
    default:
        log.Println("ERROR : LookupByRegistration()")
        log.Println(spew.Sdump(err))
        return nil, err
    }
}

func Connect() (*Registry, error) {
    db, err := sql.Open("sqlite3", DB_FILE_PATH)
    if err != nil {
        log.Println("ERROR : faa_registry.Connect()")
        log.Println(spew.Sdump(err))
        return nil, err
    }
    registry := Registry{db}
    return &registry, nil
}


/*

    return from database will be json object of aircraft data

    {
        "aircraft":{
            "category":"Land",
            "certification":"Type Certificated",
            "code":"3940032",
            "cruising_speed_mph":"None",
            "engine_type":"Turbo-fan",
            "manufacturer":"AIRBUS",
            "model":"A321-231",
            "number_of_engines":2,
            "number_of_seats":379,
            "type":"Fixed wing multi engine",
            "weight_category":"20,000 and over."
        },
        "aircraft_type":"Fixed wing multi engine",
        "airworthiness_date":"2018-03-12",
        "certificate_issue_date":"2018-03-12",
        "certification":{
            "classification":"Standard",
            "operations":[
                "Transport"
            ],
            "subclassifications":"None"
        },
        "engine":{
            "code":"34611",
            "manufacturer":"IAE",
            "model":"V2533-A5",
            "power_hp":"None",
            "thrust_lbf":31600,
            "type":"Turbo-fan"
        },
        "engine_type":"Turbo-fan",
        "expiration_date":"2024-03-31",
        "fractional_ownership":false,
        "kit_manufacturer":"",
        "kit_model":"",
        "last_action_date":"2021-02-05",
        "manufacturing_year":2018,
        "other_names":"None",
        "registrant":{
            "city":"MIRAMAR",
            "country":"US",
            "county":"011",
            "name":"SPIRIT AIRLINES INC",
            "region":"Southern",
            "state":"FL",
            "street_1":"2800 EXECUTIVE WAY",
            "street_2":"",
            "type":"Corporation",
            "zip_code":"33025-6542"
        },
        "registration_number":"N685NK",
        "serial_number":"8115",
        "source":"FAA",
        "status":"Valid Registration",
        "transponder_code":"52213400",
        "transponder_code_hex":"A91700",
        "unique_regulatory_id":"01271201"
    }

*/
