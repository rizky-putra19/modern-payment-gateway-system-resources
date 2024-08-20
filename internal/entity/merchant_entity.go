package entity

type AggregatedPaychannelEntity struct {
	Id               int    `db:"id" json:"id"`
	PaychannelCode   string `db:"merchant_paychannel_code" json:"paychannelCode"`
	PaychannelMethod string `db:"name" json:"paychannelMethod"`
	Merchant         string `db:"merchant_name" json:"merchant"`
}
