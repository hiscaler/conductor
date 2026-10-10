# 可复用参考图库

保存日期：2026-10-10。入口规则：[参考资产规则](../../platforms/reference-assets.md)。机器索引：[index.json](./index.json)。本目录是内部参考库，不是待上架图片目录。

|资产|当前使用图片|来源与状态|
|---|---|---|
|Apple iPad Air (M4) 11 英寸|[用户指定正面图](objects/ipad-air-11-m4/user-selected.png)|[Consumer Reports 原图](https://crdms.images.consumerreports.org/f_auto,w_600/prod/products/cr/models/420860-9-inch-screen-and-larger-tablets-apple-ipad-air-m4-11-10053423.png)；[Apple 规格](https://support.apple.com/en-us/126471)：高 24.76 cm、宽 17.85 cm、厚 0.61 cm；[官方尺寸图纸](objects/ipad-air-11-m4/dimensional-drawing.pdf)|
|可口可乐 Original Taste 美国版 12 fl oz 标准罐|[正面图](objects/coca-cola-original-us-12oz/front.jpg)|[Walmart 商品页](https://business.walmart.com/ip/Coca-Cola-Soda-Pop-12-fl-oz-Can/10535208)；已核对标准罐外观，品牌罐体精确尺寸待核实|
|Pepsi 美国版 12 fl oz 罐|[正面图](objects/pepsi-us-12oz/front.png)|[PepsiCo 产品页](https://www.pepsicoproductfacts.com/Home/product?PPF=&form=RTD&formula=35005*26*01-01&size=12)；品牌罐体精确尺寸待核实|
|仅文字定制马克杯操作图|[加高留白母版](customization/text-only-mug/master.png)|[用户原始附件](customization/text-only-mug/user-original.png)；图案缩小、上移，杯体保持比例，输入区不印占位字|

## 母版处理记录

工具：image_gen，参考附件编辑。提示要求：保留杯体比例、黑色内胆与把手、步骤图标与说明、红色箭头；缩小并上移小狗和南瓜整组图案，将下方文字留白高度增加约一半；杯面不添加文字。经三次位置调整后保存当前母版。已打开最终保存文件目检。

留白目测记录（1254 px 画布）：统一以杯面下方安全边界 y≈1117 计，原图图案最低点 y≈1004，调整后 y≈925；留白约 113→192 px。原要求增加约一半；本次根据杯口间距反馈再次缩小图案，保留并稍扩大下方留白。该数据是版式测量，不是商品可印尺寸。

## 使用状态

- iPad 用户指定图作为优先外观参考；官方规格用于尺寸依据。
- 两款可乐可直接加载作外观参考；精确比例图须先取得对应罐型尺寸依据。不能把通用 355 ml 罐规格默认为品牌事实。
- 本次保存来源图片，未确认公开商业再发布授权；内部参考与公开发布需区分。
- 后续按当前 SKU 替换图案及真实杯色，再逐张目检。当前调整稿由用户要求生成，尚未记录用户对调整稿的视觉验收。

2026-10-10 杯口间距修正：使用内置 image_gen 编辑当前母版，完整图案等比缩小，最高点约 y=423，与杯口下缘约留 50 px；保留杯体、步骤说明和箭头。已目检最终保存图片。
