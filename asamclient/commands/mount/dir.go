package mount

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sam701/asamstore/asamclient/schema"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type node interface {
	fs.Node
	Name() string
	Type() fuse.DirentType
}

type dir struct {
	name           string
	unixPermission os.FileMode
	unixMTime      time.Time
	entries        []schema.BlobRef
	entriesSchemas []*schema.Schema
}

func (d *dir) Name() string {
	return d.name
}

func (d *dir) Type() fuse.DirentType {
	return fuse.DT_Dir
}

func (d *dir) getEntriesSchemas() []*schema.Schema {
	if d.entriesSchemas == nil {
		d.entriesSchemas = []*schema.Schema{}
		for _, ref := range d.entries {
			s := bsClient.GetSchema(ref)
			if s == nil {
				panic("cannot get blob ref:" + ref)
			}
			d.entriesSchemas = append(d.entriesSchemas, s)
		}
	}
	return d.entriesSchemas
}

func newDir(s *schema.Schema) *dir {
	t, err := time.Parse(time.RFC3339, s.UnixMtime)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return &dir{
		name:           s.FileName,
		unixPermission: getFileMode(s.UnixPermission),
		unixMTime:      t,
		entries:        s.DirEntries,
	}
}

func (d *dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Mode = d.unixPermission
	attr.Mtime = d.unixMTime
	attr.Uid = userId
	attr.Gid = groupId
	return nil
}

func (d *dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	res := []fuse.Dirent{}
	for _, s := range d.getEntriesSchemas() {
		res = append(res, fuse.Dirent{
			Name: s.FileName,
			Type: getDtType(s),
		})
	}
	return res, nil
}

func nodeFromSchema(s *schema.Schema) fs.Node {
	switch s.Type {
	case schema.ContentTypeDir:
		return newDir(s)
	case schema.ContentTypeFile:
		t, err := time.Parse(time.RFC3339, s.UnixMtime)
		if err != nil {
			log.Fatalln("ERROR", err)
		}

		return &file{
			name:           s.FileName,
			unixPermission: getFileMode(s.UnixPermission),
			unixMTime:      t,
			parts:          s.FileParts,
		}
	}
	return nil
}

func (d *dir) Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (fs.Node, error) {
	for _, s := range d.getEntriesSchemas() {
		if s.FileName == req.Name {
			return nodeFromSchema(s), nil
		}
	}
	return nil, fuse.ENOENT
}

func getDtType(s *schema.Schema) fuse.DirentType {
	switch s.Type {
	case schema.ContentTypeDir:
		return fuse.DT_Dir
	case schema.ContentTypeFile:
		return fuse.DT_File
	default:
		logUnmountAndExit("Unknown schema type:", s.Type)
	}
	return fuse.DT_Unknown
}

func getFileMode(unixPermission string) os.FileMode {
	i, err := strconv.ParseUint(unixPermission, 8, 32)
	if err != nil {
		logUnmountAndExit("ERROR", err)
	}
	return os.FileMode(i)
}
