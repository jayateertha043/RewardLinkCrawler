package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/jayateertha043/RewardLinkCrawler/pkg/httpclient"
)

var urls []string
var count int

func printBanner() {
	fmt.Println(`
	__________                                .___.____    .__        __   _________                      .__                
	\______   \ ______  _  _______ _______  __| _/|    |   |__| ____ |  | _\_   ___ \____________ __  _  _|  |   ___________ 
	 |       _// __ \ \/ \/ /\__  \\_  __ \/ __ | |    |   |  |/    \|  |/ /    \  \/\_  __ \__  \\ \/ \/ /  | _/ __ \_  __ \
	 |    |   \  ___/\     /  / __ \|  | \/ /_/ | |    |___|  |   |  \    <\     \____|  | \// __ \\     /|  |_\  ___/|  | \/
	 |____|_  /\___  >\/\_/  (____  /__|  \____ | |_______ \__|___|  /__|_ \\______  /|__|  (____  /\/\_/ |____/\___  >__|   
			\/     \/             \/           \/         \/       \/     \/       \/            \/                 \/       
																															 
																															 
																																																														 
	`)
}

func buildHeaders() map[string]string {
	headers := make(map[string]string)
	headers["accept"] = "application/json, text/plain, */*"
	headers["x-platform-hash"] = "undefined"
	headers["x-redemption-session-id"] = "undefined"
	headers["accept-encoding"] = "gzip, deflate, br"
	return headers
}
func buildParams(token string) map[string]string {
	params := make(map[string]string)
	params["encryptedRequestCode"] = token
	params["urlVersion"] = "1"
	return params
}
func stringify(a interface{}) string {
	if a != nil {
		return fmt.Sprint(a)
	}
	return ""
}

func printData(data map[string]interface{}, csvwriter *csv.Writer) {
	if data != nil {
		var row []string
		count = count + 1
		if data["status"] == "ACTIVE" {
			//fmt.Println("availableAmount", "currencyCode", "checkoutEmail", "createdAt", "expirationDate", "redemptionType", "checkoutFirstName", "checkoutLastName", "recipientCountryCode", "referenceLineItemId", "currentAmount", "logoUrl", "expirationMessage")
			row = append(row,
				stringify(data["url"]),
				stringify(data["availableAmount"]),
				stringify(data["currencyCode"]),
				stringify(data["checkoutEmail"]),
				stringify(data["createdAt"]),
				stringify(data["expirationDate"]),
				stringify(data["redemptionType"]),
				stringify(data["checkoutFirstName"]),
				stringify(data["checkoutLastName"]),
				stringify(data["recipientCountryCode"]),
				stringify(data["referenceLineItemId"]),
				stringify((data["currentAmount"])),
				stringify(data["customization"].(interface{}).(map[string]interface{})["logoUrl"]),
				stringify(data["customization"].(interface{}).(map[string]interface{})["expirationMessage"]))
			csvwriter.Write(row)

			defer csvwriter.Flush()
			fmt.Println()
			fmt.Println("Index:", count)
			fmt.Println("URL:", stringify(data["url"]))
			fmt.Println("availableAmount:", stringify(data["availableAmount"]))
			fmt.Println("currencyCode:", stringify(data["currencyCode"]))
			fmt.Println("checkoutEmail:", (data["checkoutEmail"]))
			fmt.Println("createdAt:", stringify(data["createdAt"]))
			fmt.Println("expirationDate:", stringify(data["expirationDate"]))
			fmt.Println("redemptionType:", stringify(data["redemptionType"]))
			fmt.Println("checkoutFirstName:", stringify(data["checkoutFirstName"]))
			fmt.Println("checkoutLastName:", stringify(data["checkoutLastName"]))
			fmt.Println("recipientCountryCode:", stringify(data["recipientCountryCode"]))
			fmt.Println("referenceLineItemId:", stringify(data["referenceLineItemId"]))
			fmt.Println("currentAmount:", stringify((data["currentAmount"])))
			fmt.Println("logoUrl:", stringify(data["customization"].(interface{}).(map[string]interface{})["logoUrl"]))
			fmt.Println("expirationMessage:", stringify(data["customization"].(interface{}).(map[string]interface{})["expirationMessage"]))

		}

	}
}
func main() {

	printBanner()
	fmt.Println()
	csvFile, err := os.Create("RewardLinkCrawler.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	// fetch for all urls from stdin
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		urls = append(urls, sc.Text())
	}

	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
	}
	fields := [14]string{"url", "availableAmount", "currencyCode", "checkoutEmail", "createdAt", "expirationDate", "redemptionType", "checkoutFirstName", "checkoutLastName", "recipientCountryCode", "referenceLineItemId", "currentAmount", "logoUrl", "expirationMessage"}
	csvwriter.Write(fields[:])
	//fmt.Println("availableAmount", "currencyCode", "checkoutEmail", "createdAt", "expirationDate", "redemptionType", "checkoutFirstName", "checkoutLastName", "recipientCountryCode", "referenceLineItemId", "currentAmount", "logoUrl", "expirationMessage")
	for _, u := range urls {
		match, _ := regexp.MatchString(`^https://www\.rewardlink\.io/r/1/[A-Za-z0-9_-]{43}$`, u)

		if match {
			token := u[30:]
			headers := buildHeaders()
			params := buildParams(token)
			resp, _ := httpclient.PostRequest("https://prod-backend.rewardlink.io/redemption/load", params, true, headers, 10)

			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				data := make(map[string]interface{})
				_ = json.Unmarshal(body, &data)
				data["url"] = u
				printData(data, csvwriter)
			}

		}
	}
	defer csvwriter.Flush()
	defer csvFile.Close()
}
