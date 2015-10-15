package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/url"
	"strconv"
	"strings"
)

var DSNWithoutSchema string //like: "root:12345@tcp(cnpvgvb1ep140.pvgl.sap.corp:3306)/"

func GetDBHandler() *sql.DB {
	db, err := sql.Open("mysql", DSNWithoutSchema)

	if err != nil {
		log.Println(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		log.Println("cannot open:", DSNWithoutSchema)
		return nil
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		log.Println(err.Error()) // proper error handling instead of panic in your app
		log.Println("cannot connect:", DSNWithoutSchema)
		return nil
	}

	return db
}

func GetSchemaById(shopId int64) (result string) {
	db := GetDBHandler()
	if db == nil {
		log.Println("Open db failed...", DSNWithoutSchema)
		return ""
	}

	defer db.Close()

	DBNames := GetAllDataBaseName(db)
	log.Println("DB names :", DBNames)
	for _, schema := range DBNames {
		log.Println("Querying ", schema)
		sel := "USE " + schema
		db.Exec(sel)

		rows, err := db.Query("select option_value from wp_options where option_name = 'eshopUrlList'")
		if err != nil {
			//log.Println(err.Error()) // proper error handling instead of panic in your app
			continue
		}

		var res []byte
		// Fetch rows
		for rows.Next() {
			// get RawBytes from data
			err = rows.Scan(&res)
			if err != nil {
				//log.Println(err.Error()) // proper error handling instead of panic in your app
				continue
			}
			var dataS []interface{}
			err = json.Unmarshal(res, &dataS)
			if err != nil {
				log.Println("Unmarshal error ", err)
				continue
			}
			m := dataS[0].(map[string]interface{})
			urlstr := m["http"].(string)
			u, err := url.Parse(urlstr)
			if err != nil {
				log.Println("Parser URL failed...", err)
				continue
			}
			q := u.Query()

			if _, ok := q["eshop_id"]; !ok {
				continue
			}

			log.Println("Eshop id is ", q["eshop_id"][0])
			aId, _ := strconv.ParseInt(q["eshop_id"][0], 10, 8)
			if aId == shopId {
				result = schema
				return
			}
			break
		}
	}
	return
}

func GetSchemaByName(shopName string) (result []string) {

	db := GetDBHandler()
	if db == nil {
		log.Println("Open db failed...", DSNWithoutSchema)
		return nil
	}

	defer db.Close()

	DBNames := GetAllDataBaseName(db)
	log.Println("DB names :", DBNames)

	for _, schema := range DBNames {
		sel := "USE " + schema
		db.Exec(sel)

		rows, err := db.Query("select option_value from wp_options where option_name = 'eshopSetting'")
		if err != nil {
			//log.Println(err.Error()) // proper error handling instead of panic in your app
			continue
		}

		var res []byte
		var name string
		// Fetch rows
		for rows.Next() {
			// get RawBytes from data
			err = rows.Scan(&res)
			if err != nil {
				//log.Println(err.Error()) // proper error handling instead of panic in your app
				continue
			}
			var dataS interface{}
			json.Unmarshal(res, &dataS)
			m := dataS.(map[string]interface{})
			name = m["shopName"].(string)
			break
		}
		log.Println("Name is ", name)
		if strings.Contains(name, shopName) {
			result = append(result, schema)
		}
	}
	return
}

func GetAllDataBaseName(db *sql.DB) (result []string) {
	rows, err := db.Query("show databases")
	if err != nil {
		log.Println(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, 1)

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Println(err.Error()) // proper error handling instead of panic in your app
			continue
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for _, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				continue
			} else {
				value = string(col)
			}
			result = append(result, value)
		}
	}
	if err != nil {
		log.Println(err.Error()) // proper error handling instead of panic in your app
	}

	return
}
