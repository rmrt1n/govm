package main

import (
    "fmt"
)

func main() {
    ins := Instructions{
        // test fn call
        // OP_CONST, 0,
        // OP_CALL, 0x00, 0x08, 1,
        // OP_PRINT,
        // OP_HALT,
        // // fnc
        // OP_CONST, 1,
        // OP_CONST, 2,
        // OP_STORE, 0,
        // OP_LOAD, 0,
        // OP_IADD,
        // OP_RET,

        // test array
        // OP_CONST, 0,
        // OP_CONST, 0,
        // OP_CONST, 1,
        // OP_CONST, 0,
        // OP_CONST, 0,
        // OP_ARRAY, 0x00, 0x05,
        // OP_PRINT,
        // OP_ARRID, 0x00, 0x02,
        // OP_PRINT,
        // OP_HALT,

        // test branching
        OP_CONST, 0,
        OP_CONST, 1,
        OP_EQ,
        OP_JF, 0x00, 0x0a,
        OP_CONST, 2,
        OP_CONST, 3,
        OP_PRINT,
    }
    vm := VMNew(ins)
    vm.constants[0] = ObjInt{10}
    vm.constants[1] = ObjInt{20}
    vm.constants[2] = ObjInt{30}
    vm.constants[3] = ObjInt{40}
    e := vm.Run()
    if e != nil {
        fmt.Println(e.Error())
    }
}
