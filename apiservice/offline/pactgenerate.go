package offline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Compasses/GOProjects/apiservice/utils"
	"github.com/SEEK-Jobs/pact-go"
	"github.com/SEEK-Jobs/pact-go/provider"
)

const PactsDir = "./pacts"

type ProviderAPIClient struct {
	baseURL string
}

func (c *ProviderAPIClient) ClientRun(method, path string, reqBody []byte) error {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	newbody := make([]byte, len(reqBody))

	req, err := http.NewRequest(method, url, ioutil.NopCloser(bytes.NewReader(newbody)))
	if err != nil {
		log.Println(err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	var res interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&res); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (middleware *offlinemiddleware) GetPactFile() string {
	files, err := ioutil.ReadDir(PactsDir)
	if err != nil {
		log.Println(err)
		return ""
	}
	//now just run the first file
	for _, file := range files {
		log.Println("upload pact file name ", file.Name())
		return PactsDir + "/" + file.Name()
	}
	return ""
}

func (middleware *offlinemiddleware) buildPact(consumerName, providerName string) pact.Builder {
	return pact.
		NewConsumerPactBuilder(&pact.BuilderConfig{PactPath: PactsDir}).
		ServiceConsumer(consumerName).
		HasPactWith(providerName)
}

func (middleware *offlinemiddleware) RunPact(ms pact.ProviderService, path, method string, reqBody, respBody interface{}, statusCode int,
	consumerName, providerName, msUrl string) {
	req := utils.JsonInterfaceToByte(reqBody)

	request := provider.NewJSONRequest(method, path, "", nil)
	request.SetBody(reqBody)

	header := make(http.Header)
	header.Add("content-type", "application/json")
	response := provider.NewJSONResponse(statusCode, header)
	response.SetBody(respBody)

	//Register interaction for this test scope
	if err := ms.Given(consumerName).
		UponReceiving(providerName).
		With(*request).
		WillRespondWith(*response); err != nil {
		fmt.Println(err)
		//t.FailNow()
	}

	//test
	client := &ProviderAPIClient{baseURL: msUrl}
	if err := client.ClientRun(method, path, req); err != nil {
		log.Println(err)
		//t.FailNow()
	}

	//Verify registered interaction
	if err := ms.VerifyInteractions(); err != nil {
		log.Println(err)
		//t.FailNow()
	}

	//Clear interaction for this test scope, if you need to register and verify another interaction for another test scope
	ms.ClearInteractions()
}

func (middleware *offlinemiddleware) GenPactWithProvider() {
	builder := middleware.buildPact("EShop Online Store", "EShop Adaptor")
	ms, msUrl := builder.GetMockProviderService()
	//map[string]map[string][]interface{}
	//"Path", "Method", "[req..., rsp...,]"
	interactMap := middleware.replaydb.GetJSONMap()
	for path, value := range interactMap {
		for method, interacts := range value {
			var count int = 0
			for _, detailMapel := range interacts {
				detailMapItem := detailMapel.(map[string]interface{})
				request, ok := detailMapItem["request"]
				if !ok {
					log.Println("missing request, continue ", detailMapItem)
					continue
				}
				respose, ok := detailMapItem["response"]
				if !ok {
					log.Println("missing response, continue ", detailMapItem)
					continue
				}
				responseMap := respose.(map[string]interface{})
				count++
				for k, v := range responseMap {
					status, _ := strconv.Atoi(k)
					//fmt.Println("\r\nstore:", request, "response", v)
					consumName := "mock server for " + path + " method " + method + " " + strconv.Itoa(count)
					provideName := "pact contract for " + path + " method " + method + " " + strconv.Itoa(count)
					middleware.RunPact(ms, path, method, request, v, status, consumName,
						provideName, msUrl)
					count++
					break
				}
			}
		}
	}

	//Finally, build to produce the pact json file
	if err := builder.Build(); err != nil {
		log.Println(err)
	}
}
