export const CUSTOMIZATION_CAPABILITY = Object.freeze({
  YES: "yes",
  NO: "no",
  UNKNOWN: "unknown",
});

export const SELLER_CUSTOMIZATION_SERVICE = Object.freeze({
  NOT_OFFERED: "not_offered",
  REQUIRED: "required",
  OPTIONAL: "optional",
  PENDING: "pending",
});

export const NON_SELLER_POSITIONING = Object.freeze({
  ORDINARY: "ordinary_finished",
  DIY_BLANK: "diy_blank",
  NOT_APPLICABLE: "not_applicable",
  PENDING: "pending",
});

const labeledSellerServiceValues = new Map([
  ["不提供", SELLER_CUSTOMIZATION_SERVICE.NOT_OFFERED],
  ["不提供定制", SELLER_CUSTOMIZATION_SERVICE.NOT_OFFERED],
  ["必选", SELLER_CUSTOMIZATION_SERVICE.REQUIRED],
  ["必须定制", SELLER_CUSTOMIZATION_SERVICE.REQUIRED],
  ["可选", SELLER_CUSTOMIZATION_SERVICE.OPTIONAL],
  ["定制可选", SELLER_CUSTOMIZATION_SERVICE.OPTIONAL],
]);

const clarificationChoices = Object.freeze([
  "1. 不提供定制",
  "2. 必须定制",
  "3. 定制可选",
]);

const result = ({
  route,
  generateCustomizationMaster,
  generateCustomizationGuide,
  allowSellerUploadLabels,
  defaultDeliveryState,
  clarificationRequired = false,
  conflict = null,
}) => ({
  route,
  generateCustomizationMaster,
  generateCustomizationGuide,
  allowSellerUploadLabels,
  defaultDeliveryState,
  clarificationRequired,
  clarificationChoices: clarificationRequired ? [...clarificationChoices] : [],
  allowProductExampleGeneration: !clarificationRequired && conflict === null,
  conflict,
});

export function inferSellerCustomizationService(input = "") {
  const text = String(input).trim();
  if (!text) {
    return SELLER_CUSTOMIZATION_SERVICE.PENDING;
  }

  const labeledMatch = text.match(
    /卖家定制服务\s*[：:]\s*(不提供定制|不提供|必须定制|必选|定制可选|可选)(?:\s|$)/,
  );
  if (labeledMatch) {
    return labeledSellerServiceValues.get(labeledMatch[1]);
  }

  const candidates = new Set();
  if (
    /(?:本\s*Listing|当前\s*Listing|本商品)\s*(?:不提供|不启用|关闭)\s*(?:卖家)?定制/i.test(text) ||
    /卖家\s*(?:不提供|不承接)\s*定制/.test(text) ||
    /(?:按空白基底销售|买家自行定制)/.test(text)
  ) {
    candidates.add(SELLER_CUSTOMIZATION_SERVICE.NOT_OFFERED);
  }
  if (
    /(?:本\s*Listing|当前\s*Listing|本商品)\s*(?:必须定制|定制必选)/i.test(text) ||
    /卖家定制\s*(?:必须|为必选)/.test(text)
  ) {
    candidates.add(SELLER_CUSTOMIZATION_SERVICE.REQUIRED);
  }
  if (
    /(?:本\s*Listing|当前\s*Listing|本商品)\s*(?:定制可选|可选定制)/i.test(text) ||
    /(?:可选择定制或不定制|卖家定制\s*可选)/.test(text)
  ) {
    candidates.add(SELLER_CUSTOMIZATION_SERVICE.OPTIONAL);
  }

  return candidates.size === 1
    ? [...candidates][0]
    : SELLER_CUSTOMIZATION_SERVICE.PENDING;
}

export function resolveListingCustomizationFlow({
  capability = CUSTOMIZATION_CAPABILITY.UNKNOWN,
  sellerService = SELLER_CUSTOMIZATION_SERVICE.PENDING,
  nonSellerPositioning = NON_SELLER_POSITIONING.PENDING,
} = {}) {
  if (!Object.values(CUSTOMIZATION_CAPABILITY).includes(capability)) {
    throw new TypeError(`Unknown customization capability: ${capability}`);
  }
  if (!Object.values(SELLER_CUSTOMIZATION_SERVICE).includes(sellerService)) {
    throw new TypeError(`Unknown seller customization service: ${sellerService}`);
  }
  if (!Object.values(NON_SELLER_POSITIONING).includes(nonSellerPositioning)) {
    throw new TypeError(`Unknown non-seller positioning: ${nonSellerPositioning}`);
  }

  if (
    capability === CUSTOMIZATION_CAPABILITY.NO &&
    [SELLER_CUSTOMIZATION_SERVICE.REQUIRED, SELLER_CUSTOMIZATION_SERVICE.OPTIONAL].includes(sellerService)
  ) {
    return result({
      route: "conflict",
      generateCustomizationMaster: false,
      generateCustomizationGuide: false,
      allowSellerUploadLabels: false,
      defaultDeliveryState: "pending",
      conflict: "商品无定制能力，但本 Listing 启用了卖家定制服务",
    });
  }

  if (sellerService === SELLER_CUSTOMIZATION_SERVICE.PENDING) {
    if (capability === CUSTOMIZATION_CAPABILITY.NO) {
      return result({
        route: "ordinary_finished",
        generateCustomizationMaster: false,
        generateCustomizationGuide: false,
        allowSellerUploadLabels: false,
        defaultDeliveryState: "fixed_finished_product",
      });
    }
    return result({
      route: "clarify_seller_service",
      generateCustomizationMaster: false,
      generateCustomizationGuide: false,
      allowSellerUploadLabels: false,
      defaultDeliveryState: "pending",
      clarificationRequired: true,
    });
  }

  if (sellerService === SELLER_CUSTOMIZATION_SERVICE.NOT_OFFERED) {
    const isDiyBlank = nonSellerPositioning === NON_SELLER_POSITIONING.DIY_BLANK;
    return result({
      route: isDiyBlank ? "diy_blank" : "ordinary_finished",
      generateCustomizationMaster: false,
      generateCustomizationGuide: false,
      allowSellerUploadLabels: false,
      defaultDeliveryState: isDiyBlank ? "blank_base" : "fixed_finished_product",
    });
  }

  if (sellerService === SELLER_CUSTOMIZATION_SERVICE.REQUIRED) {
    return result({
      route: "seller_customization_required",
      generateCustomizationMaster: true,
      generateCustomizationGuide: true,
      allowSellerUploadLabels: true,
      defaultDeliveryState: "seller_personalized_product",
    });
  }

  return result({
    route: "seller_customization_optional",
    generateCustomizationMaster: true,
    generateCustomizationGuide: true,
    allowSellerUploadLabels: true,
    defaultDeliveryState: "buyer_selects_blank_or_personalized",
  });
}
