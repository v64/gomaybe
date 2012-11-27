package cpu

import (
    "fmt"
    . "../ram"
)

type Cpu struct {
    aReg, bReg, cReg, dReg, eReg, fReg, hReg, lReg uint8
    spReg, pcReg uint16
    ram *Ram
}

const (
    flag_C = uint8(4)
    flag_H = uint8(5)
    flag_N = uint8(6)
    flag_Z = uint8(7)
)

var (
    opCodeTable = map[byte]func(cpu *Cpu) int {
        // START 3.3.3.5 AND n
        0xA7: func(cpu *Cpu) int { cpu.and_A(cpu.aReg); return 4 },
        0xA0: func(cpu *Cpu) int { cpu.and_A(cpu.bReg); return 4 },
        0xA1: func(cpu *Cpu) int { cpu.and_A(cpu.cReg); return 4 },
        0xA2: func(cpu *Cpu) int { cpu.and_A(cpu.dReg); return 4 },
        0xA3: func(cpu *Cpu) int { cpu.and_A(cpu.eReg); return 4 },
        0xA4: func(cpu *Cpu) int { cpu.and_A(cpu.hReg); return 4 },
        0xA5: func(cpu *Cpu) int { cpu.and_A(cpu.lReg); return 4 },
        // END 3.3.3.5 AND n

        // START 3.3.3.6 OR n
        0xB7: func(cpu *Cpu) int { cpu.or_A(cpu.aReg); return 4 },
        0xB0: func(cpu *Cpu) int { cpu.or_A(cpu.bReg); return 4 },
        0xB1: func(cpu *Cpu) int { cpu.or_A(cpu.cReg); return 4 },
        0xB2: func(cpu *Cpu) int { cpu.or_A(cpu.dReg); return 4 },
        0xB3: func(cpu *Cpu) int { cpu.or_A(cpu.eReg); return 4 },
        0xB4: func(cpu *Cpu) int { cpu.or_A(cpu.hReg); return 4 },
        0xB5: func(cpu *Cpu) int { cpu.or_A(cpu.lReg); return 4 },
        // END 3.3.3.6 OR n

        // START 3.3.3.7 XOR n
        0xAF: func(cpu *Cpu) int { cpu.xor_A(cpu.aReg); return 4 },
        0xA8: func(cpu *Cpu) int { cpu.xor_A(cpu.bReg); return 4 },
        0xA9: func(cpu *Cpu) int { cpu.xor_A(cpu.cReg); return 4 },
        0xAA: func(cpu *Cpu) int { cpu.xor_A(cpu.dReg); return 4 },
        0xAB: func(cpu *Cpu) int { cpu.xor_A(cpu.eReg); return 4 },
        0xAC: func(cpu *Cpu) int { cpu.xor_A(cpu.hReg); return 4 },
        0xAD: func(cpu *Cpu) int { cpu.xor_A(cpu.lReg); return 4 },
        // END 3.3.3.7 XOR n

        // START 3.3.5.3 CPL
        0x2F: func(cpu *Cpu) int {
            cpu.aReg = ^cpu.aReg
            cpu.setFlag(flag_N, true)
            cpu.setFlag(flag_H, true)
            return 4
        },
        // END 3.3.5.3 CPL

        // START 3.3.5.4 CCF
        0x3F: func(cpu *Cpu) int {
            cpu.setFlag(flag_N, false)
            cpu.setFlag(flag_H, false)
            cpu.setFlag(flag_C, !cpu.getFlag(flag_C))
            return 4
        },
        // END 3.3.5.4 CCF

        // START 3.3.5.5 SCF
        0x37: func(cpu *Cpu) int {
            cpu.setFlag(flag_N, false)
            cpu.setFlag(flag_H, false)
            cpu.setFlag(flag_C, true)
            return 4
        },
        // END 3.3.5.5 SCF

        // START 3.3.5.6 NOP
        0x00: func(cpu *Cpu) int { return 4 },
        // END 3.3.5.6 NOP

        // START 3.3.8.1 JP nn
        0xC3: func(cpu *Cpu) int {
            least := cpu.ram.Read(cpu.pcReg)
            cpu.pcReg++
            most := cpu.ram.Read(cpu.pcReg)
            cpu.pcReg = bytesToUint16(least, most)
            return 16
        },
        // END 3.3.8.1 JP nn
    }
)

func (cpu *Cpu) Init(ram *Ram) {
    cpu.spReg = 0xFFFE
    cpu.pcReg = 0x0100
    cpu.ram = ram
}

func (cpu *Cpu) Step() (cycles int) {
    opCode := cpu.ram.Read(cpu.pcReg)
    cpu.pcReg++
    fmt.Printf("OP: 0x%.2X\n", opCode)

    instruction, ok := opCodeTable[opCode]
    if ok {
        cycles = instruction(cpu)
    } else {
        fmt.Printf("Unknown OP: 0x%.2X\n", opCode)
        cycles = -1
    }

    return
}

func (cpu *Cpu) getFlag(flag uint8) bool {
    return cpu.fReg & (1 << flag) == 1
}

func (cpu *Cpu) setFlag(flag uint8, set bool) {
    if set {
        cpu.fReg |= 1 << flag
    } else {
        cpu.fReg &^= 1 << flag
    }
}

func (cpu *Cpu) and_A(val byte) {
    cpu.aReg &= val
    cpu.setFlag(flag_Z, cpu.aReg == 0)
    cpu.setFlag(flag_N, false)
    cpu.setFlag(flag_H, true)
    cpu.setFlag(flag_C, false)
}

func (cpu *Cpu) or_A(val byte) {
    cpu.aReg |= val
    cpu.setFlag(flag_Z, cpu.aReg == 0)
    cpu.setFlag(flag_N, false)
    cpu.setFlag(flag_H, false)
    cpu.setFlag(flag_C, false)
}

func (cpu *Cpu) xor_A(val byte) {
    cpu.aReg ^= val
    cpu.setFlag(flag_Z, cpu.aReg == 0)
    cpu.setFlag(flag_N, false)
    cpu.setFlag(flag_H, false)
    cpu.setFlag(flag_C, false)
}

func bytesToUint16(least byte, most byte) uint16 {
    return uint16(most)<<8 + uint16(least)
}
