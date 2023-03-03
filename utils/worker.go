package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"toncap-backend/types"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/nft"
)

var worker_host = fmt.Sprintf("http://%s", os.Getenv("WORKER_HOST"))
var lite_client *liteclient.ConnectionPool
var ctx context.Context
var ton_api *ton.APIClient

func init() {
	lite_client = liteclient.NewConnectionPool()
	err := lite_client.AddConnectionsFromConfigUrl(
		context.Background(),
		"https://ton-blockchain.github.io/global.config.json",
	)
	if err != nil {
		log.Printf("Ton API: %s\n", err.Error())
	}

	ctx = lite_client.StickyContext(context.Background())
	ton_api = ton.NewAPIClient(lite_client)
}

func GetContractMeta(contract string) (contractMeta types.ContractMeta, err error) {
	master := jetton.NewJettonMasterClient(ton_api, address.MustParseAddr(contract))

	jetton_data, err := master.GetJettonData(ctx)
	if err != nil {
		return types.ContractMeta{}, err
	}

	socials, _ := GetSocials(contract)

	switch jetton_data.Content.(type) {
	case *nft.ContentOnchain:
		total_supply := jetton_data.TotalSupply.Int64()
		content := jetton_data.Content.(*nft.ContentOnchain)
		decimals, _ := strconv.Atoi(content.GetAttribute("decimals"))

		return types.ContractMeta{
			Name:        content.Name,
			Description: content.Description,
			Symbol:      content.GetAttribute("symbol"),
			Image:       content.Image,
			Socials:     socials,
			TotalSupply: float64(total_supply) / math.Pow(10, float64(decimals)),
			Decimals:    decimals,
		}, nil

	case *nft.ContentOffchain:
		total_supply := jetton_data.TotalSupply.Uint64()
		content := jetton_data.Content.(*nft.ContentOffchain)
		content_url := strings.Replace(content.URI, "ipfs://", "https://ipfs.io/ipfs/", -1)

		res, err := http.Get(content_url)
		if err != nil {
			return types.ContractMeta{}, err
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return types.ContractMeta{}, err
		}

		var ipfs_content map[string]any
		err = json.Unmarshal(body, &ipfs_content)
		if err != nil {
			return types.ContractMeta{}, err
		}

		decimals, ok := ipfs_content["decimals"].(int)
		if !ok {
			ipfs_content["decimals"] = 9
			decimals = 9
		}

		name, ok := ipfs_content["name"].(string)
		if !ok {
			name = ""
		}
		description, ok := ipfs_content["description"].(string)
		if !ok {
			description = ""
		}
		symbol, ok := ipfs_content["symbol"].(string)
		if !ok {
			symbol = ""
		}
		image, ok := ipfs_content["image"].(string)
		if !ok {
			image = ""
		}

		return types.ContractMeta{
			Name:        name,
			Description: description,
			Symbol:      symbol,
			Image:       image,
			Socials:     socials,
			TotalSupply: float64(total_supply) / math.Pow(10, float64(decimals)),
			Decimals:    decimals,
		}, nil
	}

	return types.ContractMeta{}, nil
}

func GetSocials(contract string) (socials types.SocialsResponse, err error) {
	socials = make(types.SocialsResponse)

	response, err := http.Get(worker_host + "/social/" + contract)
	if err != nil {
		return socials, err
	}

	responseBodyRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return socials, err
	}
	defer response.Body.Close()

	err = json.Unmarshal(responseBodyRaw, &socials)
	if err != nil {
		return socials, err
	}

	return socials, nil
}

func GetListedContracts() (listedContracts types.ListedContracts, err error) {
	response, err := http.Get(worker_host + "/listed")
	if err != nil {
		return nil, err
	}

	responseBodyRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var listedContractsResponse types.ListedContractsResponse
	err = json.Unmarshal(responseBodyRaw, &listedContractsResponse)
	if err != nil {
		return nil, err
	}

	return listedContractsResponse["data"], nil
}

func GetActual(contract string) (actual types.ActualResponse, err error) {
	response, err := http.Get(worker_host + "/markets/" + contract)
	if err != nil {
		return types.ActualResponse{}, err
	}

	responseBodyRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return types.ActualResponse{}, err
	}
	defer response.Body.Close()

	err = json.Unmarshal(responseBodyRaw, &actual)
	if err != nil {
		return types.ActualResponse{}, err
	}

	return actual, nil
}

func GetCurrency() (types.Currency, error) {
	response, err := http.Get(worker_host + "/currency")
	if err != nil {
		return types.Currency{}, err
	}

	responseBodyRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return types.Currency{}, err
	}
	defer response.Body.Close()

	var currency types.Currency
	err = json.Unmarshal(responseBodyRaw, &currency)
	if err != nil {
		return types.Currency{}, err
	}

	return currency, nil
}
