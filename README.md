# 财务统计机器人

| Web Framework | Log Manager     | Config Manager | Api Documentation  | Feishu Api Client     |
|:-------------:|:---------------:|:--------------:|:------------------:|:---------------------:|
| gin-gonic/gin | sirupsen/logrus | spf13/viper    | swaggo/gin-swagger | YasyaKarasu/feishuapi |

## Usage

- `app/event_handler` 自定义飞书事件处理方法
  - `approval` 审批通过事件处理
- `app/controller` 自定义 service controller
- `app/router` 为自定义的 service controller 注册路由
- `config` 在 `Config` 类型定义中添加自定义的配置字段
- `config.yaml` 添加自定义的配置字段

## Architecture

- `app` 机器人主体部分
- `config` 机器人配置
- `docs` swagger 生成的 Api 文档
- `pkg/dispatcher` 飞书事件调度器

## 自定义配置字段

配置示例：

```yaml
token:
  spreadSheetToken: shtxxxxxxxxxxxxxxxxxxxxxxxx
  sheetId: xxxxxx
  approvalCode: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

## 表格字段示例

| 序号 | 类型 | 日期 | 部门 | 经办人 | 具体事项 | 收入 | 支出 | 结余 |
|----|----|----|----|----|----|----|----|----|
| 0 |  | 2021-01-01 |  | Alice |  | 3000000 |  | 3000000 |
| 1 | 支出 | 2023-03-19 |  | Bob | 测试1：RTX 4090 200000、H100 1500000 |  | -1700000 | 1300000 |

## 部署指南

### 审批相关

在开发者后台事件订阅中，为机器人添加「审批通过」事件：

根据[审批接入指南](https://open.feishu.cn/document/ukTMukTMukTM/ukDNyUjL5QjM14SO0ITN)，获取采购审批对应的 Approval Code。

TLDR：访问<https://www.feishu.cn/approval/admin/approvalList?devMode=on>，编辑对应的审批，在 URL 中的 definitionCode 后面找到。

接下来为机器人订阅对应的审批通过事件。打开[订阅审批事件](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/approval-v4/approval/subscribe)，选择「尝试一下」，在 API 调试台中获取 token 并输入 Approval Code，调试即可完成订阅。

如果已经订阅，API 会返回 `subscription existed`

因为订阅这部分操作只需要进行一次，而且肯定需要人工完成，因此没有加入应用。下同。

### 云文档相关

根据[常见问题 - 服务端文档 - 开发文档 - 飞书开放平台](https://open.feishu.cn/document/ukTMukTMukTM/uczNzUjL3czM14yN3MTN)，将应用添加为表格协作者。

首先找到知识空间 id。打开[获取知识空间列表](https://open.feishu.cn/document/ukTMukTMukTM/uUDN04SN0QjL1QDN/wiki-v2/space/list)，记录对应的 `space_id` 字段。**务必使用 `user_access_token` 进行调试！**

接下来打开[获取知识空间子节点列表](https://open.feishu.cn/document/ukTMukTMukTM/uUDN04SN0QjL1QDN/wiki-v2/space-node/list)，记录对应的 `obj_token` 字段，并填入 `config.yaml` 的 `spreadSheetToken` 中。仍需要使用 `user_access_token`。

打开[获取工作表](https://open.feishu.cn/document/ukTMukTMukTM/uUDN04SN0QjL1QDN/sheets-v3/spreadsheet-sheet/query)，记录对应的 `sheet_id` 字段并填入 `config.yaml`。这一步可以使用 `tenant_access_token`，如果成功说明应用获得了协作者权限。
