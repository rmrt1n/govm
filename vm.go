// package govm
package main

import (
    "fmt"
    "encoding/binary"
)

const (
    MAX_GLOBALS = 256
    MAX_CONSTANTS = 256
    MAX_STACK = 256
)

type Frame struct {
    ret, bp, argc int        // ret address (last ip before call)
    caller *Frame  // parent call frame
}

type VM struct {
    instr Instructions  // instructions list
    ip, sp int              // ip
    globals []Object    // global var map
    constants []Object  // constant pool
    stack []Object      // operand stack
    frame *Frame       // call stack (linked list)
}

func VMNew(ins Instructions) *VM {
    var vm VM;
    vm.instr = ins
    vm.ip, vm.sp = -1, 0
    vm.globals = make([]Object, MAX_GLOBALS)
    vm.constants = make([]Object, MAX_CONSTANTS)
    vm.stack = make([]Object, MAX_STACK)
    return &vm
}

func FrameNew(ret, argc, bp int, caller *Frame) *Frame {
    return &Frame{ret, argc, bp, caller}
}

func (vm *VM) stackPush(o Object) error {
    if vm.sp > MAX_STACK {
        return fmt.Errorf("error: stack overflow")
    }
    vm.stack[vm.sp] = o
    vm.sp++
    return nil
}

func (vm *VM) stackPop() Object {
    vm.sp--
    return vm.stack[vm.sp]
}

func (vm *VM) stackTop() Object {
    return vm.stack[vm.sp-1]
}

func (vm *VM) fetchByte() uint8 {
    vm.ip++
    res := vm.instr[vm.ip]
    return res
}

func (vm *VM) fetch2Bytes() uint16 {
    vm.ip++
    res := binary.BigEndian.Uint16(vm.instr[vm.ip:]);
    vm.ip++
    return res
}

func (vm *VM) Run() error {
    for vm.ip < len(vm.instr) - 1 {
        vm.ip++
        op := OpCode(vm.instr[vm.ip])
        switch op {
        case OP_HALT:
            return nil
        case OP_CONST:
            i := vm.fetchByte()
            e := vm.stackPush(vm.constants[i])
            if e != nil { return e }
        case OP_IADD, OP_ISUB, OP_IMUL, OP_IDIV, OP_IREM:
            e := vm.execIntBinop(op)
            if e != nil { return e }
        case OP_EQ, OP_LT, OP_LE:
            e := vm.execConditon(op)
            if e != nil { return e }
        case OP_ARRAY:
            // elems already on stack
            nelems := int(vm.fetch2Bytes())
            arr := make([]Object, nelems)
            for i := 0; i < nelems; i++ {
                arr[nelems-i-1] = vm.stackPop()
            }
            e := vm.stackPush(ObjArr{arr})
            if e != nil { return e }
        case OP_ARRID:
            i := int(vm.fetch2Bytes())
            arr := vm.stackTop().(ObjArr)
            e := vm.stackPush(arr.Value[i])
            if e != nil { return e }
        case OP_JMP:
            vm.ip = int(vm.fetch2Bytes())
        case OP_JF:
            addr := int(vm.fetch2Bytes())
            t := vm.stackPop()
            if (falsey(t)) {
                vm.ip = addr - 1
            }
        case OP_LOAD:
            i := int(vm.fetchByte())
            e := vm.stackPush(vm.stack[vm.frame.bp-vm.frame.argc+i])
            if e != nil { return e }
        case OP_STORE:
            o := vm.stackPop()
            i := int(vm.fetchByte())
            vm.stack[vm.frame.bp-vm.frame.argc+i] = o
        case OP_GLOAD:
            i := int(vm.fetchByte())
            e := vm.stackPush(vm.globals[i])
            if e != nil { return e }
        case OP_GSTORE:
            o := vm.stackPop()
            i := int(vm.fetchByte())
            vm.globals[i] = o
        case OP_CALL:
            // args already on stack
            addr := int(vm.fetch2Bytes())
            argc := int(vm.fetchByte())
            vm.frame = FrameNew(vm.ip, argc, vm.sp, vm.frame)
            vm.ip = addr - 1
        case OP_RET:
            ret := vm.stackPop()
            vm.sp = vm.frame.bp
            for i := 0; i < vm.frame.argc; i++ {
                vm.stackPop()
            }
            vm.ip = vm.frame.ret
            vm.frame = vm.frame.caller
            e := vm.stackPush(ret)
            if e != nil { return e }
        case OP_PRINT:
            v := vm.stackTop()
            fmt.Println(v)
        default:
            return fmt.Errorf("error: unknown opcode %d", op)
        }
    }
    return nil
}

func (vm *VM) execIntBinop(op OpCode) error {
    var res ObjInt
    op2, op1 := vm.stackPop().(ObjInt), vm.stackPop().(ObjInt)
    switch (op) {
    case OP_IADD: res = ObjInt{op1.Value + op2.Value}
    case OP_ISUB: res = ObjInt{op1.Value - op2.Value}
    case OP_IMUL: res = ObjInt{op1.Value * op2.Value}
    case OP_IDIV: res = ObjInt{op1.Value / op2.Value}
    case OP_IREM: res = ObjInt{op1.Value % op2.Value}
    }
    e := vm.stackPush(res)
    return e
}

func (vm *VM) execConditon(op OpCode) error {
    var res ObjBool
    op2, op1 := vm.stackPop().(ObjInt), vm.stackPop().(ObjInt)
    switch (op) {
    case OP_EQ: res = ObjBool{op1.Value == op2.Value}
    case OP_LT: res = ObjBool{op1.Value < op2.Value}
    case OP_LE: res = ObjBool{op1.Value <= op2.Value}
    }
    e := vm.stackPush(res)
    return e
}

func falsey(o Object) bool {
    switch o.Type() {
    case ObjIntType: if o.(ObjInt).Value == 0 { return true }
    case ObjArrType: if o.(ObjArr).Value != nil { return true }
    case ObjBoolType: if o.(ObjBool).Value == false { return true }
    }
    return false
}

