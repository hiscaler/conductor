import { pathToFileURL } from "node:url";

export const START_MENU_SCOPES = Object.freeze({
  1: ["copy", "product_input"],
  2: ["copy", "images", "product_input"],
  3: ["copy", "images", "video_script", "product_input"],
  4: ["discovery"],
  5: ["product_evaluation"],
  6: ["platform_strategy"],
  7: ["competitors"],
  8: ["images", "product_input"],
  9: ["video_script", "video", "product_input"],
  10: ["listing_qa"],
  11: ["growth_review"],
  12: ["asset_view"],
  13: ["*"],
});

export const ACTION_SCOPES = Object.freeze({
  save_copy_assets: "copy",
  generate_main_image: "images",
  generate_full_image_set: "images",
  save_image_plan: "images",
  generate_missing_images: "images",
  generate_listing_package: "listing_package",
  run_listing_qa: "listing_qa",
  research_competitors: "competitors",
  calculate_profit: "profit",
  review_compliance: "compliance",
  create_test_plan: "growth_review",
  save_video_script: "video_script",
  generate_main_video: "video",
  generate_video_set: "video",
  discover_trend_video_products: "discovery",
  research_market_and_reviews: "competitors",
  regenerate_product_example: "product_input",
  research_supply_chain: "supply_chain",
  view_full_production_report: "always",
  revise_copy_assets: "copy",
  regenerate_selected_images: "images",
  regenerate_current_scope: "always",
  revise_video_script: "video_script",
  view_current_assets: "always",
  refine_discovery_criteria: "discovery",
  continue_product_evaluation: "product_evaluation",
  revise_platform_comparison: "platform_strategy",
  update_competitor_analysis: "competitors",
  revise_listing_qa: "listing_qa",
  update_growth_review: "growth_review",
  save_current_assets: "asset_view",
});

const normalizeMenuSelections = (menuSelections) => {
  const values = Array.isArray(menuSelections)
    ? menuSelections
    : String(menuSelections)
        .split(",")
        .map((value) => value.trim())
        .filter(Boolean);

  if (values.length === 0) throw new Error("至少需要一个启动菜单编号");
  return values.map((value) => {
    const menu = Number(value);
    if (!Number.isInteger(menu) || !START_MENU_SCOPES[menu]) {
      throw new Error(`未知启动菜单编号：${value}`);
    }
    return menu;
  });
};

export const validateNextActions = (menuSelections, actionIds, options = {}) => {
  const menus = normalizeMenuSelections(menuSelections);
  if (!Array.isArray(actionIds) || actionIds.length === 0) {
    throw new Error("下一步动作不能为空");
  }
  if (actionIds.length > 5) {
    throw new Error("下一步动作最多展示 5 个");
  }

  const scopes = new Set(options.extraScopes ?? []);
  for (const menu of menus) {
    for (const scope of START_MENU_SCOPES[menu]) scopes.add(scope);
  }

  const rejected = [];
  for (const actionId of actionIds) {
    const scope = ACTION_SCOPES[actionId];
    if (!scope) throw new Error(`未知下一步动作：${actionId}`);
    if (!scopes.has("*") && scope !== "always" && !scopes.has(scope)) {
      rejected.push(`${actionId}（需要 ${scope} 范围）`);
    }
  }

  if (rejected.length > 0) {
    throw new Error(
      `下一步动作超出启动菜单 ${menus.join(",")} 的任务范围：${rejected.join("、")}`,
    );
  }

  return { menus, scopes: [...scopes], actions: actionIds };
};

const runCli = () => {
  const [menus, ...actions] = process.argv.slice(2);
  try {
    const result = validateNextActions(menus, actions);
    process.stdout.write(`${JSON.stringify(result, null, 2)}\n`);
  } catch (error) {
    process.stderr.write(`下一步动作校验失败：${error.message}\n`);
    process.exitCode = 1;
  }
};

if (process.argv[1] && pathToFileURL(process.argv[1]).href === import.meta.url) {
  runCli();
}
