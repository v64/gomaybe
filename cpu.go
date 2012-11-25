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
    if (cpu.fReg & (1 << flag) == 1) {
        result = true
    } else {
        result = false
    }

    return
}

func (cpu *Cpu) SetFlag(flag uint8, set bool) {
    if (set) {
        cpu.fReg |= 1 << flag
    } else {
        cpu.fReg &^= 1 << flag
    }
}

func (cpu *Cpu) AND_A(val byte) {
    cpu.aReg &= val

    if (cpu.aReg == 0) {
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

    if (cpu.aReg == 0) {
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

    if (cpu.aReg == 0) {
        cpu.SetFlag(FLAG_Z, true)
    } else {
        cpu.SetFlag(FLAG_Z, false)
    }

    cpu.SetFlag(FLAG_N, false)
    cpu.SetFlag(FLAG_H, false)
    cpu.SetFlag(FLAG_C, false)
}

func (cpu *Cpu) Step() (cycleCount int) {
    opCode := cpu.ram.Read(cpu.pcReg)
    cpu.pcReg++
    fmt.Printf("OP: 0x%.2X\n", opCode)

    switch opCode {
        // START 3.3.3.5 AND n
        case 0xA7:
            cycleCount = 4
            cpu.AND_A(cpu.aReg)

        case 0xA0:
            cycleCount = 4
            cpu.AND_A(cpu.bReg)

        case 0xA1:
            cycleCount = 4
            cpu.AND_A(cpu.cReg)

        case 0xA2:
            cycleCount = 4
            cpu.AND_A(cpu.dReg)

        case 0xA3:
            cycleCount = 4
            cpu.AND_A(cpu.eReg)

        case 0xA4:
            cycleCount = 4
            cpu.AND_A(cpu.hReg)

        case 0xA5:
            cycleCount = 4
            cpu.AND_A(cpu.lReg)
        // END 3.3.3.5 AND n

        // START 3.3.3.6 OR n
        case 0xB7:
            cycleCount = 4
            cpu.OR_A(cpu.aReg)

        case 0xB0:
            cycleCount = 4
            cpu.OR_A(cpu.bReg)

        case 0xB1:
            cycleCount = 4
            cpu.OR_A(cpu.cReg)

        case 0xB2:
            cycleCount = 4
            cpu.OR_A(cpu.dReg)

        case 0xB3:
            cycleCount = 4
            cpu.OR_A(cpu.eReg)

        case 0xB4:
            cycleCount = 4
            cpu.OR_A(cpu.hReg)

        case 0xB5:
            cycleCount = 4
            cpu.OR_A(cpu.lReg)
        // END 3.3.3.6 OR n

        // START 3.3.3.7 XOR n
        case 0xAF:
            cycleCount = 4
            cpu.XOR_A(cpu.aReg)

        case 0xA8:
            cycleCount = 4
            cpu.XOR_A(cpu.bReg)

        case 0xA9:
            cycleCount = 4
            cpu.XOR_A(cpu.cReg)

        case 0xAA:
            cycleCount = 4
            cpu.XOR_A(cpu.dReg)

        case 0xAB:
            cycleCount = 4
            cpu.XOR_A(cpu.eReg)

        case 0xAC:
            cycleCount = 4
            cpu.XOR_A(cpu.hReg)

        case 0xAD:
            cycleCount = 4
            cpu.XOR_A(cpu.lReg)
        // END 3.3.3.7 XOR n

        // START 3.3.5.3 CPL
        case 0x2F:
            cycleCount = 4
            cpu.aReg = ^cpu.aReg
            cpu.SetFlag(FLAG_N, true)
            cpu.SetFlag(FLAG_H, true)
        // END 3.3.5.3 CPL

        // START 3.3.5.4 CCF
        case 0x3F:
            cycleCount = 4
            cpu.SetFlag(FLAG_N, false)
            cpu.SetFlag(FLAG_H, false)
            cpu.SetFlag(FLAG_C, !cpu.GetFlag(FLAG_C))
        // END 3.3.5.4 CCF

        // START 3.3.5.5 SCF
        case 0x37:
            cycleCount = 4
            cpu.SetFlag(FLAG_N, false)
            cpu.SetFlag(FLAG_H, false)
            cpu.SetFlag(FLAG_C, true)
        // END 3.3.5.5 SCF

        // START 3.3.5.6 NOP
        case 0x00:
            cycleCount = 4
        // END 3.3.5.6 NOP

        // START 3.3.8.1 JP nn
        case 0xC3:
            cycleCount = 16
            least := cpu.ram.Read(cpu.pcReg)
            cpu.pcReg++
            most := cpu.ram.Read(cpu.pcReg)
            cpu.pcReg = bytesToUint16(least, most)
        // END 3.3.8.1 JP nn

        default:
            cycleCount = 0
            fmt.Printf("Unknown OP: 0x%.2X\n", opCode)
    }

    return
}

func bytesToUint16(least byte, most byte) (uint16) {
    return uint16(most)<<8 + uint16(least)
}
