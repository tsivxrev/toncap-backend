package utils

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/nft"
)

var client *liteclient.ConnectionPool
var ctx context.Context
var api *ton.APIClient

func init() {
	client = liteclient.NewConnectionPool()
	err := client.AddConnectionsFromConfigUrl(
		context.Background(),
		"https://ton-blockchain.github.io/global.config.json",
	)
	if err != nil {
		log.Printf("Ton API: %s\n", err.Error())
	}

	ctx = client.StickyContext(context.Background())
	api = ton.NewAPIClient(client)
}

func JettonMeta(contract string) (response fiber.Map, err error) {
	jetton_addr := address.MustParseAddr(contract)
	master := jetton.NewJettonMasterClient(api, jetton_addr)

	data, err := master.GetJettonData(ctx)
	if err != nil {
		return nil, err
	}

	switch data.Content.(type) {
	case *nft.ContentOnchain:
		total_supply := data.TotalSupply.Int64()
		content := data.Content.(*nft.ContentOnchain)
		decimals_int, _ := strconv.Atoi(content.GetAttribute("decimals"))
		decimals := float64(decimals_int)

		return fiber.Map{
			"total_supply": float64(total_supply) / math.Pow(10, decimals),
			"decimals":     decimals,
			"name":         content.Name,
			"description":  content.Description,
			"image":        content.Image,
		}, nil

	case *nft.ContentOffchain:
		total_supply := data.TotalSupply.Uint64()
		content := data.Content.(*nft.ContentOffchain)
		content_url := strings.Replace(content.URI, "ipfs://", "https://ipfs.io/ipfs/", -1)

		res, err := http.Get(content_url)
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var ipfs_content map[string]interface{}
		err = json.Unmarshal(body, &ipfs_content)
		if err != nil {
			return nil, err
		}

		decimals, ok := ipfs_content["decimals"].(int)
		if !ok {
			ipfs_content["decimals"] = 9
			decimals = 9
		}

		ipfs_content["total_supply"] = float64(total_supply) / math.Pow(10, float64(decimals))

		return ipfs_content, nil
	}

	return fiber.Map{}, nil
}
