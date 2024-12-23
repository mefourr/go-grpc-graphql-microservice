package catalog

import (
	"context"
	"encoding/json"
	"errors"
	elastic "gopkg.in/olivere/elastic.v5"
	"log"
)

var (
	ErrNotFound = errors.New("catalog: not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductById(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip, take uint64) ([]Product, error)
	ListProductsWithIds(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error)
}

type ElasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (*ElasticRepository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &ElasticRepository{client: client}, nil
}

func (r *ElasticRepository) Close() {
	r.client.Stop()
}

func (r *ElasticRepository) PutProduct(ctx context.Context, p Product) error {
	_, err := r.client.Index().
		Index("catalog").
		Type("product").
		Id(p.ID).
		BodyJson(productDocument{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}).
		Do(ctx)
	return err
}

func (r *ElasticRepository) GetProductById(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get().
		Index("catalog").
		Type("product").
		Id(id).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	if !res.Found {
		return nil, ErrNotFound
	}

	var product productDocument
	if err = json.Unmarshal(*res.Source, &product); err != nil {
		return nil, err
	}

	return &Product{
		ID:          id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}, nil
}

func (r *ElasticRepository) ListProducts(ctx context.Context, skip, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []Product
	for _, hit := range res.Hits.Hits {
		var product productDocument
		if err = json.Unmarshal(*hit.Source, &product); err != nil {
			return nil, err
		}
		products = append(products, Product{
			ID:          hit.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return products, nil
}

func (r *ElasticRepository) ListProductsWithIds(ctx context.Context, ids []string) ([]Product, error) {
	var items []*elastic.MultiGetItem
	for _, id := range ids {
		items = append(
			items,
			elastic.NewMultiGetItem().
				Index("catalog").
				Type("product").
				Id(id),
		)
	}
	res, err := r.client.MultiGet().
		Add(items...).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []Product
	for _, doc := range res.Docs {
		var product productDocument
		if err := json.Unmarshal(*doc.Source, &product); err != nil {
			return nil, err
		}
		products = append(products, Product{
			ID:          doc.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return products, nil
}

func (r *ElasticRepository) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []Product
	for _, hit := range res.Hits.Hits {
		var product productDocument
		if err = json.Unmarshal(*hit.Source, &product); err != nil {
			return nil, err
		}
		products = append(products, Product{
			ID:          hit.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return products, nil
}
