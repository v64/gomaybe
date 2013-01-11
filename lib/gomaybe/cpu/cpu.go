package cpu

import (
    . "../ram"
    "../util"
    "fmt"
)

type Cpu struct {
    aReg, bReg, cReg, dReg, eReg, fReg, hReg, lReg uint8
    spReg, pcReg                                   uint16
    ram                                            *Ram
}

const (
    flag_C = uint8(4)
    flag_H = uint8(5)
    flag_N = uint8(6)
    flag_Z = uint8(7)
)

var (
    opCodeTable = map[byte]func(cpu *Cpu) int{
        0x00: func(cpu *Cpu) int { return 4 },
        0x01: func(cpu *Cpu) int { cpu.bReg, cpu.cReg = cpu.ram.ReadWordSplit(cpu.pcReg); cpu.pcReg += 2; return 12 },
        0x02: func(cpu *Cpu) int { cpu.ram.Write(cpu.bcReg(), cpu.aReg); return 8 },
        0x04: func(cpu *Cpu) int { cpu.incReg(&cpu.bReg); return 4 },
        0x06: func(cpu *Cpu) int { cpu.bReg = cpu.ram.Read(cpu.pcReg); cpu.pcReg++; return 8 },
        0x0A: func(cpu *Cpu) int { cpu.aReg = cpu.ram.Read(cpu.bcReg()); return 8 },
        0x0C: func(cpu *Cpu) int { cpu.incReg(&cpu.cReg); return 4 },
        0x0E: func(cpu *Cpu) int { cpu.cReg = cpu.ram.Read(cpu.pcReg); cpu.pcReg++; return 8 },
        0x11: func(cpu *Cpu) int { cpu.dReg, cpu.eReg = cpu.ram.ReadWordSplit(cpu.pcReg); cpu.pcReg += 2; return 12 },
        0x12: func(cpu *Cpu) int { cpu.ram.Write(cpu.deReg(), cpu.aReg); return 8 },
        0x14: func(cpu *Cpu) int { cpu.incReg(&cpu.dReg); return 4 },
        0x16: func(cpu *Cpu) int { cpu.dReg = cpu.ram.Read(cpu.pcReg); cpu.pcReg++; return 8 },
        0x1A: func(cpu *Cpu) int { cpu.aReg = cpu.ram.Read(cpu.deReg()); return 8 },
        0x1C: func(cpu *Cpu) int { cpu.incReg(&cpu.eReg); return 4 },
        0x1E: func(cpu *Cpu) int { cpu.eReg = cpu.ram.Read(cpu.pcReg); cpu.pcReg++; return 8 },
        0x20: func(cpu *Cpu) int {
            if !cpu.getFlag(flag_Z) {
                cpu.pcReg += uint16(int8(cpu.ram.Read(cpu.pcReg))) + 1
                return 12
            }

            cpu.pcReg++
            return 8
        },
        0x21: func(cpu *Cpu) int { cpu.hReg, cpu.lReg = cpu.ram.ReadWordSplit(cpu.pcReg); cpu.pcReg += 2; return 12 },
        0x24: func(cpu *Cpu) int { cpu.incReg(&cpu.hReg); return 4 },
        0x26: func(cpu *Cpu) int { cpu.hReg = cpu.ram.Read(cpu.pcReg); cpu.pcReg++; return 8 },
        0x2C: func(cpu *Cpu) int { cpu.incReg(&cpu.lReg); return 4 },
        0x2E: func(cpu *Cpu) int { cpu.lReg = cpu.ram.Read(cpu.pcReg); cpu.pcReg++; return 8 },
        0x2F: func(cpu *Cpu) int {
            cpu.aReg = ^cpu.aReg
            cpu.setFlag(flag_N, true)
            cpu.setFlag(flag_H, true)
            return 4
        },
        0x31: func(cpu *Cpu) int { cpu.spReg = cpu.ram.ReadWord(cpu.pcReg); cpu.pcReg += 2; return 12 },
        0x32: func(cpu *Cpu) int { cpu.ram.Write(cpu.hlReg(), cpu.aReg); cpu.chgRegVal("hl", -1); return 8 },
        0x36: func(cpu *Cpu) int { cpu.ram.Write(cpu.hlReg(), cpu.ram.Read(cpu.pcReg)); cpu.pcReg++; return 12 },
        0x37: func(cpu *Cpu) int {
            cpu.setFlag(flag_N, false)
            cpu.setFlag(flag_H, false)
            cpu.setFlag(flag_C, true)
            return 4
        },
        0x3C: func(cpu *Cpu) int { cpu.incReg(&cpu.aReg); return 4 },
        0x3E: func(cpu *Cpu) int { cpu.aReg = cpu.ram.Read(cpu.pcReg); cpu.pcReg++; return 8 },
        0x3F: func(cpu *Cpu) int {
            cpu.setFlag(flag_N, false)
            cpu.setFlag(flag_H, false)
            cpu.setFlag(flag_C, !cpu.getFlag(flag_C))
            return 4
        },
        0x40: func(cpu *Cpu) int { return 4 },
        0x41: func(cpu *Cpu) int { cpu.bReg = cpu.cReg; return 4 },
        0x42: func(cpu *Cpu) int { cpu.bReg = cpu.dReg; return 4 },
        0x43: func(cpu *Cpu) int { cpu.bReg = cpu.eReg; return 4 },
        0x44: func(cpu *Cpu) int { cpu.bReg = cpu.hReg; return 4 },
        0x45: func(cpu *Cpu) int { cpu.bReg = cpu.lReg; return 4 },
        0x46: func(cpu *Cpu) int { cpu.bReg = cpu.ram.Read(cpu.hlReg()); return 8 },
        0x47: func(cpu *Cpu) int { cpu.bReg = cpu.aReg; return 4 },
        0x48: func(cpu *Cpu) int { cpu.cReg = cpu.bReg; return 4 },
        0x49: func(cpu *Cpu) int { return 4 },
        0x4A: func(cpu *Cpu) int { cpu.cReg = cpu.dReg; return 4 },
        0x4B: func(cpu *Cpu) int { cpu.cReg = cpu.eReg; return 4 },
        0x4C: func(cpu *Cpu) int { cpu.cReg = cpu.hReg; return 4 },
        0x4D: func(cpu *Cpu) int { cpu.cReg = cpu.lReg; return 4 },
        0x4E: func(cpu *Cpu) int { cpu.cReg = cpu.ram.Read(cpu.hlReg()); return 8 },
        0x4F: func(cpu *Cpu) int { cpu.cReg = cpu.aReg; return 4 },
        0x50: func(cpu *Cpu) int { cpu.dReg = cpu.bReg; return 4 },
        0x51: func(cpu *Cpu) int { cpu.dReg = cpu.cReg; return 4 },
        0x52: func(cpu *Cpu) int { return 4 },
        0x53: func(cpu *Cpu) int { cpu.dReg = cpu.eReg; return 4 },
        0x54: func(cpu *Cpu) int { cpu.dReg = cpu.hReg; return 4 },
        0x55: func(cpu *Cpu) int { cpu.dReg = cpu.lReg; return 4 },
        0x56: func(cpu *Cpu) int { cpu.dReg = cpu.ram.Read(cpu.hlReg()); return 8 },
        0x57: func(cpu *Cpu) int { cpu.dReg = cpu.aReg; return 4 },
        0x58: func(cpu *Cpu) int { cpu.eReg = cpu.bReg; return 4 },
        0x59: func(cpu *Cpu) int { cpu.eReg = cpu.cReg; return 4 },
        0x5A: func(cpu *Cpu) int { cpu.eReg = cpu.dReg; return 4 },
        0x5B: func(cpu *Cpu) int { return 4 },
        0x5C: func(cpu *Cpu) int { cpu.eReg = cpu.hReg; return 4 },
        0x5D: func(cpu *Cpu) int { cpu.eReg = cpu.lReg; return 4 },
        0x5E: func(cpu *Cpu) int { cpu.eReg = cpu.ram.Read(cpu.hlReg()); return 8 },
        0x5F: func(cpu *Cpu) int { cpu.eReg = cpu.aReg; return 4 },
        0x60: func(cpu *Cpu) int { cpu.hReg = cpu.bReg; return 4 },
        0x61: func(cpu *Cpu) int { cpu.hReg = cpu.cReg; return 4 },
        0x62: func(cpu *Cpu) int { cpu.hReg = cpu.dReg; return 4 },
        0x63: func(cpu *Cpu) int { cpu.hReg = cpu.eReg; return 4 },
        0x64: func(cpu *Cpu) int { return 4 },
        0x65: func(cpu *Cpu) int { cpu.hReg = cpu.lReg; return 4 },
        0x66: func(cpu *Cpu) int { cpu.hReg = cpu.ram.Read(cpu.hlReg()); return 8 },
        0x67: func(cpu *Cpu) int { cpu.hReg = cpu.aReg; return 4 },
        0x68: func(cpu *Cpu) int { cpu.lReg = cpu.bReg; return 4 },
        0x69: func(cpu *Cpu) int { cpu.lReg = cpu.cReg; return 4 },
        0x6A: func(cpu *Cpu) int { cpu.lReg = cpu.dReg; return 4 },
        0x6B: func(cpu *Cpu) int { cpu.lReg = cpu.eReg; return 4 },
        0x6C: func(cpu *Cpu) int { cpu.lReg = cpu.hReg; return 4 },
        0x6D: func(cpu *Cpu) int { return 4 },
        0x6E: func(cpu *Cpu) int { cpu.lReg = cpu.ram.Read(cpu.hlReg()); return 8 },
        0x6F: func(cpu *Cpu) int { cpu.lReg = cpu.aReg; return 4 },
        0x70: func(cpu *Cpu) int { cpu.ram.Write(cpu.hlReg(), cpu.bReg); return 8 },
        0x71: func(cpu *Cpu) int { cpu.ram.Write(cpu.hlReg(), cpu.cReg); return 8 },
        0x72: func(cpu *Cpu) int { cpu.ram.Write(cpu.hlReg(), cpu.dReg); return 8 },
        0x73: func(cpu *Cpu) int { cpu.ram.Write(cpu.hlReg(), cpu.eReg); return 8 },
        0x74: func(cpu *Cpu) int { cpu.ram.Write(cpu.hlReg(), cpu.hReg); return 8 },
        0x75: func(cpu *Cpu) int { cpu.ram.Write(cpu.hlReg(), cpu.lReg); return 8 },
        0x77: func(cpu *Cpu) int { cpu.ram.Write(cpu.hlReg(), cpu.aReg); return 8 },
        0x78: func(cpu *Cpu) int { cpu.aReg = cpu.bReg; return 4 },
        0x79: func(cpu *Cpu) int { cpu.aReg = cpu.cReg; return 4 },
        0x7A: func(cpu *Cpu) int { cpu.aReg = cpu.dReg; return 4 },
        0x7B: func(cpu *Cpu) int { cpu.aReg = cpu.eReg; return 4 },
        0x7C: func(cpu *Cpu) int { cpu.aReg = cpu.hReg; return 4 },
        0x7D: func(cpu *Cpu) int { cpu.aReg = cpu.lReg; return 4 },
        0x7E: func(cpu *Cpu) int { cpu.aReg = cpu.ram.Read(cpu.hlReg()); return 8 },
        0x7F: func(cpu *Cpu) int { return 4 },
        0xA0: func(cpu *Cpu) int { cpu.and_A(cpu.bReg); return 4 },
        0xA1: func(cpu *Cpu) int { cpu.and_A(cpu.cReg); return 4 },
        0xA2: func(cpu *Cpu) int { cpu.and_A(cpu.dReg); return 4 },
        0xA3: func(cpu *Cpu) int { cpu.and_A(cpu.eReg); return 4 },
        0xA4: func(cpu *Cpu) int { cpu.and_A(cpu.hReg); return 4 },
        0xA5: func(cpu *Cpu) int { cpu.and_A(cpu.lReg); return 4 },
        0xA6: func(cpu *Cpu) int { cpu.and_A(cpu.ram.Read(cpu.hlReg())); return 8 },
        0xA7: func(cpu *Cpu) int { cpu.and_A(cpu.aReg); return 4 },
        0xA8: func(cpu *Cpu) int { cpu.xor_A(cpu.bReg); return 4 },
        0xA9: func(cpu *Cpu) int { cpu.xor_A(cpu.cReg); return 4 },
        0xAA: func(cpu *Cpu) int { cpu.xor_A(cpu.dReg); return 4 },
        0xAE: func(cpu *Cpu) int { cpu.xor_A(cpu.ram.Read(cpu.hlReg())); return 8 },
        0xAB: func(cpu *Cpu) int { cpu.xor_A(cpu.eReg); return 4 },
        0xAC: func(cpu *Cpu) int { cpu.xor_A(cpu.hReg); return 4 },
        0xAD: func(cpu *Cpu) int { cpu.xor_A(cpu.lReg); return 4 },
        0xAF: func(cpu *Cpu) int { cpu.xor_A(cpu.aReg); return 4 },
        0xB0: func(cpu *Cpu) int { cpu.or_A(cpu.bReg); return 4 },
        0xB1: func(cpu *Cpu) int { cpu.or_A(cpu.cReg); return 4 },
        0xB2: func(cpu *Cpu) int { cpu.or_A(cpu.dReg); return 4 },
        0xB3: func(cpu *Cpu) int { cpu.or_A(cpu.eReg); return 4 },
        0xB4: func(cpu *Cpu) int { cpu.or_A(cpu.hReg); return 4 },
        0xB5: func(cpu *Cpu) int { cpu.or_A(cpu.lReg); return 4 },
        0xB6: func(cpu *Cpu) int { cpu.or_A(cpu.ram.Read(cpu.hlReg())); return 8 },
        0xB7: func(cpu *Cpu) int { cpu.or_A(cpu.aReg); return 4 },
        0xC3: func(cpu *Cpu) int { cpu.pcReg = cpu.ram.ReadWord(cpu.pcReg); return 16 },
        0xCD: func(cpu *Cpu) int { cpu.call(); return 24 },
        0xE0: func(cpu *Cpu) int {
            cpu.ram.Write(0xFF00+uint16(cpu.ram.Read(cpu.pcReg)), cpu.aReg)
            cpu.pcReg++
            return 12
        },
        0xE2: func(cpu *Cpu) int { cpu.ram.Write(0xFF00+uint16(cpu.cReg), cpu.aReg); return 8 },
        0xE6: func(cpu *Cpu) int { cpu.and_A(cpu.ram.Read(cpu.pcReg)); cpu.pcReg++; return 8 },
        0xEA: func(cpu *Cpu) int { cpu.ram.Write(cpu.ram.ReadWord(cpu.pcReg), cpu.aReg); cpu.pcReg += 2; return 16 },
        0xEE: func(cpu *Cpu) int { cpu.xor_A(cpu.ram.Read(cpu.pcReg)); cpu.pcReg++; return 8 },
        0xFA: func(cpu *Cpu) int { cpu.aReg = cpu.ram.Read(cpu.ram.ReadWord(cpu.pcReg)); cpu.pcReg += 2; return 16 },
        0xF6: func(cpu *Cpu) int { cpu.or_A(cpu.ram.Read(cpu.pcReg)); cpu.pcReg++; return 8 },
        0xF9: func(cpu *Cpu) int { cpu.spReg = cpu.hlReg(); return 8 },
    }

    cbOpCodeTable = map[byte]func(cpu *Cpu) int{
        0x7C: func(cpu *Cpu) int {
            cpu.setFlag(flag_Z, !cpu.getBit(cpu.hReg, 7))
            cpu.setFlag(flag_N, false)
            cpu.setFlag(flag_H, true)
            return 8
        },
    }
)

func (cpu *Cpu) Init(ram *Ram) {
    cpu.spReg = 0xFFFE
    cpu.pcReg = 0x0000
    cpu.ram = ram
}

func (cpu *Cpu) Step() (cycles int) {
    var (
        instruction func(cpu *Cpu) int
        ok          bool
    )

    opCode := cpu.nextOpCode()

    if opCode == 0xCB {
        opCode = cpu.nextOpCode()
        instruction, ok = cbOpCodeTable[opCode]
    } else {
        instruction, ok = opCodeTable[opCode]
    }

    if ok {
        cycles = instruction(cpu)
        fmt.Println(cpu)
    } else {
        fmt.Printf("Unknown OP: 0x%.2X\n", opCode)
        cycles = -1
    }

    return
}

func (cpu *Cpu) nextOpCode() (opCode byte) {
    opCode = cpu.ram.Read(cpu.pcReg)
    cpu.pcReg++
    fmt.Printf("OP: 0x%.2X\n", opCode)
    return
}

func (cpu *Cpu) String() string {
    return fmt.Sprintf("pc:0x%.4X sp:0x%.4X a:0x%.2X b:0x%.2X c:0x%.2X d:0x%.2X e:0x%.2X f:0x%.2X h:0x%.2X l:0x%.2X",
        cpu.pcReg, cpu.spReg, cpu.aReg, cpu.bReg, cpu.cReg, cpu.dReg, cpu.eReg, cpu.fReg, cpu.hReg, cpu.lReg)
}

func (cpu *Cpu) getBit(val uint8, bit uint8) bool {
    return val&(1<<bit) != 0
}

func (cpu *Cpu) setBit(reg *uint8, bit uint8, set bool) {
    if set {
        *reg |= 1 << bit
    } else {
        *reg &^= 1 << bit
    }
}

func (cpu *Cpu) getFlag(flag uint8) bool {
    return cpu.getBit(cpu.fReg, flag)
}

func (cpu *Cpu) setFlag(flag uint8, set bool) {
    cpu.setBit(&cpu.fReg, flag, set)
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

func (cpu *Cpu) incReg(reg *byte) {
    *reg += 1
    cpu.setFlag(flag_Z, *reg == 0)
    cpu.setFlag(flag_N, false)
    cpu.setFlag(flag_H, *reg&0x0F == 0)
}

func (cpu *Cpu) afReg() uint16 {
    return util.B2W(cpu.aReg, cpu.fReg)
}

func (cpu *Cpu) bcReg() uint16 {
    return util.B2W(cpu.bReg, cpu.cReg)
}

func (cpu *Cpu) deReg() uint16 {
    return util.B2W(cpu.dReg, cpu.eReg)
}

func (cpu *Cpu) hlReg() uint16 {
    return util.B2W(cpu.hReg, cpu.lReg)
}

func (cpu *Cpu) chgRegVal(reg string, chgVal int) {
    var (
        origVal         uint16
        highReg, lowReg *byte
    )

    switch reg {
    case "af":
        origVal = cpu.afReg()
        highReg = &cpu.aReg
        lowReg = &cpu.fReg
    case "bc":
        origVal = cpu.bcReg()
        highReg = &cpu.bReg
        lowReg = &cpu.cReg
    case "de":
        origVal = cpu.deReg()
        highReg = &cpu.dReg
        lowReg = &cpu.eReg
    case "hl":
        origVal = cpu.hlReg()
        highReg = &cpu.hReg
        lowReg = &cpu.lReg
    }

    if chgVal >= 0 {
        origVal += uint16(chgVal)
    } else {
        origVal -= uint16(chgVal * -1)
    }

    *highReg, *lowReg = util.W2B(origVal)
}

func (cpu *Cpu) push(val uint16) {
    cpu.spReg -= 2
    cpu.ram.WriteWord(cpu.spReg, val)
}

func (cpu *Cpu) pop() (val uint16) {
    val = cpu.ram.ReadWord(cpu.spReg)
    cpu.spReg += 2
    return
}

func (cpu *Cpu) call() {
    loc := cpu.ram.ReadWord(cpu.pcReg)
    cpu.push(cpu.pcReg)
    cpu.pcReg = loc
}
