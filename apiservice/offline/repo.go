package offline

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"strconv"
)

var GlobalDB *bolt.DB

//tables define
const (
	ProductTable                 = "PRODUCTS"
	SKUTable                     = "PRODUCT-SKU"
	DefaultATS           TableId = 10
	CustomerTable                = "CUSTOMER"
	OrderTable                   = "ORDERS"
	RecommandProductsNum         = 15
)

func RepoGetSalesOrder(orderId TableId) (result interface{}) {

	GlobalDB.Update(func(tx *bolt.Tx) error {
		orderBucket, err := tx.CreateBucketIfNotExists([]byte(OrderTable))
		if err != nil {
			HandleError(err)
			return err
		}

		orderBytes := orderBucket.Get(orderId.ToBytes())
		json.Unmarshal(orderBytes, &result)

		return nil
	})
	return
}

func RepoCreateOrder(order interface{}) map[string]interface{} {
	var newOrder map[string]interface{}

	GlobalDB.Update(func(tx *bolt.Tx) error {
		orderBucket, err := tx.CreateBucketIfNotExists([]byte(OrderTable))
		if err != nil {
			HandleError(err)
			return err
		}
		newOrder = order.(map[string]interface{})

		newId, _ := orderBucket.NextSequence()
		newOrder["id"] = TableId(newId)
		newOrder["billingAddr"] = GetAddressObj(newOrder["billingAddress"])
		newOrder["shippingAddr"] = GetAddressObj(newOrder["shippingAddress"])

		orderBytes, _ := json.Marshal(order)
		orderBucket.Put(TableId(newId).ToBytes(), orderBytes)

		return nil
	})

	return newOrder
}

func GetProductATS(ProductId TableId) int64 {
	var atsQua int64

	GlobalDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(SKUTable))

		if err != nil {
			HandleError(err)
			return err
		}

		bb, errbb := b.CreateBucketIfNotExists(ProductId.ToBytes())

		if errbb != nil {
			HandleError(errbb)
			return err
		}

		ats := bb.Get([]byte("ats"))
		if ats == nil {
			//new product id, need initialize the ats info
			bb.Put([]byte("ats"), DefaultATS.ToBytes())
			atsQua = int64(DefaultATS)
		} else {
			result := ToInt64FromBytes(ats)
			if err != nil {
				HandleError(err)
				return err
			}
			atsQua = result
		}
		return nil
	})
	return atsQua
}

func RepoCreateATSRsp(req *ATSReq) []ATSRsp {
	var rest []ATSRsp

	for _, atsR := range req.SkuIds {

		atsQua := GetProductATS(atsR)
		rsp := ATSRsp{
			SkuId:          atsR,
			Ats:            atsQua,
			AllowBackOrder: true,
		}
		rest = append(rest, rsp)
	}

	log.Printf("ATS Rsp %+v\n", rest)

	return rest
}

func RepoCreateRecommandationProducts(Id TableId) []TableId {
	var ProductId []TableId

	GlobalDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(ProductTable))
		if err != nil {
			HandleError(err)
			return err
		}

		c := b.Cursor()
		num := 0

		for k, _ := c.First(); k != nil && (num < RecommandProductsNum); k, _ = c.Next() {
			ProductId = append(ProductId, TableId(ToInt64FromBytes(k)))
			num++
		}
		_, errbb := b.CreateBucketIfNotExists(Id.ToBytes())

		if errbb != nil {
			HandleError(errbb)
			return err
		}

		return nil
	})

	return ProductId
}

func RepoCreateAccount(customer CustomerCreate) CustomerCreateRsp {
	var result CustomerCreateRsp
	key := []byte(CustomerTable)

	GlobalDB.Update(func(tx *bolt.Tx) error {
		pb, err := tx.CreateBucketIfNotExists(key)

		if err != nil {
			HandleError(err)
			return err
		}

		customerId, _ := pb.NextSequence()

		b, err := pb.CreateBucketIfNotExists([]byte(customer.Account))

		if err != nil {
			HandleError(err)
			return err
		}

		user := b.Get([]byte("User"))
		if user == nil {
			//create new user
			customer.AccountInfo.CustomerID = TableId(customerId)
			customer.AccountInfo.CustomerCode = "offline" + strconv.FormatInt(customer.AccountInfo.CustomerID.ToInt(), 10)
			customer.AccountInfo.ChannelAccountID = TableId(customer.ChannelId * customer.AccountInfo.CustomerID)
			cusStream, _ := json.Marshal(&customer)
			b.Put([]byte("User"), cusStream)
			//use Id to get account
			pb.Put(customer.AccountInfo.CustomerID.ToBytes(), []byte(customer.Account))
			result = customer.AccountInfo
		}

		json.Unmarshal(user, &customer)
		result = customer.AccountInfo
		return nil
	})

	return result
}

func RepoCreateAddress(customer *CustomerAddress) (result interface{}) {
	key := []byte(CustomerTable)

	GlobalDB.Update(func(tx *bolt.Tx) error {
		pb, err := tx.CreateBucketIfNotExists(key)

		if err != nil {
			HandleError(err)
			return err
		}

		account := pb.Get(customer.CustomerInfo.Id.ToBytes())
		if account == nil {
			result = "not found this user:" + customer.CustomerInfo.Id.ToString()
			return nil
		} else {
			customerBucket := pb.Bucket(account)
			if customerBucket == nil {
				result = "not found this account:" + string(account)
				return nil
			}
			//create the bucket for store address info
			addressBucket, _ := customerBucket.CreateBucketIfNotExists([]byte("addresses"))
			addressId, _ := addressBucket.NextSequence()
			customer.Id = TableId(addressId)

			streamD, _ := json.Marshal(customer)
			addressBucket.Put(TableId(addressId).ToBytes(), streamD)
			result = customer
		}

		return nil
	})
	return
}

func RepoUpdateAddress(addressId TableId, customer *CustomerAddress) (result interface{}) {
	key := []byte(CustomerTable)

	GlobalDB.Update(func(tx *bolt.Tx) error {
		pb, err := tx.CreateBucketIfNotExists(key)

		if err != nil {
			HandleError(err)
			return err
		}
		account := pb.Get(customer.CustomerInfo.Id.ToBytes())
		if account == nil {
			result = "not found this user:" + customer.CustomerInfo.Id.ToString()
			return nil
		} else {
			customerBucket := pb.Bucket(account)
			if customerBucket == nil {
				result = "not found this account:" + string(account)
				return nil
			}

			//create the bucket for store address info
			addressBucket, _ := customerBucket.CreateBucketIfNotExists([]byte("addresses"))
			oldAddress := addressBucket.Get(addressId.ToBytes())
			var oldCustomerAddr CustomerAddress
			json.Unmarshal(oldAddress, &oldCustomerAddr)

			//set to new one
			oldCustomerAddr.AddressInfo = customer.AddressInfo
			streamD, _ := json.Marshal(oldCustomerAddr)

			addressBucket.Put(TableId(addressId).ToBytes(), streamD)
			result = oldCustomerAddr
		}
		return nil
	})
	return
}

func RepoGetCustomerAddress(customerId TableId) (result interface{}) {
	key := []byte(CustomerTable)

	GlobalDB.Update(func(tx *bolt.Tx) error {
		pb, err := tx.CreateBucketIfNotExists(key)

		if err != nil {
			HandleError(err)
			return err
		}
		account := pb.Get(customerId.ToBytes())
		if account == nil {
			result = "not found this user:" + customerId.ToString()
			return nil
		} else {
			customerBucket := pb.Bucket(account)
			if customerBucket == nil {
				result = "not found this account:" + string(account)
				return nil
			}

			countinfo := make(map[string]interface{})
			var count int64 = 0

			//create the bucket for store address info
			addressBucket := customerBucket.Bucket([]byte("addresses"))
			if addressBucket == nil {
				countinfo["odata.count"] = 0
				result = countinfo
			} else {
				var bos []interface{}
				cur := addressBucket.Cursor()
				for k, v := cur.First(); k != nil; k, v = cur.Next() {
					count++
					var Addr CustomerAddress
					json.Unmarshal(v, &Addr)
					bos = append(bos, Addr)
				}

				countinfo["odata.count"] = count
				countinfo["value"] = bos
				result = countinfo
			}
		}
		return nil
	})
	return
}

func init() {
	GlobalDB, _ = bolt.Open("./EshopOfflineServerDB", 0666, nil)
	//	go func() {
	//		prev := GlobalDB.Stats()
	//		for {
	//			time.Sleep(10 * time.Second)

	//			states := GlobalDB.Stats()
	//			diff := states.Sub(&prev)
	//			//json.NewEncoder(os.Stderr).Encode(diff)
	//			prev = states
	//		}
	//	}()
}
