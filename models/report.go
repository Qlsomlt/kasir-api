package models

type SalesReport struct {
	TotalRevenue   int            `json:"total_revenue"`
	TotalTransaksi int            `json:"total_transaksi"`
	ProdukTerlaris BestSellerItem `json:"produk_terlaris"`
}

type BestSellerItem struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}
