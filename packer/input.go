package packer

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

// Asset represents a single input source into the texture packer.
// Many assets are supplied and packed together to create a single atlas.
//
// Assets commonly represent files in a filesystem, but could also
// be blobs in a blobstore or images on a remote server.
type Asset interface {
	// Asset returns the name of the given asset
	Asset() string
	// CreateReader creates a readcloser capable of
	// reading an image from the asset source
	CreateReader() (io.ReadCloser, error)
}

// AssetStreamer is a factory responsible for creating readers that
// can read the contents of a given asset
type AssetStreamer interface {
	AssetStream(ctx context.Context) (<-chan Asset, <-chan error)
}

// AssetStreamerFunc is a function that conforms to the AssetStreamer interface
type AssetStreamerFunc func(ctx context.Context) (<-chan Asset, <-chan error)

// AssetStream implements the AssetStreamer interface
func (f AssetStreamerFunc) AssetStream(ctx context.Context) (<-chan Asset, <-chan error) {
	return f(ctx)
}

type FileAsset struct {
	Path string
}

func (a *FileAsset) Asset() string {
	return a.Path
}

func (a *FileAsset) CreateReader() (io.ReadCloser, error) {
	return os.Open(a.Path)
}

// NewFileStream creates an asset streamer that streams files from a given
// input directory. The input directory will be walked and readers will be
// created using the standard os package.
func NewFileStream(inputDirectory string) AssetStreamerFunc {
	return AssetStreamerFunc(func(ctx context.Context) (<-chan Asset, <-chan error) {
		stream := make(chan Asset)
		errc := make(chan error, 1)
		go func() {
			defer close(stream)
			// No select needed for this send, since errc is buffered.
			errc <- filepath.Walk(inputDirectory, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					// TODO walk nested directories...
					return nil
				}
				if !info.Mode().IsRegular() {
					return nil
				}
				select {
				case stream <- &FileAsset{Path: path}:
				case <-ctx.Done():
					return ctx.Err()
				}
				return nil
			})
		}()
		return stream, errc
	})
}

// TODO we could match globs too, that'd be cool
