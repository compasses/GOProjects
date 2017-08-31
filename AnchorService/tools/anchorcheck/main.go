package main

import (
	"fmt"
	"os"

	"AnchorService/anchor"
	"encoding/json"
	"github.com/FactomProject/factom"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "checkanchor",
			Aliases: []string{"c"},
			Usage:   "specify factom server to check",
			Action: func(c *cli.Context) error {
				fmt.Println("factom server : ", c.Args().First())
				CheckFactomAnchor(c.Args().First())
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func CheckFactomAnchor(factomServer string) {
	// get heights
	heightReq := factom.NewJSON2Request("heights", 0, nil)
	heightResp, err := anchor.DoFactomReq(heightReq, factomServer)
	if err != nil {
		fmt.Println("error on do factom req ", err)
		return
	}
	var result factom.HeightsResponse
	err = json.Unmarshal(heightResp.Result, &result)
	if err != nil {
		fmt.Println(" error on unmarshal ", err)
		return
	}

	height := result.DirectoryBlockHeight
	fmt.Println("start check anchor total height", height)

	for i := int64(1); i < height; i++ {
		params := struct {
			Height int64 `json:"height"`
		}{
			Height: i,
		}

		req := factom.NewJSON2Request("directory-blockinfo-by-height", 0, params)
		resp, err := anchor.DoFactomReq(req, factomServer)

		if resp.Error != nil {
			fmt.Println("directory-blockinfo-by-height error happen height", i, resp.Error.Message)
			continue
		}

		var result map[string]interface{}

		err = json.Unmarshal(resp.Result, &result)
		if err != nil {
			fmt.Println("Unmarshal error on height ", i, err)
			continue
		}

		confirm := result["BTCConfirmed"].(bool)
		btchash := result["BTCTxHash"].(string)

		fmt.Println("height ", i, " confirmd ", confirm, "BTCTxHash ", btchash)
		//fmt.Println("got blockinfo ", spew.Sdump(result))

	}

}
