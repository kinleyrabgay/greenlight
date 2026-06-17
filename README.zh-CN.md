<h1 align="center"><code>git push greenlight</code></h1>
<p align="center">
  <a href="https://github.com/kinleyrabgay/greenlight/actions/workflows/release.yml"
    ><img
      alt="Release"
      src="https://img.shields.io/github/actions/workflow/status/kinleyrabgay/greenlight/release.yml?style=flat-square&label=release"
  /></a>
  <a href="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-blue?style=flat-square"
    ><img
      alt="Platform"
      src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-blue?style=flat-square"
  /></a>
  <a href="https://x.com/kinleyrabgay"
    ><img
      alt="X"
      src="https://img.shields.io/badge/X-@kinleyrabgay-black?style=flat-square"
  /></a>
  <a href="https://discord.gg/Wsy2NpnZDu"
    ><img
      alt="Discord"
      src="https://img.shields.io/discord/1439901831038763092?style=flat-square&label=discord"
  /></a>
</p>

<h3 align="center">干掉所有 slop，开出干净的 PR。</h3>

<p align="center"><a href="README.md">English</a> · <strong>简体中文</strong></p>

<p align="center">
  <img src="https://raw.githubusercontent.com/kinleyrabgay/greenlight/main/demo.gif" alt="greenlight demo" width="800" />
</p>

`greenlight` 在你真实的远端前面放了一个本地 git 代理。把分支推给 `greenlight` 而不是 `origin`，它会拉起一个用完即弃的 worktree，跑一条 AI 驱动的校验流水线，**只有每一项检查都通过后**才转发到上游，并自动开出一个干净的 PR。

- **不阻塞** —— 流水线在隔离的 worktree 里跑，不打断你手头的工作。
- **不挑 agent** —— 支持 `claude`、`codex`、`rovodev`、`opencode`、`pi`，或通过 `acpx` 用 `acp:<target>`。
- **agent 原生** —— `/greenlight` 既能让编码 agent 完成一个任务再过网关，也能直接为已提交的工作过网关：它跑完流水线、让流水线应用安全的修复，剩下的升级给你。
- **人始终说了算** —— 自动修复，还是逐条审查 findings，你决定。
- **默认就是干净 PR** —— 推送、开 PR、盯 CI、自动修复失败，一气呵成。

完整文档：<https://kinleyrabgay.github.io/greenlight/>

## 工作原理

```
        你的分支
            │  git push greenlight
            ▼
   ┌──────────────────────────────────────────────┐
   │  用完即弃的 worktree —— 你的工作原地不动        │
   │  review → test → docs → lint → push → PR → CI  │
   └──────────────────────────────────────────────┘
            │  每项检查变绿
            ▼
        干净的 PR，已替你开好
```

每一步要么自己通过，要么停下来给你一条 **finding** 让你处理。安全、机械性的修复会自动应用；任何牵涉到你**意图**的，都会升级给你来 **approve（批准）**、**fix（修复）** 或 **skip（跳过）**。在每项检查都变绿之前，没有任何东西会到达你真实的远端。

## 安装

```sh
curl -fsSL https://raw.githubusercontent.com/kinleyrabgay/greenlight/main/docs/install.sh | sh
```

Windows、Go install 以及从源码构建的说明，见[安装指南](https://kinleyrabgay.github.io/greenlight/start-here/installation/)。

## 快速上手

```sh
$ greenlight init
  ✓ Gate initialized

    repo  /Users/you/src/my-repo
    gate  greenlight → /Users/you/.greenlight/repos/abc123def456.git
  remote  git@github.com:you/my-repo.git
   skill  /greenlight installed for agents at user level

  Push through the gate with:
  git push greenlight <branch>

$ git checkout my-branch

# 在分支里干点活……

$ git push greenlight
  * Pipeline started

  Run greenlight to review.

$ greenlight
# 打开当前运行的 TUI
```

在 TUI 里你逐条处理 **finding**：**auto-fix** 类自动替你应用（或由你 approve 放行），**ask-user** 类需要你判断，由你 approve、fix 或 skip。每项检查变绿后，网关会把你的分支转发到上游并替你开好 PR —— 不用手动 `git push origin`，也不用手写 PR 正文。想让编码 agent 无人值守地走完同一套流程？用 `/greenlight`（见下文）。

## 触发网关的三种方式

每一处改动都走同一条流水线。改动就绪时，挑一个最贴合你当下工作方式的入口：

- **`git push greenlight`** —— 显式的 Git 路径。把已提交的分支推给网关 remote，而不是 `origin`。
- **`greenlight`** —— TUI。改完之后运行它（无需先提交），向导会带你建分支、提交、推过网关，然后挂到这次运行上。`greenlight -y` 会把这一切自动做完。
- **`/greenlight`** —— agent skill。用 `/greenlight <task>` 让编码 agent 完成一个任务再过网关，或用裸 `/greenlight` 为已提交的工作过网关。它跑完流水线、让流水线应用安全的修复，并在任何需要人来拍板的地方停下来问你。

`greenlight init` 会为 Claude Code 及其他 agent 安装 `/greenlight` skill。底层上这个 skill 驱动的是 `greenlight axi` —— 同一套审批流程的非交互式 TOON 接口。

完整的首次运行走查见[快速上手](https://kinleyrabgay.github.io/greenlight/start-here/quick-start/)。

## 开发

```sh
make build   # 构建 bin/greenlight（带版本信息）
make test    # 运行 go test -race ./...（不含 e2e 套件）
make e2e     # 运行打了标签的端到端 agent 旅程套件
make e2e-record # agent 线格式变化时，重新录制 e2e fixtures
make lint    # 检查生成的 skill 是否漂移，并跑 go vet ./...
make skill   # 重新生成已提交的 greenlight skill 文件
make fmt     # 运行 gofmt -w .
make demo    # 重新生成 demo.gif 和 demo.mp4（需要 vhs 和 ffmpeg）
make docs    # 在 docs/dist 构建 Astro 文档站
```

完整 target 列表见 `Makefile`。

`make e2e-record` 会用真实的 `claude`、`codex`、`opencode` CLI 覆盖 `internal/e2e/fixtures/`，会消耗真实 API 额度，提交前应当审查。

## Star 历史

<a href="https://www.star-history.com/?repos=kinleyrabgay%2Fgreenlight&type=date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/chart?repos=kinleyrabgay/greenlight&type=date&theme=dark&legend=top-left" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/chart?repos=kinleyrabgay/greenlight&type=date&legend=top-left" />
   <img alt="Star History Chart" src="https://api.star-history.com/chart?repos=kinleyrabgay/greenlight&type=date&legend=top-left" />
 </picture>
</a>
