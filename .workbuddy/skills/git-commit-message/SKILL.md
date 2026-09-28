---
name: git-commit-message
slug: git-commit-message
version: 1.0.0
displayName: "Git 提交消息（桥接到 Cursor 原规则）"
summary: "撰写本仓库 Git 提交消息，使用 Sign: Message（New/Chg/Enh/Bug/Doc），仅基于 staged diff。规则单一来源为 Cursor 版 skill，本技能只做桥接引用，不重写。"
description: "当用户说 撰写提交消息、git message、commit message、提交信息、写 commit、生成 git message 时触发。本技能不重复定义规则，而是先全量 Read 本仓库 .cursor/skills/git-commit-message/SKILL.md（Cursor 版，唯一事实源），严格按其中的 Sign: Message 格式与「仅处理 staged 文件」等注意事项执行；撰写后提示用户确认再提交。"
tags: ["git", "commit", "提交消息", "cursor-bridge"]
agent_created: true
---

# Git 提交消息（WorkBuddy 桥接技能）

本技能是 Cursor 版 `git-commit-message` 的 **WorkBuddy 桥接层**。**不重写规则，只引用。**

## 单一来源（禁止内联拷贝规则正文）

- 规则唯一事实源：`<项目根目录>/.cursor/skills/git-commit-message/SKILL.md`（Cursor 版）。
- 项目根目录定义：同时包含 `.cursor` 与 `.workbuddy` 两个子目录的那一层（即 `.workbuddy` 的直接父目录，也是 `AGENTS.md` 所在目录）。Windows / macOS / Linux 路径不同，禁止写死绝对路径，一律按此规则解析。
- 本文件保持纯桥接：若 Cursor 版规则后续更新，本技能无需改动即可同步生效——这正是「引用而非副本」的目的。

## 执行流程（每次触发都执行）

1. 用 Read 工具 **全量读取** `<项目根目录>/.cursor/skills/git-commit-message/SKILL.md`。
2. 严格按刚刚读取的 Cursor 版 SKILL.md 中的**全部规则**撰写提交消息并执行提交流程（消息格式、Sign 取值、仅处理 staged、staged 为空的处理、提交前格式化、确认后再提交等）。本文件不重复罗列这些规则。
3. 若 `.cursor` 版规则与本文有任何出入，以 `.cursor` 版为准（它才是事实源）。

## 维护约定

- 规则若要修改，请在 `.cursor/skills/git-commit-message/SKILL.md` 单一来源处修改。
- 本文件不得内联拷贝规则正文；如需调整触发词，只改本文件 frontmatter 的 `description`，不改变指向。
