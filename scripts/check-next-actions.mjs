import assert from "node:assert/strict";
import {
  ACTION_SCOPES,
  START_MENU_SCOPES,
  validateNextActions,
} from "./next-action-scope.mjs";

const menuTwoActions = [
  "revise_copy_assets",
  "regenerate_selected_images",
  "regenerate_current_scope",
  "view_current_assets",
];
assert.doesNotThrow(() => validateNextActions("2", menuTwoActions));

for (const unrelatedAction of [
  "run_listing_qa",
  "generate_listing_package",
  "research_supply_chain",
  "save_video_script",
  "calculate_profit",
]) {
  assert.throws(
    () => validateNextActions("2", [unrelatedAction]),
    /超出启动菜单 2 的任务范围/,
    `菜单 2 不得推荐 ${unrelatedAction}`,
  );
}

assert.doesNotThrow(() =>
  validateNextActions("2,7", ["revise_copy_assets", "research_competitors"]),
);
assert.doesNotThrow(() => validateNextActions("3", ["revise_video_script"]));
assert.throws(() => validateNextActions("2", ["revise_video_script"]), /超出启动菜单 2/);
assert.doesNotThrow(() =>
  validateNextActions("13", [
    "run_listing_qa",
    "generate_listing_package",
    "research_supply_chain",
    "save_video_script",
    "calculate_profit",
  ]),
);

for (let menu = 1; menu <= 12; menu += 1) {
  const allowedScopes = new Set(START_MENU_SCOPES[menu]);
  for (const [actionId, actionScope] of Object.entries(ACTION_SCOPES)) {
    const shouldAllow = actionScope === "always" || allowedScopes.has(actionScope);
    if (shouldAllow) {
      assert.doesNotThrow(
        () => validateNextActions(String(menu), [actionId]),
        `菜单 ${menu} 应允许同范围动作 ${actionId}`,
      );
    } else {
      assert.throws(
        () => validateNextActions(String(menu), [actionId]),
        /超出启动菜单/,
        `菜单 ${menu} 必须拒绝跨范围动作 ${actionId}`,
      );
    }
  }
}

assert.doesNotThrow(() => validateNextActions("4", ["refine_discovery_criteria"]));
assert.doesNotThrow(() => validateNextActions("5", ["continue_product_evaluation"]));
assert.doesNotThrow(() => validateNextActions("6", ["revise_platform_comparison"]));
assert.doesNotThrow(() => validateNextActions("7", ["update_competitor_analysis"]));
assert.doesNotThrow(() => validateNextActions("10", ["revise_listing_qa"]));
assert.doesNotThrow(() => validateNextActions("11", ["update_growth_review"]));
assert.doesNotThrow(() => validateNextActions("12", ["save_current_assets"]));

process.stdout.write("下一步动作范围检查通过\n");
