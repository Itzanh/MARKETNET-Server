# MARKETNET

Marketnet is an un-deprecated ERP software, that will make your enterprise management easier. This software is based in a web application where all the user interaction happens, connected to a powerful backend with a PostgreSQL database.

Both server and client on this software are open source and free software, feel free to download the source code and add all those interesting features that your business needs and aren't inclued in this base software. Pull requests are also welcome.

The UI is a React web application, that you can host with an apache or nginx web server. When running the app, it will attempt to connect via WebSocket to the backend, on the same host as the webpage has loaded.

** CURRENTLY UNDER CONSRTUCTION **

[Official web page](https://www.marketnet.io/)

## Features

### Sales orders lifecycle

Get an exact status for all the sales orders, for the lifecycle and the invoicing.

### Traceability

Track the exact status of the sales order detail by detail.

### E-commerce integration

Powerful PrestShop integration that allows to syncronize all the business data between the e-commerce and the ERP.

### Ease to manage data

Quickly sort, filter, copy, import and export data from the webpage as if it was a native application.


## Installation

### Install PostgreSQL

Install PostgreSQL in your system if you haven't done it already. You can get PostgreSQL from the [official website](https://www.postgresql.org/download/).

### Import the database

Download the db.sql file in the repository, and import it into PostgreSQL.
Create a user account for the backend, and grant the permissions to access the database.

### Install the Go backend

Go to the [backend repository](https://github.com/Itzanh/MARKETNET-Server), download the code, and compile the Go application. Next, put the config.json file in the same directory as the binary, and edit the configuration file with your parameters.

Example of command to compile:
```
go build main.go sales_order.go address.go billing_series.go currency.go payment_method.go warehouse.go language.go country.go state.go customer.go product.go product_family.go color.go sales_order_detail.go stock.go sales_order_discount.go sales_invoice.go sales_invoice_detail.go manufacturing_order.go manufacturing_order_type.go packages.go packaging.go sales_order_detail_packaged.go warehouse_movement.go sales_delivery_note.go incoterms.go carrier.go shipping.go user.go group.go user_group.go login_token.go supplier.go purchase_order.go purchase_order_detail.go needs.go purchase_delivery_note.go purchase_invoice.go purchase_invoice_detail.go initial_data.go settings.go document_container.go document.go prestashop.go reports.go mail.go pallet.go connections.go journal.go account.go accounting_movement.go accounting_movement_detail.go config_accounts_vat.go collection_operation.go charges.go payment_transaction.go payment.go data_generator.go logs.go analytics.go api_rest.go api_key.go sendcloud.go shipping_tag.go woocommerce.go ecommerce.go connection_log.go shopify.go report_template.go saas.go google_authenticator.go transactional_log.go shipping_status_history.go sales_order_detail_digital_product_data.go email_log.go utils.go enterprise_logo.go manufacturing_order_type_components.go complex_manufacturing_order.go pos_terminals.go permission_dictionary.go crashreporter.go product_account.go vat_number_check.go report_translation.go hs_codes.go accounting_reports.go inventory.go webhook.go transfer_between_warehouses.go custom_fields.go orm_models.go label_printer_profile.go transfer_between_warehouses_minimum_stock.go product_included_products.go
```

Example of config.json file:
```
{
    "db": {
        "host": "localhost",
        "port": 5432,
        "user": "marketnet",
        "password": "** PASTE POSTGRESQL PASSWORD HERE **",
        "dbname": "marketnet"
    },
    "server": {
        "port": 12279,
        "hashIterations": 25000,
        "tokenExpirationHours": 48,
        "tls": {
            "useTLS": false,
            "crtPath": "./certificates/fullchain.pem",
            "keyPath": "./certificates/privkey.pem"
        }
    }
}
```

### Set up a web server and serve the front-end

Go to the [frontend repository](https://github.com/Itzanh/MARKETNET-Web), and download the build. Next, unzip the file in the root directory of your apache or nginx web server.

If you need to run the web application from a different path than the root of the web server, you must download the frontend code and build it yourself. First, update the homepage attribute of the package.json file, and set the relative path. Then, build it running the command:

```
npm run-script build
```

### Do the first steps

#### Configure SSL

It is highly recommended to set up SSL on the server. Copy the certificate and key file in a folder than can be accesed from the server, and edit the config.json file to useTLS=true.

Example for let's encrypt:

```
"tls": {
    "useTLS": true,
    "crtPath": "./certificates/fullchain.pem",
    "keyPath": "./certificates/privkey.pem"
}
```

#### Add users and groups

Open the web, and login with the default credentials:

```
User: marketnet
Password: admin1234
```

On the navigation bar, open "Utils" -> "Users" and change the password for the default user. This default user is automatically added to the Administrators group.

Create more groups with the according permissions and add more users to set up your production environment.

#### Enable PrestaShop integration

Go to the "Utils" menu, and click "Settings". On the window that will open, go to the "E-Commerce" tab, and type all the details for your PrestaShop integration.

#### Install utilities
##### MARKETNET Tag Printer

[MARKETNET Tag Printer](https://github.com/Itzanh/MARKETNET-Tag-printer) is a Windows utility that allows to print product tags. This application is called from a custom sheme registered from the windows registry.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

This code is distributed under [AGPL](https://spdx.org/licenses/AGPL-3.0-or-later.html) license.

```
.\main.exe --dev-no-upgrade-database --add-enterprise enterprise_key=PRODTEST enterprise_name=ProductionTest enterprise_desc=ProductionTest user_password=12345678 license_code=c312cff3-113a-4bdf-8dd3-c8da2094f7f1 license_chance=c312cff3-113a-4bdf-8dd3-c8da20
```

```
C:\Program Files\PostgreSQL\13\bin\pg_dump.exe --file "C:\\Users\\Itzan\\Desktop\\go\\MARKET~1\\db.sql" --host "localhost" --port "5432" --username "postgres" --no-password --format=p --schema-only --no-owner --no-privileges --no-comments "marketnet"
```
