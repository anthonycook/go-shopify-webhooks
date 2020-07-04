package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// ShopifyCustomer is a payload received from Shopify's customer create/update webhook
// See https://shopify.dev/docs/admin-api/rest/reference/events/webhook
type ShopifyCustomer struct {
	ID                  int64     `json:"id"`
	Email               string    `json:"email"`
	AcceptsMarketing    bool      `json:"accepts_marketing"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	OrdersCount         int64     `json:"orders_count"`
	State               string    `json:"state"`
	TotalSpent          string    `json:"total_spent"`
	LastOrderID         int64     `json:"last_order_id"`
	Note                string    `json:"note"`
	VerifiedEmail       bool      `json:"verified_email"`
	MultipassIdentifier string    `json:"multipass_identifier"`
	TaxExempt           bool      `json:"tax_exempt"`
	Phone               string    `json:"phone"`
	Tags                string    `json:"tags"`
	LastOrderName       string    `json:"last_order_name"`
	Currency            string    `json:"currency"`
	Addresses           []struct {
		ID           int64  `json:"id"`
		CustomerID   int64  `json:"customer_id"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Company      string `json:"company"`
		AddressOne   string `json:"address1"`
		AdressTwo    string `json:"address2"`
		City         string `json:"city"`
		Province     string `json:"province"`
		Country      string `json:"country"`
		Zip          string `json:"zip"`
		Phone        string `json:"phone"`
		ProvinceCode string `json:"province_code"`
		CountryCode  string `json:"country_code"`
		CountryName  string `json:"country_name"`
		Default      bool   `json:"default"`
	} `json:"addresses"`
}

// CustomerWebhook accepts a Shopify customer create/update webhook
// and creates or updates it in our database
func CustomerWebhook(c *gin.Context) {

	// Get JSON body
	customer := ShopifyCustomer{}

	err := c.BindJSON(&customer)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(400)
		return
	}

	// Create/update customers
	_, err = db.Exec(`
		INSERT INTO shopify_customers(
			id, 
			email, 
			accepts_marketing, 
			created_at, 
			updated_at, 
			first_name, 
			last_name, 
			orders_count, 
			state, 
			total_Spent, 
			last_order_id, 
			note, 
			verified_email, 
			multipass_identifier, 
			tax_exempt, 
			phone, 
			tags, 
			last_order_name, 
			currency
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		ON CONFLICT (id) DO UPDATE
		SET id = $1, 
				email = $2, 
				accepts_marketing = $3, 
				created_at = $4, 
				updated_at = $5, 
				first_name = $6, 
				last_name = $7, 
				orders_count = $8, 
				state = $9, 
				total_Spent = $10, 
				last_order_id = $11, 
				note = $12, 
				verified_email = $13, 
				multipass_identifier = $14, 
				tax_exempt = $15, 
				phone = $16, 
				tags = $17, 
				last_order_name = $18, 
				currency = $19;`,
		customer.ID,
		customer.Email,
		customer.AcceptsMarketing,
		customer.CreatedAt,
		customer.UpdatedAt,
		customer.FirstName,
		customer.LastName,
		customer.OrdersCount,
		customer.State,
		customer.TotalSpent,
		customer.LastOrderID,
		customer.Note,
		customer.VerifiedEmail,
		customer.MultipassIdentifier,
		customer.TaxExempt,
		customer.Phone,
		customer.Tags,
		customer.LastOrderName,
		customer.Currency,
	)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(400)
		return
	}

	// Create/update product images
	for _, address := range customer.Addresses {
		_, err = db.Exec(`
			INSERT INTO shopify_customer_addresses(
				id, 
				customer_id, 
				first_name, 
				last_name, 
				company, 
				address_one, 
				address_two, 
				city, 
				province, 
				country, 
				zip, 
				phone, 
				province_code, 
				country_code, 
				country_name, 
				default_address
			) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
			ON CONFLICT (id) DO UPDATE
			SET id = $1, 
					customer_id = $2, 
					first_name = $3, 
					last_name = $4, 
					company = $5, 
					address_one = $6, 
					address_two = $7, 
					city = $8, 
					province = $9, 
					country = $10, 
					zip = $11, 
					phone = $12, 
					province_code = $14, 
					country_code = $15, 
					default_address = $16;`,
			address.ID,
			address.CustomerID,
			address.FirstName,
			address.LastName,
			address.Company,
			address.AddressOne,
			address.AdressTwo,
			address.City,
			address.Province,
			address.Country,
			address.Zip,
			address.Phone,
			address.ProvinceCode,
			address.CountryCode,
			address.CountryName,
			address.Default,
		)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(400)
			return
		}

	}

	c.Status(200)
}
