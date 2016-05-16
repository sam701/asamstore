package mount

import (
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/sam701/asamstore/asamclient/schema"
	"golang.org/x/net/context"
)

type commit struct {
	contentRef    schema.BlobRef
	contentSchema *schema.Schema
}

func (c *commit) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Mode = os.ModeDir | 0700
	attr.Uid = userId
	attr.Gid = groupId
	return nil
}

func (c *commit) getSchema() *schema.Schema {
	if c.contentSchema == nil {
		c.contentSchema = bsClient.GetSchema(c.contentRef)
		if c.contentSchema == nil {
			logUnmountAndExit("Schema does not exist", c.contentRef)
		}
	}
	return c.contentSchema
}

func (c *commit) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	s := c.getSchema()
	res := []fuse.Dirent{fuse.Dirent{
		Name: s.FileName,
		Type: getDtType(s),
	}}
	return res, nil
}

func (c *commit) Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (fs.Node, error) {
	s := c.getSchema()
	if req.Name == s.FileName {
		return nodeFromSchema(s), nil
	}
	return nil, fuse.ENOENT
}
