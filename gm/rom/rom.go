package rom

import (
    "fmt"
    "../ram"
)

type Rom struct {
    title string
    cartridgeType byte
}

func (rom *Rom) Init(romData []byte, ram *ram.Ram) {
    rom.title = string(romData[0x0134:0x0142])
    fmt.Println("Title: " + rom.title)

    rom.cartridgeType = romData[0x0147]
    fmt.Println("Cartridge: " + getCartridgeTypeStr(rom.cartridgeType))

    switch rom.cartridgeType {
        case 0x00:
            ram.WriteBlock(0x0000, romData)
        default:
            fmt.Printf("Don't know how to handle cartridge type %.2X\n", rom.cartridgeType)
    }
}

func getCartridgeTypeStr(cartridgeType byte) (cartridgeTypeStr string) {
    switch cartridgeType {
        case 0x00:
            cartridgeTypeStr = "ROM"
        case 0x01:
            cartridgeTypeStr = "ROM+MBC1"
        case 0x02:
            cartridgeTypeStr = "ROM+MBC1+RAM"
        case 0x03:
            cartridgeTypeStr = "ROM+MBC1+RAM+BATTERY"
        case 0x05:
            cartridgeTypeStr = "ROM+MBC2"
        case 0x06:
            cartridgeTypeStr = "ROM+MBC2+BATTERY"
        case 0x08:
            cartridgeTypeStr = "ROM+RAM"
        case 0x09:
            cartridgeTypeStr = "ROM+RAM+BATTERY"
        case 0x0B:
            cartridgeTypeStr = "ROM+MMMD1"
        case 0x0C:
            cartridgeTypeStr = "ROM+MMMD1+SRAM"
        case 0x0D:
            cartridgeTypeStr = "ROM+MMMD1+SRAM+BATTERY"
        case 0x0F:
            cartridgeTypeStr = "ROM+MBC3+TIMER+BATTERY"
        case 0x10:
            cartridgeTypeStr = "ROM+MBC3+TIMER+RAM+BATTERY"
        case 0x11:
            cartridgeTypeStr = "ROM+MBC3"
        case 0x12:
            cartridgeTypeStr = "ROM+MBC3+RAM"
        case 0x13:
            cartridgeTypeStr = "ROM+MBC3+RAM+BATTERY"
        case 0x19:
            cartridgeTypeStr = "ROM+MBC5"
        case 0x1A:
            cartridgeTypeStr = "ROM+MBC5+RAM"
        case 0x1B:
            cartridgeTypeStr = "ROM+MBC5+RAM+BATTERY"
        case 0x1C:
            cartridgeTypeStr = "ROM+MBC5+RUMBLE"
        case 0x1D:
            cartridgeTypeStr = "ROM+MBC5+RUMBLE+SRAM"
        case 0x1E:
            cartridgeTypeStr = "ROM+MBC5+RUMBLE+SRAM+BATTERY"
        case 0x1F:
            cartridgeTypeStr = "Pocket Camera"
        case 0xFD:
            cartridgeTypeStr = "Bandai TAMA5"
        case 0xFE:
            cartridgeTypeStr = "Hudson HuC-3"
        case 0xFF:
            cartridgeTypeStr = "Hudson HuC-1"
        default:
            cartridgeTypeStr = "Unknown"
    }

    return
}
