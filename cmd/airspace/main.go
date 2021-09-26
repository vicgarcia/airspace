package main

import (
    "os"
    "fmt"
    "strings"
    "time"
    "github.com/joho/godotenv"
    "log"
    "github.com/davecgh/go-spew/spew"
    "database/sql"
    "github.com/mattn/go-tty"
    "github.com/vicgarcia/airspace/pkg/console"
    "github.com/vicgarcia/airspace/pkg/dump1090"
    "github.com/vicgarcia/airspace/pkg/faa_registry"
    "github.com/vicgarcia/airspace/pkg/aviation_stack"
)


func makeAircraftLineOutput(aircraft dump1090.Aircraft, registration *faa_registry.Registration) string {
    output := fmt.Sprintf("%s %s %s",
        aircraft.Transponder,
        registration.Aircraft.Manufacturer,
        registration.Aircraft.Model,
    )
    if len(aircraft.Flight) > 0 && aircraft.Flight != registration.RegistrationNumber {
        output += fmt.Sprintf(" FLIGHT %s", aircraft.Flight)
    }
    return output
}

func makeRegistrationOutput(registration *faa_registry.Registration) string {
    var output string
    output += fmt.Sprintf("%s %s %s \n",
        registration.Aircraft.Manufacturer,
        registration.Aircraft.Model,
        registration.AirworthinessDate,
    )
    if len(registration.Engine.Manufacturer) > 0 && len(registration.Engine.Model) > 0 {
        output += fmt.Sprintf("%dx %s %s \n",
            registration.Aircraft.Engines,
            registration.Engine.Manufacturer,
            registration.Engine.Model,
        )
    }
    output += fmt.Sprintf("%s %d SEATS\n",
        strings.ToUpper(registration.Aircraft.Type),
        registration.Aircraft.Seats,
    )
    output += fmt.Sprintf("REG# %s TRNSP %s \n",
        registration.RegistrationNumber,
        registration.Transponder,
    )
    output += fmt.Sprintf("%s %s %s \n",
        registration.Owner.Name,
        registration.Owner.City,
        registration.Owner.State,
    )
    return output
}

func makeFlightOutput(flight *aviation_stack.Flight) string {
    var output string
    output += fmt.Sprintf("%s %s \n",
        flight.Airline.Name,
        flight.Number.IATA,
    )
    output += fmt.Sprintf("%s (%s) -> %s (%s) \n",
        flight.Departure.Airport,
        flight.Departure.Gate,
        flight.Arrival.Airport,
        flight.Arrival.Gate,
    )
    output += fmt.Sprintf("https://flightaware.com/live/flight/%s \n",
        flight.Number.ICAO,
    )
    return output
}

func Ping() {
    faa, err := faa_registry.Connect()
    if err != nil {
        console.Renderln("error connecting to the faa registry")
        return
    }
    aircrafts, err := dump1090.GetAircraft()
    if err != nil {
        console.Renderln("error connecting to ther ADS-B receiver")
        return
    }
    for _, aircraft := range aircrafts {
        registration, _ := faa.LookupByTransponder(aircraft.Transponder)
        if registration != nil {
            output := makeAircraftLineOutput(aircraft, registration)
            console.Renderln(output)
        }
    }
}

func Live(tty *tty.TTY) {
    faa, err := faa_registry.Connect()
    if err != nil {
        console.Renderln("error connecting to the faa registry")
        return
    }
    quit := make(chan bool)
    go func() {
        elapsedTime := 60
        for {
            select {
            case <-quit:
                return
            default:
                if elapsedTime == 60 {
                    console.ClearScreen()
                    console.Render("airspace | live tracking\n\n")
                    aircrafts, err := dump1090.GetAircraft()
                    if err != nil {
                        console.Renderln("error connecting to ther ADS-B receiver")
                        return
                    }
                    for _, aircraft := range aircrafts {
                        registration, _ := faa.LookupByTransponder(aircraft.Transponder)
                        if registration != nil {
                            output := makeAircraftLineOutput(aircraft, registration)
                            console.Renderln(output)
                        }
                    }
                    console.Render("\n\npress any key to exit live tracking\n\n")
                    elapsedTime = 0
                } else {
                    time.Sleep(1 * time.Second)
                    elapsedTime += 1
                }
            }
        }
    }()
    _, err = tty.ReadRune()
    if err != nil {
        log.Println("ERROR : calling tty.ReadRune() in Live()")
        log.Println(spew.Sdump(err))
        panic(err)
    }
    quit <- true
    return
}

func Transponder(command []string) {
    faa, err := faa_registry.Connect()
    if err != nil {
        console.Renderln("error connecting to the faa registry")
        return
    }
    registration, err := faa.LookupByTransponder(command[1])
    switch err {
    case nil:
        // spew.Dump(registration)
        output := makeRegistrationOutput(registration)
        output += fmt.Sprintf("https://flightaware.com/live/flight/%s \n",
            registration.RegistrationNumber,
        )
        console.Render(output)
    case sql.ErrNoRows:
        output := fmt.Sprintf("no registration with transponder '%s' found", command[1])
        console.Renderln(output)
    default:
        console.Renderln("error while seraching the faa registry")
    }
}

func Registration(command []string) {
    faa, err := faa_registry.Connect()
    if err != nil {
        console.Renderln("error connecting to the faa registry")
        return
    }
    registration, err := faa.LookupByRegistration(command[1])
    switch err {
    case nil:
        // spew.Dump(registration)
        output := makeRegistrationOutput(registration)
        output += fmt.Sprintf("https://flightaware.com/resources/registration/%s \n",
            registration.RegistrationNumber,
        )
        console.Render(output)
    case sql.ErrNoRows:
        output := fmt.Sprintf("no registration with registration number '%s' found", command[1])
        console.Renderln(output)
    default:
        console.Renderln("error while seraching the faa registry")
    }
}

func Flight(command []string) {
    flight, err := aviation_stack.LookupActiveFlightByICAO(command[1])
    if err != nil {
        console.Renderln("error connecting to the aviation stack api")
        return
    }
    if flight != nil {
        output := makeFlightOutput(flight)
        console.Render(output)
    } else {
        output := fmt.Sprintf("no flight information for flight '%s' found", command[1])
        console.Renderln(output)
    }
}

func main() {

    // load configuration into env
    configFilePath = os.Getenv("APPLICATION_PATH") + "/.config"
    err := godotenv.Load(configFilePath)
    if err != nil {
        log.Println("ERROR : opening config file in main()")
        log.Println(spew.Sdump(err))
        panic(err)
    }

    // setup logging
    logFilePath = os.Getenv("APPLICATION_PATH") + "/error.log"
    logFile, err := os.OpenFile(logFilePath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Println("ERROR : opening log file in main()")
        log.Println(spew.Sdump(err))
        panic(err)
    }
    defer logFile.Close()
    log.SetOutput(logFile)

    // setup tty
    tty, err := tty.Open()
    if err != nil {
        log.Println("ERROR : starting tty in main()")
        log.Println(spew.Sdump(err))
        panic(err)
    }
    defer tty.Close()

    // startup banner and render initial prompt
    console.Render("\nairspace | tracking console\n")
    console.RenderPrompt()

    // app loop : wait for input, execute command, render prompt for next iteration
    for {
        input, err := tty.ReadString()
        if err != nil {
            log.Println("ERROR : calling tty.ReadString() in main()")
            log.Println(spew.Sdump(err))
            panic(err)
        }
        // spew.Dump(input)

        if len(input) == 0 {
            console.Renderln("use 'help' for available commands, 'exit' to exit the console")
        } else {
            command := strings.Fields(input)
            // spew.Dump(command)

            switch command[0] {

            case "ping":
                Ping()

            case "live":
                Live(tty)

            case "transponder":
                Transponder(command)

            case "registration":
                Registration(command)

            case "flight":
                Flight(command)

            case "help":
                console.Renderln("help command output")

            case "clear":
                console.ClearScreen()

            case "exit":
                console.Renderln("")
                return

            default:
                output := fmt.Sprintf("no valid command '%s', use 'help' for valid commands", command[0])
                console.Renderln(output)

            }
        }

        console.RenderPrompt()
    }

}
