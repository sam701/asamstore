package mount

import (
	"os"

	"bazil.org/fuse"
	"github.com/sam701/asamstore/asamclient/schema"
	"golang.org/x/net/context"
)

type file struct {
	name           string
	unixPermission os.FileMode
	parts          []*schema.BytesPart
}

func newFile(name string, permission string, parts []*schema.BytesPart) *file {
	return &file{
		name:           name,
		unixPermission: getFileMode(permission),
		parts:          parts,
	}
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Type() fuse.DirentType {
	return fuse.DT_File
}

func (f *file) size() uint64 {
	var s uint64
	for _, bp := range f.parts {
		s += bp.Size
	}
	return s
}

func (f *file) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Size = f.size()
	attr.Mode = f.unixPermission
	attr.Uid = userId
	attr.Gid = groupId
	return nil
}
