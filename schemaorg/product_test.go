package schemaorg

import (
	"testing"
)

func TestNewProduct_SetsFieldsAndDefaults(t *testing.T) {
	brand := &Brand{Name: "Example Brand"}
	offer := &Offer{Price: "29.99", PriceCurrency: "USD"}
	agg := &AggregateRating{RatingValue: 4.5, ReviewCount: 10}
	review := &Review{ReviewBody: "Good product", Author: &Person{Name: "Alice"}, ReviewRating: &Rating{RatingValue: 5}}

	product := NewProduct(
		"Example Product",
		"Description",
		[]string{"https://example.com/image.jpg"},
		"SKU123",
		brand,
		offer,
		"Electronics",
		agg,
		[]*Review{review},
	)

	if product.Context != "https://schema.org" {
		t.Errorf("expected context schema.org, got %s", product.Context)
	}
	if product.Type != "Product" {
		t.Errorf("expected type Product, got %s", product.Type)
	}
	if product.Brand == nil || product.Brand.Type != "Brand" {
		t.Errorf("expected brand type Brand, got %v", product.Brand)
	}
	if product.Offers == nil || product.Offers.Type != "Offer" {
		t.Errorf("expected offer type Offer, got %v", product.Offers)
	}
	if product.AggregateRating == nil || product.AggregateRating.Type != "AggregateRating" {
		t.Errorf("expected aggregateRating type AggregateRating, got %v", product.AggregateRating)
	}
	if len(product.Review) != 1 || product.Review[0].Type != "Review" {
		t.Errorf("expected review type Review, got %v", product.Review)
	}
	if product.Review[0].Author == nil || product.Review[0].Author.Type != "Person" {
		t.Errorf("expected author type Person, got %v", product.Review[0].Author)
	}
	if product.Review[0].ReviewRating == nil || product.Review[0].ReviewRating.Type != "Rating" {
		t.Errorf("expected rating type Rating, got %v", product.Review[0].ReviewRating)
	}
}

func TestProduct_EnsureDefaults_NilSafe(t *testing.T) {
	product := &Product{}
	product.ensureDefaults()

	if product.Context != "https://schema.org" {
		t.Errorf("expected context schema.org, got %s", product.Context)
	}
	if product.Type != "Product" {
		t.Errorf("expected type Product, got %s", product.Type)
	}
}

func TestBrand_EnsureDefaults(t *testing.T) {
	b := &Brand{}
	b.ensureDefaults()
	if b.Type != "Brand" {
		t.Errorf("expected Brand type, got %s", b.Type)
	}
}

func TestOffer_EnsureDefaults(t *testing.T) {
	o := &Offer{}
	o.ensureDefaults()
	if o.Type != "Offer" {
		t.Errorf("expected Offer type, got %s", o.Type)
	}
}

func TestAggregateRating_EnsureDefaults(t *testing.T) {
	a := &AggregateRating{}
	a.ensureDefaults()
	if a.Type != "AggregateRating" {
		t.Errorf("expected AggregateRating type, got %s", a.Type)
	}
}

func TestReview_EnsureDefaults_WithNested(t *testing.T) {
	r := &Review{
		Author:       &Person{},
		ReviewRating: &Rating{},
	}
	r.ensureDefaults()

	if r.Type != "Review" {
		t.Errorf("expected Review type, got %s", r.Type)
	}
	if r.Author.Type != "Person" {
		t.Errorf("expected Author type Person, got %s", r.Author.Type)
	}
	if r.ReviewRating.Type != "Rating" {
		t.Errorf("expected Rating type Rating, got %s", r.ReviewRating.Type)
	}
}

func TestRating_EnsureDefaults(t *testing.T) {
	r := &Rating{}
	r.ensureDefaults()
	if r.Type != "Rating" {
		t.Errorf("expected Rating type, got %s", r.Type)
	}
}

func TestProduct_ToGoHTMLJsonLd(t *testing.T) {
	product := NewProduct(
		"Test Product",
		"Sample product description",
		[]string{"https://example.com/image.jpg"},
		"SKU123",
		&Brand{Name: "TestBrand"},
		&Offer{Price: "99.99", PriceCurrency: "USD"},
		"Electronics",
		nil,
		nil,
	)

	html, err := product.ToGoHTMLJsonLd()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if html == "" {
		t.Error("expected non-empty HTML output")
	}
}
