package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gen2brain/beeep"
	"github.com/joho/godotenv"
)

type CustomField struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type NameId struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Issue struct {
	Id              int           `json:"id"`
	Project         NameId        `json:"project"`
	Tracker         NameId        `json:"tracker"`
	Status          NameId        `json:"status"`
	Priority        NameId        `json:"priority"`
	Author          NameId        `json:"author"`
	Assigned_to     NameId        `json:"assigned_to"`
	Subject         string        `json:"subject"`
	Description     string        `json:"description"`
	Start_date      string        `json:"start_date"`
	Due_date        string        `json:"due_date"`
	Done_ratio      int           `json:"done_ration"`
	Is_private      bool          `json:"is_private"`
	Estimated_hours string        `json:"estimated_hours"`
	Custom_fields   []CustomField `json:"custom_fields"`
	Created_on      string        `json:"created_on"`
	Updated_on      string        `json:"updated_on"`
	Closed_on       string        `json:"closed_on"`
}

type Payload struct {
	Issues      []Issue `json:"issues"`
	Total_count int     `json:"total_count"`
	Offset      int     `json:"offset"`
	Limit       int     `json:"limit"`
}

var NotifiedIssues = make(map[int]bool)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	http.HandleFunc("/track", func(w http.ResponseWriter, r *http.Request) {
		var issues Payload = fetch()
		if issues.Total_count == 0 {
			fmt.Fprintf(w, "done\n")
			return
		}
		for i := 0; i < issues.Total_count-issues.Limit; i++ {
			if !NotifiedIssues[issues.Issues[i].Id] {
				NotifiedIssues[issues.Issues[i].Id] = true

				err := beeep.Notify(issues.Issues[i].Subject, issues.Issues[i].Author.Name, "")
				if err != nil {
					panic(err)
				}
			}
		}
		fmt.Fprintf(w, "done\n")
	})
	http.ListenAndServe(":8999", nil)
}

func fetch() Payload {
	IP := os.Getenv("IP")
	url := fmt.Sprintf("http://%v/redmine/issues.json?assigned_to_id=22", IP)
	method := "GET"

	payload := strings.NewReader(fmt.Sprint(`http://%v/redmine/projects/admin-panel/issues?c%5B%5D=cf_1&c%5B%5D=cf_5&c%5B%5D=tracker&c%5B%5D=status&c%5B%5D=done_ratio&c%5B%5D=priority&c%5B%5D=subject&c%5B%5D=cf_3&c%5B%5D=assigned_to&c%5B%5D=author&c%5B%5D=updated_on&c%5B%5D=project&c%5B%5D=created_on&c%5B%5D=last_updated_by&c%5B%5D=attachments&f%5B%5D=status_id&f%5B%5D=assigned_to_id&f%5B%5D=&group_by=&op%5Bassigned_to_id%5D=%3D&op%5Bstatus_id%5D=o&set_filter=1&sort=status%2Cid%3Adesc&t%5B%5D=&utf8=%E2%9C%93&v%5Bassigned_to_id%5D%5B%5D=me`, IP))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err, " error")
		return Payload{}
	}
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Cookie", "Redmine=ZXkxYm5IcFVFK3FpUGxORS9FcXdKMlp6bGVxK2txd0dVYmZhT216WHI3bHRrbC9SSkVFMXBHVXJBVDZDQTA5Z0JYeXhHMkdkNTltMjg5SEZYVGZtcW5IT1lqVm96MjdHMEhDb0lxVGF1b3hyb0Q3dURMcWdwVllIdm5uRHE2dHlLZ2kwSkNsS0tMUGY1NVdFcmJHNVJyNm1kSFJNU3BOSlZPQmFoa29mZTB6VmxENWc3bkJvcC8wRStzZk1ZYnVwaXI5RUd2Z3oxR3ZFVExZZkpJSXFUdng4SnRLSDJmTE1JUUdMOTRVSkhDb0orRVZNMDVuamtwZ3l3NmZXdnIwcHl2Q3M3VHZDMWZ4cStGMjZIbkNKLzc1TjU4QmdyODl5TlhhaVd3MFVHVlhacG5LeWw1djVtemlTTUQ2eXpZNFFmV2hrUG9ZUzhZOGU3QUFJRkRDcGlINGY2YWllMDdrL1V3QVB4bzFRYzl3MXcxTnl2V0xCaEpqVk9IWS93cDJrLS1ac3ZmOTZOT2VML2NXUXQ0bnFCN3FRPT0%3D--ac0e51a70110a8c0523ceb6eb288bcaa2c18a47d")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err, " error")
		return Payload{}
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err, " error")
		return Payload{}
	}

	var issues Payload
	_ = json.Unmarshal(body, &issues)
	// fmt.Println(issues.Total_count)
	return issues
}
