package structs

import (
	"time"
)

type Incoming_Order struct {
	ID              uint       `gorm:"primaryKey;autoIncrement;column:id"`
	IDPo            string     `gorm:"size:30;column:id_po;default:null"`
	PoDate          time.Time  `gorm:"column:po_date;default:null"`
	Tgljt           time.Time  `gorm:"column:tgljt;default:null"`
	IDSupplier      int64      `gorm:"column:id_supplier;default:null"`
	SupplierName    string     `gorm:"size:255;column:supplier_name;default:null"`
	IDCurrency      int        `gorm:"column:id_currency;default:null"`
	SupplierAddress string     `gorm:"size:255;column:supplier_address;default:null"`
	PaymentTerm     string     `gorm:"size:100;column:payment_term;default:null"`
	PoNumber        string     `gorm:"size:100;column:po_number;default:null"`
	TotalQty        float32    `gorm:"column:total_qty;default:null"`
	TotalPrice      float32    `gorm:"type:decimal(15,2);column:total_price;default:null"`
	TotalDisc       float32    `gorm:"type:decimal(15,2);column:total_disc;default:null"`
	TotalPpn        float32    `gorm:"type:decimal(15,2);column:total_ppn;default:null"`
	GrandTotal      float32    `gorm:"type:decimal(15,2);column:grand_total;default:null"`
	State           string     `gorm:"size:50;column:state;default:null"`
	Catatan         string     `gorm:"type:text;column:catatan"`
	NoCar           string     `gorm:"size:50;column:nocar;default:null"`
	SuratJalan      string     `gorm:"size:255;column:suratjalan;default:null"`
	From            string     `gorm:"size:50;column:from;default:null"`
	Import          string     `gorm:"size:255;column:import;default:null"`
	Posting         string     `gorm:"size:255;column:posting;default:null"`
	TypePpn         string     `gorm:"size:50;column:type_ppn;default:null"`
	PpnPercent      float32    `gorm:"type:decimal(5,2);column:ppn_percent;default:null"`
	DiscPercent     float32    `gorm:"type:decimal(5,2);column:disc_percent;default:null"`
	Invoice         string     `gorm:"size:100;column:invoice;default:null"`
	DateInvoice     *time.Time `gorm:"column:date_invoice;default:null"`
	Remarks         string     `gorm:"type:text;column:remakrs"`
	TypePayment     string     `gorm:"size:100;column:typepayment;default:null"`
	SoNum           string     `gorm:"size:100;column:sonum;default:null"`
	DateWr          *time.Time `gorm:"column:date_wr;default:null"`
}
type Incoming_order_detail struct {
	IDDetail      int64     `gorm:"primaryKey;autoIncrement" json:"id_detail"`
	IDFromPo      string    `gorm:"size:40" json:"id_from_po"`
	IDInv         string    `gorm:"size:30" json:"id_inv"`
	IDPo          string    `gorm:"size:30" json:"id_po"`
	PoDate        time.Time `gorm:"type:date" json:"po_date"`
	DateReceipt   time.Time `gorm:"type:date" json:"date_receipt"`
	Ponum         string    `gorm:"size:100" json:"ponum"`
	IDSupplier    int64     `json:"id_supplier"`
	SupplierName  string    `gorm:"size:255" json:"supplier_name"`
	PN            string    `gorm:"size:100" json:"PN"`
	PNas          string    `gorm:"size:100" json:"PNas"`
	PNAME         string    `gorm:"size:255" json:"PNAME"`
	PNAMEas       string    `gorm:"size:255" json:"PNAMEas"`
	Substitute    string    `gorm:"size:100" json:"subtitute"`
	Qty           float32   `json:"qty"`
	QtyOrder      float32   `json:"qtyorder"`
	QtyIn         float32   `json:"qtyin"`
	QtyReal       float32   `json:"qtyreal"`
	PriceUSD      float64   `gorm:"type:decimal(15,2)" json:"priceusd"`
	IncomingDate  time.Time `gorm:"type:date" json:"incomingDate"`
	Stock         int       `json:"stok"`
	Kurs          float64   `gorm:"type:decimal(15,2)" json:"kurs"`
	Disc          float64   `gorm:"type:decimal(10,2)" json:"disc"`
	Price         float64   `gorm:"type:decimal(15,2)" json:"price"`
	Subtotal      float64   `gorm:"type:decimal(15,2)" json:"subtotal"`
	Foto          string    `gorm:"size:255" json:"foto"`
	State         string    `gorm:"size:50" json:"state"`
	Posting       string    `gorm:"size:30" json:"posting"`
	DateInvoice   time.Time `gorm:"type:date" json:"date_invoice"`
	Invoice       string    `gorm:"size:100" json:"invoice"`
	From          string    `gorm:"size:50" json:"from"`
	Brand         string    `gorm:"size:100" json:"brand"`
	User          string    `gorm:"size:100" json:"user"`
	InvName       string    `gorm:"size:100" json:"invname"`
	QtyInv        int       `json:"qty_inv"`
	StatusInvoice string    `gorm:"size:50" json:"statusInvoice"`
	QtyInvFix     int       `json:"qtyinvfix"`
}
