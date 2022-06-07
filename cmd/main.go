package main

import (
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	"github.com/mennanov/fieldmask-utils/testproto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// A function that maps field mask field names to the names used in Go structs.
// It has to be implemented according to your needs.
func naming(s string) string {
	if s == "foo" {
		return "Foo"
	}
	return s
}

func main() {
	var request = testproto.UpdateUserRequest{
		User:      &testproto.User{},
		FieldMask: &fieldmaskpb.FieldMask{},
	}
	userDst := &testproto.User{} // a struct to copy to
	mask, _ := fieldmask_utils.MaskFromPaths(request.FieldMask.Paths, naming)
	fieldmask_utils.StructToStruct(mask, request.User, userDst)
	// Only the fields mentioned in the field mask will be copied to userDst, other fields are left intact
}
