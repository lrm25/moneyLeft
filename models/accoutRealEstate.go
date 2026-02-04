package models

import (
	"fmt"
	"github.com/lrm25/moneyLeft/logger"
)

type RealEstateConfig struct {
	Person              *Person
	Name                string
	Amount              float64
	AppreciationPerYear float64
	PercentSaleCosts    int
	PercentCash         int
	PercentBonds        int
	PercentStock        int
	StockReturn         float64
	StockSaleFee        float64
	BondReturn          float64
}

// RealEstate represents a real estate investment that can be sold
type RealEstate struct {
	closed bool
	config *RealEstateConfig
}

type RealEstateInvestments []*RealEstate

func NewRealEstate(config *RealEstateConfig) *RealEstate {
	totalPercent := config.PercentCash + config.PercentSaleCosts + config.PercentStock + config.PercentBonds
	if totalPercent != 100 {
		panic("real estate costs and percentages on sale must add up to 100 (currently " + fmt.Sprintf("%d", totalPercent) + ")")
	}
	return &RealEstate{
		config: config,
		closed: false,
	}
}

func (r *RealEstate) Amount() float64 {
	return r.config.Amount
}

func (r *RealEstate) Name() string {
	return r.config.Name
}

func (r *RealEstate) Type() int {
	return TypeRealEstate
}

func (r *RealEstate) Closed() bool {
	return r.closed
}

func (r *RealEstate) Removable() bool {
	return true
}

func (r *RealEstate) Close() {
	r.closed = true
}

func (r *RealEstate) Increase() {
	r.config.Amount *= 1 + (r.config.AppreciationPerYear / 1200.0)
}

func (r *RealEstate) Deduct(amount float64) (float64, float64) {
	var cash, stock, bonds float64
	if 0 < r.config.PercentCash {
		cash = float64(r.config.PercentCash) * r.config.Amount / 100
		c := NewCashAccount(r.config.Name+" sale cash", TypeSavingsNoInterest, cash)
		r.config.Person.accounts = append(r.config.Person.accounts, c)
	}
	if 0 < r.config.PercentStock {
		stock = float64(r.config.PercentStock) * r.config.Amount / 100
		c := NewAccountStockBrokerage(r.config.Name+" sale stock", stock, r.config.StockReturn, r.config.StockSaleFee, r.config.Person)
		r.config.Person.interestAccounts = append(r.config.Person.interestAccounts, c)
	}
	if 0 < r.config.PercentBonds {
		bonds = float64(r.config.PercentBonds) * r.config.Amount / 100
		c := NewAccountStockBrokerage(r.config.Name+" sale bonds", bonds, r.config.BondReturn, r.config.StockSaleFee, r.config.Person)
		r.config.Person.interestAccounts = append(r.config.Person.interestAccounts, c)
	}
	r.config.Amount = 0
	r.closed = true
	logger.Get().Debug(fmt.Sprintf("Real estate sale, %.2f outstanding", amount))
	return 0, amount
}

func (r *RealEstate) String() string {
	return fmt.Sprintf("Name: %s, accountType: %d, amount: %.2f, closed: %t, removable: %t", r.config.Name, TypeRealEstate, r.config.Amount, r.closed, false)
}
