// Code generated by blevebindings generator. DO NOT EDIT.

package index

import (
	"bytes"
	"context"
	bleve "github.com/blevesearch/bleve"
	metrics "github.com/stackrox/rox/central/metrics"
	mappings "github.com/stackrox/rox/central/policycategory/index/mappings"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	batcher "github.com/stackrox/rox/pkg/batcher"
	ops "github.com/stackrox/rox/pkg/metrics"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
	"time"
)

const batchSize = 5000

const resourceName = "PolicyCategory"

type indexerImpl struct {
	index bleve.Index
}

type policyCategoryWrapper struct {
	*storage.PolicyCategory `json:"policy_category"`
	Type                    string `json:"type"`
}

func (b *indexerImpl) AddPolicyCategory(policycategory *storage.PolicyCategory) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Add, "PolicyCategory")
	if err := b.index.Index(policycategory.GetId(), &policyCategoryWrapper{
		PolicyCategory: policycategory,
		Type:           v1.SearchCategory_POLICY_CATEGORIES.String(),
	}); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) AddPolicyCategories(policycategories []*storage.PolicyCategory) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.AddMany, "PolicyCategory")
	batchManager := batcher.New(len(policycategories), batchSize)
	for {
		start, end, ok := batchManager.Next()
		if !ok {
			break
		}
		if err := b.processBatch(policycategories[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (b *indexerImpl) processBatch(policycategories []*storage.PolicyCategory) error {
	batch := b.index.NewBatch()
	for _, policycategory := range policycategories {
		if err := batch.Index(policycategory.GetId(), &policyCategoryWrapper{
			PolicyCategory: policycategory,
			Type:           v1.SearchCategory_POLICY_CATEGORIES.String(),
		}); err != nil {
			return err
		}
	}
	return b.index.Batch(batch)
}

func (b *indexerImpl) Count(ctx context.Context, q *v1.Query, opts ...blevesearch.SearchOption) (int, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Count, "PolicyCategory")
	return blevesearch.RunCountRequest(v1.SearchCategory_POLICY_CATEGORIES, q, b.index, mappings.OptionsMap, opts...)
}

func (b *indexerImpl) DeletePolicyCategory(id string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Remove, "PolicyCategory")
	if err := b.index.Delete(id); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) DeletePolicyCategories(ids []string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.RemoveMany, "PolicyCategory")
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

func (b *indexerImpl) Search(ctx context.Context, q *v1.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "PolicyCategory")
	return blevesearch.RunSearchRequest(v1.SearchCategory_POLICY_CATEGORIES, q, b.index, mappings.OptionsMap, opts...)
}
