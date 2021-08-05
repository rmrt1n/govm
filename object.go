// package govm
package main

type ObjType uint8

const (
    ObjIntType = iota
    // ObjStrType
    ObjBoolType
    // ObjFloatType
    ObjArrType
)

type Object interface {
    Type() ObjType
}

type ObjInt struct {
    Value int
}

type ObjArr struct {
    Value []Object
}

type ObjBool struct {
    Value bool
}

func (o ObjInt) Type() ObjType { return ObjIntType }
func (o ObjArr) Type() ObjType { return ObjArrType }
func (o ObjBool) Type() ObjType { return ObjBoolType }

