package app

import (
	"fmt"
	"os"
	"time"

	"github.com/caiiiyua/csv2go"
	"github.com/revel/revel"
)

const Version = "0.1"

type Account struct {
	CardId     string        `csv:"card_id"`
	Name       string        `csv:"vip_name"`
	Mobile     string        `csv:"mobile"`
	CreateDate time.Time     `csv:"vip_start_date"`
	SaveAmt    float64       `csv:"save_amt"`
	ConsumeAmt float64       `csv:"consum_amt"`
	LeftAmt    float64       `csv:"-"`
	Items      []accountItem `csv:"-"`
}

type accountItem struct {
	ActId    string  `csv:"card_id"`
	ItemId   string  `csv:"item_no"`
	Quantity int64   `csv:"real_qty"`
	Price    float64 `csv:"sale_price"`
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

	itemCsv, err := os.Open(revel.BasePath + "/public/mssql2csv2.csv")
	if err != nil {
		fmt.Println("open csv itemCsv failed")
		return nil
	}
	defer itemCsv.Close()

	it := csv2go.NewDecoder(itemCsv)
	it.Comma(';')

	var tmpItems []accountItem

	for {
		act := &Account{}
		err := d.Decode(act)
		if err != nil {
			break
		}
		if len(tmpItems) > 0 {
			act.Items = append(act.Items, tmpItems[0])
			tmpItems = tmpItems[:0]
		}
		for {
			item := &accountItem{ActId: act.CardId}
			err = it.Decode(item)
			if err != nil {
				break
			}
			if item.ActId != act.CardId {
				tmpItems = append(tmpItems, *item)
				break
			}
			act.Items = append(act.Items, *item)
		}

		act.LeftAmt = act.SaveAmt - act.LeftAmt
		acts = append(acts, *act)
	}
	return acts
}
