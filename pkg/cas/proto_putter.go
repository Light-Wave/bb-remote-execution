package cas

import (
	"bytes"
	"context"

	"github.com/buildbarn/bb-storage/pkg/digest"

	"google.golang.org/protobuf/proto"
)

// ProtoPutter stores Protobuf messages in the Content Addressable
// Storage (CAS). The digest of the stored message is returned, so that
// the object may be referenced.
type ProtoPutter interface {
	PutProto(ctx context.Context, message proto.Message, digestFunction digest.Function) (digest.Digest, error)
}

// casProtoPutter is a ProtoPutter that stores messages as blobs in the
// CAS.
type casProtoPutter struct {
	putReader PutReaderFunc
}

// NewProtoPutter returns a ProtoPutter that stores messages in the CAS
// using the provided PutReaderFunc.
func NewProtoPutter(putReader PutReaderFunc) ProtoPutter {
	return &casProtoPutter{
		putReader: putReader,
	}
}

func (p *casProtoPutter) PutProto(ctx context.Context, message proto.Message, digestFunction digest.Function) (digest.Digest, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return digest.BadDigest, err
	}
	digestGenerator := digestFunction.NewGenerator(int64(len(data)))
	if _, err := digestGenerator.Write(data); err != nil {
		return digest.BadDigest, err
	}
	blobDigest := digestGenerator.Sum()
	if err := p.putReader(ctx, blobDigest, bytes.NewReader(data)); err != nil {
		return digest.BadDigest, err
	}
	return blobDigest, nil
}
