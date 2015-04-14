package app

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/caiiiyua/csv2go"
	"github.com/revel/revel"
)

const Version = "0.1"

type Account struct {
	CardId     string    `csv:"card_id"`
	Name       string    `csv:"vip_name"`
	Mobile     string    `csv:"mobile"`
	CreateDate time.Time `csv:"vip_start_date"`
	SaveAmt    float64   `csv:"save_amt"`
	ConsumeAmt float64   `csv:"consum_amt"`
	LeftAmt    float64   `csv:"-"`
	// Items      []accountItem `csv:"-"`
}

func ImportAccounts() []Account {
	var acts []Account
	path := revel.BasePath
	accountCsv, err := os.Open(path + "/public/mssql2csv.csv")
	if err != nil {
		fmt.Println("open csv accountCsv failed", err)
		return nil
	}
	defer accountCsv.Close()

	d := csv2go.NewDecoder(accountCsv)
	d.Comma(';')

	for {
		act := &Account{}
		err := d.Decode(act)
		if err != nil {
			log.Println("decode from csv file failed")
			break
		}
		acts = append(acts, *act)
	}
	return acts
}
