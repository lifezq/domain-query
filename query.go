package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type ModuleDef struct {
	Avail int    `json:"avail"`
	Name  string `json:"name"`
	Tld   string `json:"tld"`
}

type AliResponse struct {
	ErrorCode int         `json:"errorCode"`
	Module    []ModuleDef `json:"module"`
	Success   string      `json:"success"`
}

func QueryDomains(domains []string) []string {

	var set []string

	resp, err := http.Get("https://checkapi.aliyun.com/check/checkdomain?domain=" + strings.Join(domains, ",") + "&command=&token=Ya7ebb0c1de672df836a9f8c3b06697f5&ua=&currency=&site=&bid=&_csrf_token=")
	if err != nil {
		return set
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return set
	}

	var res AliResponse
	if err = json.Unmarshal(body, &res); err != nil {
		return set
	}

	if len(res.Module) > 0 {
		for _, m := range res.Module {

			if m.Avail == 1 {
				fmt.Printf("******************************* res------%v\n", m)
				set = append(set, m.Name)

			}
		}
	}

	if len(set) > 0 {
		EmailNoticeResult(set)
	}

	return set
}

func EmailNoticeResult(res []string) {
	cmd := exec.Command("sh", "-c", "echo \"Avali domain list : "+strings.Join(res, ", ")+"\"|mail -s \"OK Query Domain notice...Important\" 15001046768@139.com")
	cmd.Run()

	fmt.Printf("OK Query Domain notice : %s\n", strings.Join(res, ", "))
}

func QueryDomainValid(chars []byte, fixed int8) {

	fmt.Printf("%s------%s----fixed:%d\n", time.Now().String(), string(chars), fixed)

	if fixed == 0 {
		QueryDomains([]string{string(chars) + ".com"})
		return
	}

	var domains []string
	switch len(chars) {
	case 3:

		switch fixed {
		case 1:

			for j := chars[0]; j < 122; j++ {

				domains = []string{}
				for k := chars[1]; k < 122; k++ {
					for l := chars[2]; l < 122; l++ {
						domains = append(domains, fmt.Sprintf("%c%c%c", j, k, l)+".com")
					}
				}

				QueryDomains(domains)
			}

		case 2:
			domains = []string{}
			for h := 97; h < 122; h++ {
				domains = append(domains, fmt.Sprintf("%c%c%c%c", h, chars[0], chars[1], chars[2])+".com")
			}
			QueryDomains(domains)
		}

	case 4:

		switch fixed {
		case 1:

			for i := chars[0]; i < 122; i++ {
				for j := chars[1]; j < 122; j++ {

					domains = []string{}
					for k := chars[2]; k < 122; k++ {
						for l := chars[3]; l < 122; l++ {
							domains = append(domains, fmt.Sprintf("%c%c%c%c", i, j, k, l)+".com")
						}
					}

					QueryDomains(domains)
				}
			}

		case 2:
			domains = []string{}
			for h := 97; h < 122; h++ {
				domains = append(domains, fmt.Sprintf("%c%c%c%c%c", h, chars[0], chars[1], chars[2], chars[3])+".com")
			}
			QueryDomains(domains)
		}
	}
}

func main() {

	for {

		QueryDomainValid([]byte{0x30, 0x70, 0x61, 0x69}, 0x00)
		QueryDomainValid([]byte{0x31, 0x70, 0x61, 0x69}, 0x00)

		QueryDomainValid([]byte{0x61, 0x61, 0x61}, 0x01)
		QueryDomainValid([]byte{0x61, 0x61, 0x61, 0x61}, 0x01)

		QueryDomainValid([]byte{0x70, 0x61, 0x69}, 0x02)
		QueryDomainValid([]byte{0x63, 0x6f, 0x69, 0x6e}, 0x02)

		time.Sleep(time.Hour * 24)
	}
}
