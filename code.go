// package govm
package main

type Instructions []byte
type OpCode byte

const (
    OP_HALT = iota
    OP_CONST // push constant

    // integer arithmetic ops
    OP_IADD
    OP_ISUB
    OP_IMUL
    OP_IDIV
    OP_IREM

    // arrays
    OP_ARRAY // create array from values on the stack
    OP_ARRID // get elem at index i

    // conditionals
    OP_EQ
    OP_LT
    OP_LE

    // branching
    OP_JMP
    OP_JF // jump if tos is falsey

    // load/store global variables
    OP_GLOAD
    OP_GSTORE

    // load/store local variables
    OP_LOAD
    OP_STORE
    
    // function call & return
    OP_CALL
    OP_RET

    OP_PRINT
)

