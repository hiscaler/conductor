import assert from "node:assert/strict";
import {
  CUSTOMIZATION_CAPABILITY,
  inferSellerCustomizationService,
  NON_SELLER_POSITIONING,
  SELLER_CUSTOMIZATION_SERVICE,
  resolveListingCustomizationFlow,
} from "./listing-customization-flow.mjs";

assert.equal(
  inferSellerCustomizationService("卖家定制服务：不提供"),
  SELLER_CUSTOMIZATION_SERVICE.NOT_OFFERED,
);
assert.equal(
  inferSellerCustomizationService("卖家定制服务：必须定制"),
  SELLER_CUSTOMIZATION_SERVICE.REQUIRED,
);
assert.equal(
  inferSellerCustomizationService("卖家定制服务：定制可选"),
  SELLER_CUSTOMIZATION_SERVICE.OPTIONAL,
);
assert.equal(
  inferSellerCustomizationService("本 Listing 不提供定制，买家可自行加工"),
  SELLER_CUSTOMIZATION_SERVICE.NOT_OFFERED,
);
assert.equal(
  inferSellerCustomizationService("当前 Listing 定制可选"),
  SELLER_CUSTOMIZATION_SERVICE.OPTIONAL,
);
assert.equal(
  inferSellerCustomizationService("商品具备定制能力：是；支持文字和图片热转印"),
  SELLER_CUSTOMIZATION_SERVICE.PENDING,
);
assert.equal(
  inferSellerCustomizationService("本 Listing 必须定制；本 Listing 定制可选"),
  SELLER_CUSTOMIZATION_SERVICE.PENDING,
);

const explicitServiceFromBasicInput = resolveListingCustomizationFlow({
  capability: CUSTOMIZATION_CAPABILITY.YES,
  sellerService: inferSellerCustomizationService("卖家定制服务：可选"),
});
assert.equal(explicitServiceFromBasicInput.clarificationRequired, false);
assert.equal(explicitServiceFromBasicInput.allowProductExampleGeneration, true);

const ordinaryFivePack = resolveListingCustomizationFlow({
  capability: CUSTOMIZATION_CAPABILITY.YES,
  sellerService: SELLER_CUSTOMIZATION_SERVICE.NOT_OFFERED,
  nonSellerPositioning: NON_SELLER_POSITIONING.ORDINARY,
});
assert.equal(ordinaryFivePack.route, "ordinary_finished");
assert.equal(ordinaryFivePack.generateCustomizationMaster, false);
assert.equal(ordinaryFivePack.generateCustomizationGuide, false);
assert.equal(ordinaryFivePack.allowSellerUploadLabels, false);

const diyBlankFivePack = resolveListingCustomizationFlow({
  capability: CUSTOMIZATION_CAPABILITY.YES,
  sellerService: SELLER_CUSTOMIZATION_SERVICE.NOT_OFFERED,
  nonSellerPositioning: NON_SELLER_POSITIONING.DIY_BLANK,
});
assert.equal(diyBlankFivePack.route, "diy_blank");
assert.equal(diyBlankFivePack.defaultDeliveryState, "blank_base");
assert.equal(diyBlankFivePack.generateCustomizationGuide, false);
assert.equal(diyBlankFivePack.allowSellerUploadLabels, false);

const requiredPersonalization = resolveListingCustomizationFlow({
  capability: CUSTOMIZATION_CAPABILITY.YES,
  sellerService: SELLER_CUSTOMIZATION_SERVICE.REQUIRED,
  nonSellerPositioning: NON_SELLER_POSITIONING.NOT_APPLICABLE,
});
assert.equal(requiredPersonalization.route, "seller_customization_required");
assert.equal(requiredPersonalization.generateCustomizationMaster, true);
assert.equal(requiredPersonalization.generateCustomizationGuide, true);
assert.equal(requiredPersonalization.allowSellerUploadLabels, true);

const optionalPersonalization = resolveListingCustomizationFlow({
  capability: CUSTOMIZATION_CAPABILITY.YES,
  sellerService: SELLER_CUSTOMIZATION_SERVICE.OPTIONAL,
  nonSellerPositioning: NON_SELLER_POSITIONING.NOT_APPLICABLE,
});
assert.equal(optionalPersonalization.route, "seller_customization_optional");
assert.equal(optionalPersonalization.defaultDeliveryState, "buyer_selects_blank_or_personalized");
assert.equal(optionalPersonalization.generateCustomizationGuide, true);

const ambiguousCapableProduct = resolveListingCustomizationFlow({
  capability: CUSTOMIZATION_CAPABILITY.YES,
  sellerService: SELLER_CUSTOMIZATION_SERVICE.PENDING,
});
assert.equal(ambiguousCapableProduct.route, "clarify_seller_service");
assert.equal(ambiguousCapableProduct.clarificationRequired, true);
assert.equal(ambiguousCapableProduct.allowProductExampleGeneration, false);
assert.deepEqual(ambiguousCapableProduct.clarificationChoices, [
  "1. 不提供定制",
  "2. 提供定制，买家必须提交定制内容",
  "3. 提供定制，买家可以选择是否定制",
]);
assert.equal(ambiguousCapableProduct.generateCustomizationGuide, false);

const incapableProduct = resolveListingCustomizationFlow({
  capability: CUSTOMIZATION_CAPABILITY.NO,
  sellerService: SELLER_CUSTOMIZATION_SERVICE.PENDING,
});
assert.equal(incapableProduct.route, "ordinary_finished");
assert.equal(incapableProduct.clarificationRequired, false);
assert.equal(incapableProduct.allowProductExampleGeneration, true);

const contradictoryInput = resolveListingCustomizationFlow({
  capability: CUSTOMIZATION_CAPABILITY.NO,
  sellerService: SELLER_CUSTOMIZATION_SERVICE.REQUIRED,
});
assert.equal(contradictoryInput.route, "conflict");
assert.match(contradictoryInput.conflict, /无定制能力/);
assert.equal(contradictoryInput.generateCustomizationGuide, false);
assert.equal(contradictoryInput.allowProductExampleGeneration, false);

for (const resolvedFlow of [
  ordinaryFivePack,
  diyBlankFivePack,
  requiredPersonalization,
  optionalPersonalization,
]) {
  assert.equal(resolvedFlow.allowProductExampleGeneration, true);
}

console.log("Listing 定制流程检查通过：普通成品、DIY 空白基底、定制必选、定制可选、待确认与冲突分支均符合预期");
