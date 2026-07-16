# 目标市场数据源规则

趋势找品不能默认所有国家都使用同一套数据源。指挥家必须先识别目标国家/地区，再选择对应的搜索趋势、内容平台、电商验证和本地语言关键词。

本文件列出常用候选数据源。实际执行时必须以当前可访问性和目标市场为准；无法访问、需要登录、没有公开数据或数据源不适用时，必须标注“未获取”或“受访问限制”，不得编造数据。

## 执行顺序

```text
目标国家/地区
-> 本地语言和英文关键词
-> 搜索/趋势数据源
-> 短视频/内容平台
-> 电商平台验证
-> 供应链平台验证
-> 候选商品评分
```

如果用户未提供目标市场，默认先按美国市场处理，但必须在输出中标注“目标市场未提供，暂按美国市场假设”。

## 关键词规则

- 美国、英国、加拿大、澳大利亚：优先英文关键词。
- 德国：必须生成德语关键词，并保留英文对照。
- 法国：必须生成法语关键词，并保留英文对照。
- 日本：必须生成日语关键词，并保留英文对照。
- 韩国：必须生成韩语关键词，并保留英文对照。
- 西班牙、墨西哥：必须生成西班牙语关键词，并保留英文对照。
- 巴西：必须生成葡萄牙语关键词，并保留英文对照。
- 中国：必须生成中文关键词，并根据平台补充口语化表达。
- 东南亚：根据目标国家生成本地语言关键词；如果无法确认，至少使用英文 + 本地常用表达假设，并标注待确认。

不得只用中文或英文关键词验证非中文/非英语市场。

## 常用数据源映射

| 目标市场 | 搜索/趋势候选 | 短视频/内容候选 | 电商验证候选 | 关键词语言 |
| --- | --- | --- | --- | --- |
| 美国 | Google Trends、Google Search、搜索建议 | TikTok、YouTube、Instagram、Pinterest、Reddit | Amazon、Walmart、Etsy、TikTok Shop、eBay | 英文 |
| 英国 | Google Trends、Google Search | TikTok、YouTube、Instagram、Pinterest | Amazon UK、eBay UK、Etsy、TikTok Shop | 英文 |
| 加拿大 | Google Trends、Google Search | TikTok、YouTube、Instagram、Pinterest、Reddit | Amazon CA、Walmart Canada、Etsy、eBay CA | 英文，可补法语 |
| 澳大利亚 | Google Trends、Google Search | TikTok、YouTube、Instagram、Pinterest | Amazon AU、eBay AU、Kogan、Etsy | 英文 |
| 德国 | Google Trends、Google Search | YouTube、Instagram、TikTok、Pinterest | Amazon DE、Otto、eBay DE、Kaufland | 德语 + 英文 |
| 法国 | Google Trends、Google Search | YouTube、Instagram、TikTok、Pinterest | Amazon FR、Cdiscount、Etsy、ManoMano | 法语 + 英文 |
| 日本 | Google Trends、Google Search、Yahoo Japan | YouTube、Instagram、TikTok、X | Amazon JP、Rakuten、Yahoo Shopping、Mercari | 日语 + 英文 |
| 韩国 | Naver DataLab、Naver Search、Google Trends | YouTube、Instagram、TikTok、Naver Blog | Coupang、Naver Shopping、Gmarket、11st | 韩语 + 英文 |
| 中国 | 百度指数、百度搜索、巨量算数、微信指数 | 抖音、小红书、B站、快手、微博 | 淘宝、天猫、京东、拼多多、抖音电商 | 中文 |
| 东南亚 | Google Trends、Google Search、本地搜索建议 | TikTok、YouTube、Facebook、Instagram | Shopee、Lazada、TikTok Shop、Tokopedia、Bukalapak | 英文 + 本地语言 |
| 墨西哥 | Google Trends、Google Search | TikTok、YouTube、Instagram、Facebook | Mercado Libre、Amazon MX、Walmart MX | 西班牙语 + 英文 |
| 巴西 | Google Trends、Google Search | TikTok、YouTube、Instagram、Facebook | Mercado Livre、Amazon BR、Shopee BR | 葡萄牙语 + 英文 |
| 中东 | Google Trends、Google Search | TikTok、YouTube、Instagram、Snapchat | Amazon UAE/SA、Noon、Namshi | 阿拉伯语 + 英文 |

## 验证原则

- 趋势平台只回答“方向是否变热”，不能直接证明商品能卖。
- 内容平台用于验证真实场景、痛点、商品出现频率和购买意图。
- 电商平台用于验证竞品、价格带、图片表达、评价痛点和品牌集中度。
- 供应链平台用于验证 MOQ、成本、交期、规格稳定性和小批量可行性。

## 输出要求

在 `production-output.md` 的 `## 3. 选品初筛` 中必须写明：

- 目标市场数据源：
- 本地语言关键词：
- 趋势数据获取状态：
- 短视频/内容数据获取状态：
- 电商验证数据获取状态：

如果某个市场的数据源未覆盖，输出“需补充目标市场数据源清单”，不要套用美国市场结论。
