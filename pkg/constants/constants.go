package constants

const (
	MIMEApplicationJSON = "application/json"
)
const (
	HeaderContentType = "Content-Type"
	HeaderAccept      = "Accept"
)

// Table name
const (
	TableNameCustomers = "customers"
	TableInvoiceItems  = "invoice_items"
	TableInvoices      = "invoices"
	TableItems         = "items"
)

// Date and time formats
const (
	DateTimeFormat = "2006-01-02 15:04:05" // Date and time format
	DateFormat     = "02/01/2006"
	DateFormatSQL  = "2006-01-02"
)

// Invoice Status

const (
	InvStatusUnpaid = "Unpaid"
	InvStatusPaid   = "Paid"
)
