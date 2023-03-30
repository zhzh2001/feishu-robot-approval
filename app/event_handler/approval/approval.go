package approval

import (
	"fmt"
	"strconv"
	"strings"
	"xlab-feishu-robot/config"
	"xlab-feishu-robot/pkg/global"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

type ApprovalInfo struct {
	Date    string
	Dept    string
	Manager string
	Detail  string
	Expense string
}

func ApprovalInfoByInstance(InstanceCode string) *ApprovalInfo {
	info := global.Cli.ApprovalInstanceById(InstanceCode)
	if info.Status == "APPROVED" {
		date := info.EndTime.Format("2006-01-02")
		dept := info.DepartmentId
		if dept != "" {
			dept = global.Cli.DepartmentGetInfoById(dept).Name
		}
		manager := ""
		for _, v := range info.Timeline {
			if v.Type == "PASS" {
				manager = global.Cli.UserInfoById(v.OpenId, feishuapi.OpenId).Name
			}
		}
		form := info.Form
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

func Receive(event map[string]any) {
	ranges := global.Cli.SheetAppendData(config.C.Token.SpreadSheetToken, config.C.Token.SheetId, "A3:A3", [][]interface{}{{-1}})
	var row int
	fmt.Sscanf(strings.Split(ranges, "!")[1], "A%d", &row)
	ranges = "A" + strconv.Itoa(row-1) + ":I" + strconv.Itoa(row-1)
	values := global.Cli.SheetGetData(config.C.Token.SpreadSheetToken, config.C.Token.SheetId, ranges)
	id := values[0].([]interface{})[0].(float64)
	if id < 0 {
		logrus.Warn("Spreadsheet layout is not correct, please check it.")
		return
	}
	remain := values[0].([]interface{})[8].(float64)
	ainfo := ApprovalInfoByInstance(event["instance_code"].(string))
	expense, _ := strconv.ParseFloat(ainfo.Expense, 64)
	ranges = "A" + strconv.Itoa(row) + ":I" + strconv.Itoa(row)
	global.Cli.SheetWriteData(config.C.Token.SpreadSheetToken, config.C.Token.SheetId, ranges, [][]interface{}{{id + 1, "支出", ainfo.Date, ainfo.Dept, ainfo.Manager, ainfo.Detail, nil, -expense, remain - expense}})
}
