package gomaybe

import (
    "fmt"
    "io/ioutil"
    "os"
)

var (
    cpu Cpu
    rom Rom
    ram Ram
)

func main() {
    fmt.Println("GoMaybe")
    fmt.Println("Jahn Veach <v64@v64.net>")
    fmt.Println("https://github.com/v64/gomaybe")

    if len(os.Args) < 2 {
        fmt.Println("Usage: gomaybe rom.gb")
        return
    }

    file := os.Args[1];

    if romData, err := ioutil.ReadFile(file); err == nil {
        fmt.Println("Loading ROM: " + file)
        ram.Init()
        rom.Init(romData, &ram)
        cpu.Init(&ram)
    } else {
        fmt.Println("Error loading ROM: " + err.Error())
        return
    }

    for {
        cycleCount := cpu.Step()
        if cycleCount == 0 {
            fmt.Println("Unknown opcode encountered, exiting")
            break
        }
    }
}
