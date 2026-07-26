import assert from "node:assert/strict";
import {
  CUSTOMIZATION_CAPABILITY,
  NON_SELLER_POSITIONING,
  SELLER_CUSTOMIZATION_SERVICE,
  resolveListingCustomizationFlow,
} from "./listing-customization-flow.mjs";

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
assert.deepEqual(ambiguousCapableProduct.clarificationChoices, [
  "1. 不提供定制",
  "2. 必须定制",
  "3. 定制可选",
]);
assert.equal(ambiguousCapableProduct.generateCustomizationGuide, false);

const incapableProduct = resolveListingCustomizationFlow({
  capability: CUSTOMIZATION_CAPABILITY.NO,
  sellerService: SELLER_CUSTOMIZATION_SERVICE.PENDING,
});
assert.equal(incapableProduct.route, "ordinary_finished");
assert.equal(incapableProduct.clarificationRequired, false);

const contradictoryInput = resolveListingCustomizationFlow({
  capability: CUSTOMIZATION_CAPABILITY.NO,
  sellerService: SELLER_CUSTOMIZATION_SERVICE.REQUIRED,
});
assert.equal(contradictoryInput.route, "conflict");
assert.match(contradictoryInput.conflict, /无定制能力/);
assert.equal(contradictoryInput.generateCustomizationGuide, false);

console.log("Listing 定制流程检查通过：普通成品、DIY 空白基底、定制必选、定制可选、待确认与冲突分支均符合预期");
