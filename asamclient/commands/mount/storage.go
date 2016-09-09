package mount

import (
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type storageNode struct{}

func (s *storageNode) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Mode = os.ModeDir | 0700
	attr.Uid = userId
	attr.Gid = groupId
	return nil
}

func (s *storageNode) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	res := []fuse.Dirent{}
	for name, _ := range bsIndex.GetRoots() {
		res = append(res, fuse.Dirent{
			Name: name,
			Type: fuse.DT_Dir,
		})
	}
	return res, nil
}

func (s *storageNode) Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (fs.Node, error) {
	for name, ref := range bsIndex.GetRoots() {
		if name == req.Name {
			return &rootNode{
				name: name,
				ref:  ref,
			}, nil
		}
	}
	return nil, fuse.ENOENT
}
