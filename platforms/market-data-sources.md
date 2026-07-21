# 目标市场数据源规则

趋势找品不能默认所有国家都使用同一套数据源。指挥家必须先识别目标国家/地区，再选择对应的搜索趋势、内容平台、电商验证和本地语言关键词。

本文件列出常用候选数据源。实际执行时必须以当前可访问性和目标市场为准；无法访问、需要登录、没有公开数据或数据源不适用时，必须标注“未获取”或“受访问限制”，不得编造数据。

## 执行顺序

```text
目标国家/地区
-> 本地语言和英文关键词
-> 搜索/趋势数据源
-> 短视频/内容平台
-> 电商平台热销榜与搜索结果
-> 代表性商品详情页与评论
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
- 电商验证应优先组合四类证据：平台热销榜、同类商品搜索结果、代表性商品详情页、可访问的用户评论。
- 热销榜用于判断当前热卖款式、规格组合、价格带和主流视觉表达；搜索结果用于观察竞争密度；详情页用于提取卖点和属性；评论用于提取购买动机、好评理由、差评痛点和未满足期待。
- 热销榜和排名属于动态数据，必须记录目标市场、榜单名称或入口、抓取日期和获取状态，不得把当前排名写成长期事实。
- 每次研究优先覆盖头部热销商品、新品或增长型商品、高销量但差评明显的商品，避免只复制榜首商品。
- 不同平台的数据只解释对应平台和市场，不得用 Amazon 美国榜单代替 Temu 或其他国家市场结论。
- 评论、销量、评分和排名无法访问时必须写“未获取”或“受访问限制”，不得用常识补写。
- 供应链平台用于验证 MOQ、成本、交期、规格稳定性和小批量可行性。
- 供应链验证优先使用 1688、Alibaba、义乌购、Made-in-China 或用户指定的供应商平台；中国供应链优先检查 1688 和义乌购，国际供应链可补充 Alibaba。
- 供应链候选供应商默认查询 3-5 家，记录店铺名称、商品/店铺网址、主营产品、开店时间或店铺年限、销量/成交公开口径、评分、MOQ、价格区间、交期、是否支持定制、小批量能力、联系电话、邮箱、微信、联系人、抓取日期和访问限制。
- 供应商联系方式只采集公开展示的企业或店铺联系方式，不绕过登录、验证码、权限或平台限制。未公开或无法访问的字段写“未公开”“需登录”或“未获取”，不得推断或编造。

## 未来节日筛选

- 节日必须根据目标国家/地区使用当地日历和当地名称。
- 只考虑当前日期之后且仍有生产、定制、运输和营销准备时间的节日；已经过去或准备周期不足的节日直接排除。
- 节日按“强相关、可测试、不建议”分级，商品与节日关联牵强时不得硬套。
- 定制示例、购买动机和广告故事可以参考未来节日，但不能因此编造包装、交期、折扣或产品功能。
- 必须记录节日日期、建议最晚准备时间、相关性和排除理由；日期无法核实时写“待确认”。

## 输出要求

在 `production-output.md` 的 `## 3. 选品初筛` 中必须写明：

- 目标市场数据源：
- 本地语言关键词：
- 趋势数据获取状态：
- 短视频/内容数据获取状态：
- 电商验证数据获取状态：
- 热销榜名称、目标市场和抓取日期：
- 商品详情页与评论样本数：
- 未来节日筛选及准备周期：

如果某个市场的数据源未覆盖，输出“需补充目标市场数据源清单”，不要套用美国市场结论。
