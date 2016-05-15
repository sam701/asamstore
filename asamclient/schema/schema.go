package schema

type Schema struct {
	Version int         `json:"asamVersion"`
	Type    ContentType `json:"contentType"`

	// for a root
	RootName string `json:"rootName,omitempty"`

	// for a commit
	CommitTime string  `json:"commitTime,omitempty"`
	RootRef    BlobRef `json:"root,omitempty"`
	ContentRef BlobRef `json:"contentRef,omitempty"`

	FileName string `json:"fileName,omitempty"`

	UnixPermission string `json:"unixPermission,omitempty"`

	FileParts  []*BytesPart `json:"fileParts,omitempty"`
	DirEntries []BlobRef    `json:"dirEntries,omitempty"`
}

type BytesPart struct {
	Size       uint64  `json:"size"`
	Offset     uint64  `json:"offset"`
	ContentRef BlobRef `json:"contentRef"`
}

type ContentType string

const (
	ContentTypeFile   ContentType = "file"
	ContentTypeDir                = "dir"
	ContentTypeRoot               = "root"
	ContentTypeCommit             = "commit"
)

func NewSchema(contentType ContentType) *Schema {
	return &Schema{
		Version: 1,
		Type:    contentType,
	}
}
