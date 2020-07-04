package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// ShopifyProduct is a payload received from Shopify's product create/update webhook
// See https://shopify.dev/docs/admin-api/rest/reference/events/webhook
type ShopifyProduct struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	BodyHTML    string    `json:"body_html"`
	Vendor      string    `json:"vendor"`
	ProductType string    `json:"product_type"`
	Handle      string    `json:"handle"`
	Tags        string    `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at"`
	Variants    []struct {
		ID                int64     `json:"id"`
		ProductID         int64     `json:"product_id"`
		Title             string    `json:"title"`
		Price             string    `json:"price"`
		SKU               string    `json:"sku"`
		Position          int       `json:"position"`
		InventoryPolicy   string    `json:"inventory_policy"`
		CompareAtPrice    string    `json:"compare_at_price"`
		Taxable           bool      `json:"taxable"`
		Barcode           string    `json:"barcode"`
		Weight            float64   `json:"weight"`
		WeightUnit        string    `json:"weight_unit"`
		InventoryQuantity int       `json:"inventory_quantity"`
		RequiresShipping  bool      `json:"requires_shipping"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
	} `json:"variants"`
	Images []struct {
		ID         int64     `json:"id"`
		ProductID  int64     `json:"product_id"`
		Position   int       `json:"position"`
		Src        string    `json:"src"`
		Alt        string    `json:"alt"`
		Width      int       `json:"width"`
		Height     int       `json:"height"`
		VariantIDs []int64   `json:"variant_ids"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	} `json:"images"`
}

// ProductWebhook accepts a Shopify product create/update webhook
// and creates or updates it in our database
func ProductWebhook(c *gin.Context) {

	// Get JSON body
	product := ShopifyProduct{}

	err := c.BindJSON(&product)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(400)
		return
	}

	// Create/update product
	_, err = db.Exec(`
		INSERT INTO shopify_products(
			id, 
			title, 
			body_html, 
			vendor, 
			product_type, 
			handle, 
			tags, 
			created_at, 
			updated_at, 
			published_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE
		SET id = $1, 
				title = $2, 
				body_html = $3, 
				vendor = $4, 
				product_type = $5, 
				handle = $6, 
				tags = $7, 
				created_at = $8, 
				updated_at = $9, 
				published_at = $10;`,
		product.ID,
		product.Title,
		product.BodyHTML,
		product.Vendor,
		product.ProductType,
		product.Handle,
		product.Tags,
		product.CreatedAt,
		product.UpdatedAt,
		product.PublishedAt,
	)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(400)
		return
	}

	// Create/update product variants
	for _, variant := range product.Variants {

		price, _ := strconv.ParseFloat(variant.Price, 64)
		comparePrice, _ := strconv.ParseFloat(variant.CompareAtPrice, 64)

		_, err = db.Exec(`
			INSERT INTO shopify_variants(
				id, 
				product_id, 
				title, 
				price, 
				sku, 
				position, 
				inventory_policy, 
				compare_at_price, 
				taxable, 
				barcode, 
				weight, 
				weight_unit, 
				inventory_quantity, 
				requires_shipping, 
				created_at, 
				updated_at
			) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
			ON CONFLICT (id) DO UPDATE
			SET id = $1, 
					product_id = $2, 
					title = $3, 
					price = $4, 
					sku = $5, 
					position = $6, 
					inventory_policy = $7, 
					compare_at_price = $8, 
					taxable = $9, 
					barcode = $10, 
					weight = $11, 
					weight_unit = $12, 
					inventory_quantity = $13, 
					requires_shipping = $14, 
					created_at = $15, 
					updated_at = $16;`,
			variant.ID,
			variant.ProductID,
			variant.Title,
			price,
			variant.SKU,
			variant.Position,
			variant.InventoryPolicy,
			comparePrice,
			variant.Taxable,
			variant.Barcode,
			variant.Weight,
			variant.WeightUnit,
			variant.InventoryQuantity,
			variant.RequiresShipping,
			variant.CreatedAt,
			variant.UpdatedAt,
		)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(400)
			return
		}
	}

	// Create/update product images
	for _, image := range product.Images {
		_, err = db.Exec(`
		INSERT INTO shopify_images(
			id, 
			product_id, 
			position, 
			src, 
			alt, 
			width, 
			height, 
			variant_ids, 
			created_at, 
			updated_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE
		SET id = $1, 
				product_id = $2, 
				position = $3, 
				src = $4, 
				alt = $5, 
				width = $6, 
				height = $7, 
				variant_ids = $8, 
				created_at = $9, 
				updated_at = $10;`,
			image.ID,
			image.ProductID,
			image.Position,
			image.Src,
			image.Alt,
			image.Width,
			image.Height,
			pq.Array(image.VariantIDs),
			image.CreatedAt,
			image.UpdatedAt,
		)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(400)
			return
		}
	}

	c.Status(200)
}
