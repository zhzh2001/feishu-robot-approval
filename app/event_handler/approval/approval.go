package approval

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"xlab-feishu-robot/config"
	"xlab-feishu-robot/pkg/global"
)

type ApprovalInfo struct {
	Date    string
	Dept    string
	Manager string
	Detail  string
	Expense string
}

func ApprovalInfoByInstance(instanceCode string) *ApprovalInfo {
	resp := global.Cli.Request("get", "open-apis/approval/v4/instances/"+instanceCode, nil, nil, nil)
	if resp["status"].(string) == "APPROVED" {
		end_time := resp["end_time"].(string)
		end_time_int, _ := strconv.ParseInt(end_time, 10, 64)
		tm := time.Unix(end_time_int/1000, 0)
		date := tm.Format("2006-01-02")
		dept := resp["department_id"].(string)
		if dept != "" {
			dept = global.Cli.DepartmentGetInfoById(dept).Name
		}
		manager := ""
		for _, v := range resp["timeline"].([]interface{}) {
			if v.(map[string]interface{})["type"].(string) == "PASS" {
				uid := v.(map[string]interface{})["open_id"].(string)
				resp := global.Cli.Request("get", "open-apis/contact/v3/users/"+uid, nil, nil, nil)
				manager = resp["user"].(map[string]interface{})["name"].(string)
			}
		}
		var form []map[string]interface{}
		json.Unmarshal([]byte(resp["form"].(string)), &form)
		detail := ""
		expense := ""
		for _, v := range form {
			if v["name"].(string) == "采购事由" {
				detail = v["value"].(string) + "："
			} else if v["name"].(string) == "费用明细" {
				for _, ext := range v["ext"].([]interface{}) {
					if ext.(map[string]interface{})["type"].(string) == "amount" {
						expense = ext.(map[string]interface{})["value"].(string)
					}
				}
				for i, v1 := range v["value"].([]interface{}) {
					if i > 0 {
						detail += "、"
					}
					for _, v2 := range v1.([]interface{}) {
						if v2.(map[string]interface{})["name"].(string) == "名称" {
							detail += v2.(map[string]interface{})["value"].(string) + " "
						} else if v2.(map[string]interface{})["name"].(string) == "金额" {
							detail += strconv.FormatFloat(v2.(map[string]interface{})["value"].(float64), 'g', -1, 64)
						}
					}
				}
			}
		}
		return &ApprovalInfo{
			Date:    date,
			Dept:    dept,
			Manager: manager,
			Detail:  detail,
			Expense: expense,
		}
	}
	return nil
}

type ValueRange struct {
	Range  string          `json:"range"`
	Values [][]interface{} `json:"values"`
}

type AppendBody struct {
	VRange ValueRange `json:"valueRange"`
}

type AppendBody2 struct {
	VRanges []ValueRange `json:"valueRanges"`
}

func Receive(event map[string]any) {
	var body AppendBody
	body.VRange.Range = config.C.Token.SheetId + "!A3:A3"
	body.VRange.Values = [][]interface{}{{"-1"}}
	resp := global.Cli.Request("post", "open-apis/sheets/v2/spreadsheets/"+config.C.Token.SpreadSheetToken+"/values_append", nil, nil, body)
	ranges := resp["tableRange"].(string)
	var row int
	fmt.Sscanf(strings.Split(ranges, "!")[1], "A%d", &row)
	url := fmt.Sprintf("open-apis/sheets/v2/spreadsheets/%s/values/%s!A%d:I%d", config.C.Token.SpreadSheetToken, config.C.Token.SheetId, row-1, row-1)
	resp = global.Cli.Request("get", url, nil, nil, nil)
	values := resp["valueRange"].(map[string]interface{})["values"].([]interface{})
	id := values[0].([]interface{})[0].(float64)
	remain := values[0].([]interface{})[8].(float64)
	ainfo := ApprovalInfoByInstance(event["instance_code"].(string))
	expense, _ := strconv.ParseFloat(ainfo.Expense, 64)
	var body2 ValueRange
	body2.Range = config.C.Token.SheetId + "!A" + strconv.Itoa(row) + ":I" + strconv.Itoa(row)
	body2.Values = [][]interface{}{{id + 1, "支出", ainfo.Date, ainfo.Dept, ainfo.Manager, ainfo.Detail, nil, -expense, remain - expense}}
	global.Cli.Request("post", "open-apis/sheets/v2/spreadsheets/"+config.C.Token.SpreadSheetToken+"/values_batch_update", nil, nil, AppendBody2{VRanges: []ValueRange{body2}})
}
