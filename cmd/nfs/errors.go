package nfs

import (
	"errors"
	"fmt"
)

var (
	errEmptyExportString      = errors.New("export string is empty")
	errWrongExportID          = errors.New("error on parsing NFS export config string, ID is not a number")
	errWrongACLUID            = errors.New("UID in ACLs is not a number")
	errWrongACLGID            = errors.New("GID in ACLs is not a number")
	errWrongIdentityType      = errors.New("identity type is not one of [cidr hostname netgroup]")
	errWrongSquash            = errors.New("squash is not one of [none root rootuid all]")
	errWrongSquashAccessLevel = errors.New("access level is not one of [ro rw]")
)

type errInvalidExportConfigArg struct {
	got string
}

func (e *errInvalidExportConfigArg) Error() string {
	s := "invalid NFS export config argument, got %s, please use the following notations: \n"
	s += "exportString ::= ID,PATH,PSEUDOPATH,[ ACL [+ACL]... ]\n"
	s += "ACL ::= [cidr|hostname|netgroup];[MATCHER];[UID];[GID];[none|root|rootuid|all];[rw|ro]\n"
	return fmt.Sprintf(s, e.got)
}

func newErrInvalidExportConfigArg(got string) *errInvalidExportConfigArg {
	return &errInvalidExportConfigArg{
		got: got,
	}
}
