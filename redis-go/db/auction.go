package db

type Auction struct {
	Command     string `json:"command" binding:"required"`
	AuctionID   int    `json:"auctionid" binding:"required"`
	ServerIP    string `json:"ipaddr"`
	Description string `json:"desc"`
	ConnUsers   int    `json:"connected users"`
}
