# Shopify Webhooks

Shopify Webhooks is a Go library for receiving Shopify webhooks and saving them to a Postgres database.

## Installation

Setup environmental variables for your database connection

```bash
export DB_HOST=""
export DB_PORT=""
export DB_USER=""
export DB_PASS=""
export DB_NAME=""
```

Build the Go binary

```bash
go build
```

Run the Go binary

```bash
./shopify-hooks
```

Setup your webhooks in your Shopify admin and point them to the following endpoints:

**Product create/update: YOURDOMAIN/sync/product**

**Customer create/update: YOURDOMAIN/sync/customer**


## Important notes

This software is currently in beta, I plan to add integrity checks before the first release.

The plan is to update this repo to keep it inline with the latest Shopify API version, it's currently using the 2020-04 version.


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)