package cas

import (
	"context"
	"io"

	"github.com/buildbarn/bb-storage/pkg/blobstore"
	"github.com/buildbarn/bb-storage/pkg/blobstore/chunk"
	"github.com/buildbarn/bb-storage/pkg/capabilities"
	bb_cas "github.com/buildbarn/bb-storage/pkg/cas"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/buildbarn/bb-storage/pkg/zstd"
)

// PutReaderFunc stores the contents of r in the Content Addressable
// Storage (CAS) as blob d. Implementations may elide data that is
// already present in storage. This type hides how blobs are decomposed
// and stored, so that callers do not need to be aware of chunking.
type PutReaderFunc func(ctx context.Context, d digest.Digest, r io.Reader) error

// NewPutReaderFunc returns a PutReaderFunc that stores blobs by
// chunking them, as implemented by cas.PutReader.
func NewPutReaderFunc(zstdPool zstd.Pool, chunkStorage blobstore.BlobAccess[*chunk.Chunk], chunkMappingStorage blobstore.BlobAccess[chunk.Mapping], cdcParametersFetcher capabilities.CDCParametersFetcher) PutReaderFunc {
	return func(ctx context.Context, d digest.Digest, r io.Reader) error {
		params, err := cdcParametersFetcher.FetchCDCParameters(ctx, d.GetInstanceName())
		if err != nil {
			return err
		}
		return bb_cas.PutReader(ctx, zstdPool, chunkStorage, chunkMappingStorage, params, d, r)
	}
}
