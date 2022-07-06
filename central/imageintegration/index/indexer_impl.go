// Code generated by blevebindings generator. DO NOT EDIT.

package index

import (
	"bytes"
	bleve "github.com/blevesearch/bleve"
	mappings "github.com/stackrox/rox/central/imageintegration/mappings"
	metrics "github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	batcher "github.com/stackrox/rox/pkg/batcher"
	ops "github.com/stackrox/rox/pkg/metrics"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
	"time"
)

const batchSize = 5000

const resourceName = "ImageIntegration"

type indexerImpl struct {
	index bleve.Index
}

type imageIntegrationWrapper struct {
	*storage.ImageIntegration `json:"image_integration"`
	Type                      string `json:"type"`
}

func (b *indexerImpl) AddImageIntegration(imageintegration *storage.ImageIntegration) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Add, "ImageIntegration")
	if err := b.index.Index(imageintegration.GetId(), &imageIntegrationWrapper{
		ImageIntegration: imageintegration,
		Type:             v1.SearchCategory_IMAGE_INTEGRATIONS.String(),
	}); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) AddImageIntegrations(imageintegrations []*storage.ImageIntegration) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.AddMany, "ImageIntegration")
	batchManager := batcher.New(len(imageintegrations), batchSize)
	for {
		start, end, ok := batchManager.Next()
		if !ok {
			break
		}
		if err := b.processBatch(imageintegrations[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (b *indexerImpl) processBatch(imageintegrations []*storage.ImageIntegration) error {
	batch := b.index.NewBatch()
	for _, imageintegration := range imageintegrations {
		if err := batch.Index(imageintegration.GetId(), &imageIntegrationWrapper{
			ImageIntegration: imageintegration,
			Type:             v1.SearchCategory_IMAGE_INTEGRATIONS.String(),
		}); err != nil {
			return err
		}
	}
	return b.index.Batch(batch)
}

func (b *indexerImpl) Count(q *v1.Query, opts ...blevesearch.SearchOption) (int, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Count, "ImageIntegration")
	return blevesearch.RunCountRequest(v1.SearchCategory_IMAGE_INTEGRATIONS, q, b.index, mappings.OptionsMap, opts...)
}

func (b *indexerImpl) DeleteImageIntegration(id string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Remove, "ImageIntegration")
	if err := b.index.Delete(id); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) DeleteImageIntegrations(ids []string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.RemoveMany, "ImageIntegration")
	batch := b.index.NewBatch()
	for _, id := range ids {
		batch.Delete(id)
	}
	if err := b.index.Batch(batch); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) MarkInitialIndexingComplete() error {
	return b.index.SetInternal([]byte(resourceName), []byte("old"))
}

func (b *indexerImpl) NeedsInitialIndexing() (bool, error) {
	data, err := b.index.GetInternal([]byte(resourceName))
	if err != nil {
		return false, err
	}
	return !bytes.Equal([]byte("old"), data), nil
}

func (b *indexerImpl) Search(q *v1.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "ImageIntegration")
	return blevesearch.RunSearchRequest(v1.SearchCategory_IMAGE_INTEGRATIONS, q, b.index, mappings.OptionsMap, opts...)
}
