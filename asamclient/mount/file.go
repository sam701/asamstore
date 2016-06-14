package mount

import (
	"log"
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

	lastPart        *schema.BytesPart
	lastPartContent []byte
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

func (f *file) loadBytePart(offset uint64) {
	if f.lastPart != nil && f.lastPart.Offset <= offset && offset < f.lastPart.Offset+f.lastPart.Size {
		return
	}
	for _, v := range f.parts {
		if offset >= v.Offset && offset < v.Offset+v.Size {
			f.lastPart = v
			f.lastPartContent = bsClient.Get(v.ContentRef)
			if f.lastPartContent == nil {
				log.Fatalln("Cannot get content", v.ContentRef)
			}
			return
		}
	}
	log.Fatalln("Cannot found byte part for offset", offset, "in", f.name)
}

func (f *file) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	f.loadBytePart(uint64(req.Offset))

	offset := int(req.Offset) - int(f.lastPart.Offset)
	end := offset + req.Size
	if end > len(f.lastPartContent) {
		end = len(f.lastPartContent)
	}

	resp.Data = f.lastPartContent[offset:end]

	return nil
}
