package mount

import (
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/sam701/asamstore/asamclient/index"
	"github.com/sam701/asamstore/asamclient/schema"
	"golang.org/x/net/context"
)

type rootNode struct {
	name    string
	ref     schema.BlobRef
	commits []*index.Commit
}

func (r *rootNode) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Mode = os.ModeDir | 0700
	attr.Uid = userId
	attr.Gid = groupId
	return nil
}

func (r *rootNode) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	res := []fuse.Dirent{}
	if r.commits == nil {
		r.commits = bsIndex.GetCommits(r.ref)
	}
	for _, c := range r.commits {
		res = append(res, fuse.Dirent{
			Name: c.CommitTime,
			Type: fuse.DT_Dir,
		})
	}
	return res, nil
}

func (r *rootNode) Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (fs.Node, error) {
	for _, c := range r.commits {
		if c.CommitTime == req.Name {
			dirSchema := bsClient.GetSchema(c.Content)
			if dirSchema == nil {
				panic("No such blob: " + c.Content)
			}
			return newDir(c.CommitTime, dirSchema.UnixPermission, dirSchema.DirEntries), nil
		}
	}
	return nil, fuse.ENOENT
}
