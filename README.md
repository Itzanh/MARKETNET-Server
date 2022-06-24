# MARKETNET

Marketnet is an un-deprecated ERP software, that will make your enterprise management easier. This software is based in a web application where all the user interaction happens, connected to a powerful backend with a PostgreSQL database.

Both server and client on this software are open source and free software, feel free to download the source code and add all those interesting features that your business needs and aren't inclued in this base software. Pull requests are also welcome.

The UI is a React web application, that you can host with an apache or nginx web server. When running the app, it will attempt to connect via WebSocket to the backend, on the same host as the webpage has loaded.

[Official web page](https://www.marketneterp.io/)

## Languages
- English
- Spanish

## Features

### Sales orders lifecycle

Get the exact status of all your orders, and billing and shipping information.

### Traceability

Easily browse through database record relationships.

Keep track of everything that has happened in a purchase or sale order without breaking your head.

### E-commerce integration

The powerful integration with several e-commerce softwares allows you to synchronize all your business data between e-commerce and ERP.

It couldn't be easier to manage the orders you receive online.

#### E-commerce platforms supported so far:
- PrestaShop
- WooCommerce
- Shopify

### Label printing

Functionality to print the tags with a utility from MARKETNET, what allows to print the labels with barcode for the products.

This tools algo generated the EAN13 barcodes of the products.

### Process automation

All the processes of the administration of your company are automated and with error control.

MARKETNET organizes work in your company automatically with its simplified workflow.

### Ready for day to day

All tools in the software package are tested, robustly designed, and tested by developers.

### Report generator
Easily display documents to your customers and suppliers using the reporting functionality. These reports are easily customizable, so that you can adapt them to the image of your company without complications.

Customize templates and deliver business-style reports to your clients, without leaving MARKETNET.


### Good documentation

Your users will have no problem using your application with the official MARKETNET documentation.

It is very easy to use the MARKETNET REST API, as well as to install and update the server, thanks to the good documentation of the project.

### Free and open source software

MARKETNET is free and open source software, and it respects it's users' freedoms.

This software is distributed under the AGPL license. [GNU AGPL v3.0-only](https://spdx.org/licenses/AGPL-3.0-only.html) 

### E-mail integration included (SendGrid and SMTP)

MARKETNET allows the final user to send emails directly from the server side using the web.

Forget about sending emails with the local Outlook/Thunderbird, MARKETNET allows the email sending centralized and completely transparent to the final users.

### Shipping integration

MARKETNET is fully integrated with the SendCloud infrastructure.

SendCloud is a platform that allows you to integrate with several transport agencies at the same time.



## Why MARKETNET?

Because this software is capable of automating all the management of your enterprise, without worrying. MARKETNET is designed to save you working time on the administration, from the first minute.

Some management programs are really expensive, require from installation and manual activation, you have to make youself the backups, manage your licenses, you need to install a local server on your enterpise, they dont' respect you freedoms, they don't let you modify the software... you'll never have this problems with MARKETNET.



## Useful links

[GitLab MARKETNET Server](https://gitlab.com/itzanh/marketnet-server)

[GitLab MARKETNET Web](https://gitlab.com/itzanh/marketnet-web)

[Project's home page](https://marketneterp.io/)

[Documentation: Wiki](https://marketneterp.io:3443/)



## Installation

Go to the wiki to lean [how to install MARKETNET](https://marketneterp.io:3443/en/technical-documentation/installing-marketnet) in your server.



## Compilation

Command to compile:
```
go build main.go sales_order.go address.go billing_series.go currency.go payment_method.go warehouse.go language.go country.go state.go customer.go product.go product_family.go color.go sales_order_detail.go stock.go sales_order_discount.go sales_invoice.go sales_invoice_detail.go manufacturing_order.go manufacturing_order_type.go packages.go packaging.go sales_order_detail_packaged.go warehouse_movement.go sales_delivery_note.go incoterms.go carrier.go shipping.go user.go group.go user_group.go login_token.go supplier.go purchase_order.go purchase_order_detail.go needs.go purchase_delivery_note.go purchase_invoice.go purchase_invoice_detail.go initial_data.go settings.go document_container.go document.go prestashop.go reports.go mail.go pallet.go connections.go journal.go account.go accounting_movement.go accounting_movement_detail.go config_accounts_vat.go collection_operation.go charges.go payment_transaction.go payment.go data_generator.go logs.go analytics.go api_rest.go api_key.go sendcloud.go shipping_tag.go woocommerce.go ecommerce.go connection_log.go shopify.go report_template.go saas.go google_authenticator.go transactional_log.go shipping_status_history.go sales_order_detail_digital_product_data.go email_log.go utils.go enterprise_logo.go manufacturing_order_type_components.go complex_manufacturing_order.go pos_terminals.go permission_dictionary.go product_account.go vat_number_check.go report_translation.go hs_codes.go accounting_reports.go inventory.go webhook.go transfer_between_warehouses.go custom_fields.go orm_models.go label_printer_profile.go transfer_between_warehouses_minimum_stock.go product_included_products.go
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

This code is distributed under the [AGPL](https://spdx.org/licenses/AGPL-3.0-only.html) license (AGPL-3.0-only).

You can find the full text of this license in the `COPYING` file.

