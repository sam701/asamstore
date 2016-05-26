package mount

import (
	"bytes"
	"os"
	"time"

	"bazil.org/fuse"
	"github.com/sam701/asamstore/asamclient/schema"
	"golang.org/x/net/context"
)

type file struct {
	name           string
	unixPermission os.FileMode
	unixMTime      time.Time
	parts          []*schema.BytesPart
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
	attr.Mtime = f.unixMTime
	attr.Uid = userId
	attr.Gid = groupId
	return nil
}

func (f *file) ReadAll(ctx context.Context) ([]byte, error) {
	var buf bytes.Buffer
	for _, part := range f.parts {
		if ok := bsClient.Get(part.ContentRef, &buf); !ok {
			logUnmountAndExit("Cannot read blob", part.ContentRef)
		}
	}
	return buf.Bytes(), nil
}
