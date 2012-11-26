package main

import (
    "fmt"
)

const (
    FLAG_C = uint8(4)
    FLAG_H = uint8(5)
    FLAG_N = uint8(6)
    FLAG_Z = uint8(7)
)

var (
    opCodeCycles = map[byte]int{
        0xA7: 4,
        0xA0: 4,
        0xA1: 4,
        0xA2: 4,
        0xA3: 4,
        0xA4: 4,
        0xA5: 4,
        0xB7: 4,
        0xB0: 4,
        0xB1: 4,
        0xB2: 4,
        0xB3: 4,
        0xB4: 4,
        0xB5: 4,
        0xAF: 4,
        0xA8: 4,
        0xA9: 4,
        0xAA: 4,
        0xAB: 4,
        0xAC: 4,
        0xAD: 4,
        0x2F: 4,
        0x3F: 4,
        0x37: 4,
        0x00: 4,
        0xC3: 16,
    }
)

type Cpu struct {
    aReg, bReg, cReg, dReg, eReg, fReg, hReg, lReg uint8
    spReg, pcReg uint16
    ram *Ram
}

func (cpu *Cpu) Init(ram *Ram) {
    cpu.spReg = 0xFFFE
    cpu.pcReg = 0x0100
    cpu.ram = ram
}

func (cpu *Cpu) GetFlag(flag uint8) (result bool) {
    if cpu.fReg & (1 << flag) == 1 {
        result = true
    } else {
        result = false
    }

    return
}

func (cpu *Cpu) SetFlag(flag uint8, set bool) {
    if set {
        cpu.fReg |= 1 << flag
    } else {
        cpu.fReg &^= 1 << flag
    }
}

func (cpu *Cpu) AND_A(val byte) {
    cpu.aReg &= val

    if cpu.aReg == 0 {
        cpu.SetFlag(FLAG_Z, true)
    } else {
        cpu.SetFlag(FLAG_Z, false)
    }

    cpu.SetFlag(FLAG_N, false)
    cpu.SetFlag(FLAG_H, true)
    cpu.SetFlag(FLAG_C, false)
}

func (cpu *Cpu) OR_A(val byte) {
    cpu.aReg |= val

    if cpu.aReg == 0 {
        cpu.SetFlag(FLAG_Z, true)
    } else {
        cpu.SetFlag(FLAG_Z, false)
    }

    cpu.SetFlag(FLAG_N, false)
    cpu.SetFlag(FLAG_H, false)
    cpu.SetFlag(FLAG_C, false)
}

func (cpu *Cpu) XOR_A(val byte) {
    cpu.aReg ^= val

    if cpu.aReg == 0 {
        cpu.SetFlag(FLAG_Z, true)
    } else {
        cpu.SetFlag(FLAG_Z, false)
    }

    cpu.SetFlag(FLAG_N, false)
    cpu.SetFlag(FLAG_H, false)
    cpu.SetFlag(FLAG_C, false)
}

func (cpu *Cpu) Step() int {
    opCode := cpu.ram.Read(cpu.pcReg)
    cpu.pcReg++
    fmt.Printf("OP: 0x%.2X\n", opCode)

    switch opCode {
        // START 3.3.3.5 AND n
        case 0xA7:
            cpu.AND_A(cpu.aReg)

        case 0xA0:
            cpu.AND_A(cpu.bReg)

        case 0xA1:
            cpu.AND_A(cpu.cReg)

        case 0xA2:
            cpu.AND_A(cpu.dReg)

        case 0xA3:
            cpu.AND_A(cpu.eReg)

        case 0xA4:
            cpu.AND_A(cpu.hReg)

        case 0xA5:
            cpu.AND_A(cpu.lReg)
        // END 3.3.3.5 AND n

        // START 3.3.3.6 OR n
        case 0xB7:
            cpu.OR_A(cpu.aReg)

        case 0xB0:
            cpu.OR_A(cpu.bReg)

        case 0xB1:
            cpu.OR_A(cpu.cReg)

        case 0xB2:
            cpu.OR_A(cpu.dReg)

        case 0xB3:
            cpu.OR_A(cpu.eReg)

        case 0xB4:
            cpu.OR_A(cpu.hReg)

        case 0xB5:
            cpu.OR_A(cpu.lReg)
        // END 3.3.3.6 OR n

        // START 3.3.3.7 XOR n
        case 0xAF:
            cpu.XOR_A(cpu.aReg)

        case 0xA8:
            cpu.XOR_A(cpu.bReg)

        case 0xA9:
            cpu.XOR_A(cpu.cReg)

        case 0xAA:
            cpu.XOR_A(cpu.dReg)

        case 0xAB:
            cpu.XOR_A(cpu.eReg)

        case 0xAC:
            cpu.XOR_A(cpu.hReg)

        case 0xAD:
            cpu.XOR_A(cpu.lReg)
        // END 3.3.3.7 XOR n

        // START 3.3.5.3 CPL
        case 0x2F:
            cpu.aReg = ^cpu.aReg
            cpu.SetFlag(FLAG_N, true)
            cpu.SetFlag(FLAG_H, true)
        // END 3.3.5.3 CPL

        // START 3.3.5.4 CCF
        case 0x3F:
            cpu.SetFlag(FLAG_N, false)
            cpu.SetFlag(FLAG_H, false)
            cpu.SetFlag(FLAG_C, !cpu.GetFlag(FLAG_C))
        // END 3.3.5.4 CCF

        // START 3.3.5.5 SCF
        case 0x37:
            cpu.SetFlag(FLAG_N, false)
            cpu.SetFlag(FLAG_H, false)
            cpu.SetFlag(FLAG_C, true)
        // END 3.3.5.5 SCF

        // START 3.3.5.6 NOP
        case 0x00:
        // END 3.3.5.6 NOP

        // START 3.3.8.1 JP nn
        case 0xC3:
            least := cpu.ram.Read(cpu.pcReg)
            cpu.pcReg++
            most := cpu.ram.Read(cpu.pcReg)
            cpu.pcReg = bytesToUint16(least, most)
        // END 3.3.8.1 JP nn

        default:
            fmt.Printf("Unknown OP: 0x%.2X\n", opCode)
    }

    return opCodeCycles[opCode]
}

func bytesToUint16(least byte, most byte) (uint16) {
    return uint16(most)<<8 + uint16(least)
}
