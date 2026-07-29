两种部署方式

  1. tar 导出方式（make deploy-ecs / scripts/deploy.sh）

  不使用镜像仓库，直接在当前机器构建 Docker 镜像，docker save 打包压缩后通过 scp 传到 ECS，在 ECS 上 docker load 并 docker compose up -d 启动。部署完成后自动调用 ISR revalidate 刷新所有缓存。

  2. Harbor 镜像仓库方式（make deploy / scripts/docker-deploy.sh）

  将构建好的镜像推送到 Harbor 私有仓库，然后 SSH 到 ECS docker pull 并启动。

  关键配置

  - Next.js 使用 output: "standalone" 模式，使其可以自托管运行
  - 两个 Docker 容器均以非 root 用户（appuser）运行
  - PocketBase 只对内暴露 8090 端口，不对外开放
  - 没有 CI/CD 流水线，部署完全是手动执行脚本
  - 没有使用 Vercel 等平台服务

✻ Sautéed for 1m 6s

❯ 如果PocketBase 只对内暴露 8090 端口，不对外开放， 那打开网站应用的network 会有http://127.0.0.1:8090/api/collections/pages/records?page=1&perPage=1&filter=title_key%3D%22footer_plain_items%22&skipTotal=1 这样子的请求？为什么？

  Thought for 10s (ctrl+o to expand)

⏺ Explore(Find PocketBase client usage)
  ⎿  Done (28 tool uses · 46.8k tokens · 2m 47s)
  (ctrl+o to expand)

  Thought for 46s, read 4 files (ctrl+o to expand)

⏺ 原因找到了。这是项目中客户端组件双重获取数据的模式导致的。

  数据流分析

  以首页为例：

  page.tsx:13（Server Component）在服务端获取全部数据，并通过 props 传递给子组件：

  // page.tsx (服务端, POCKETBASE_URL=http://pb:8090)
  const heroPage = await getPageContent("home_hero");
  return <HeroSection pageRecord={heroPage} />;  // 数据已作为 prop 传入

  HeroSection.tsx:16-18（Client Component）收到 prop 后，又自己调用了 hook 去获取数据：

  export default function HeroSection({ pageRecord }) {
    const { data: record } = usePageRecordData("home_hero");  // 又发了一次请求
    const finalRecord = pageRecord || record;  // 优先用 prop

  为什么浏览器能看到 127.0.0.1:8090 的请求

  usePageRecordData → useContent 是一个 "use client" hook，它的 useEffect 会在浏览器端执行 fetch。

  关键在 src/lib/pocketbase.ts:19：

  const PB_URL = process.env.POCKETBASE_URL ?? "http://127.0.0.1:8090";

  POCKETBASE_URL 没有 NEXT_PUBLIC_ 前缀，所以在浏览器端 process.env.POCKETBASE_URL 是 undefined，直接 fallback 到了 http://127.0.0.1:8090。PocketBase SDK 就从这个 URL 发请求，你在浏览器的 Network 面板里自然就看到了。

  三个问题

  1. 本地开发：PocketBase 确实跑在 127.0.0.1:8090，请求能成功，但你看到了冗余的双重请求（服务端拿了一次，客户端又拿了一次）。
  2. 生产环境：浏览器里的 127.0.0.1 是用户自己的电脑，不是 ECS 服务器。这些请求会失败（Connection Refused），触发 reportPbError，进而导致 PbStatusBanner 显示错误横幅。
  3. 浪费带宽：即使 pageRecord prop 已经有数据（finalRecord = pageRecord || record 优先用 prop），useEffect 还是会无条件触发一次多余的 fetch。

  本质原因

  这些 Client Component 设计了"自给自足"的 fallback 数据获取（通过 useContent hook），但它们的父级 Server Component 已经通过 props 把数据传下来了。客户端 fetch 变成了一次多余的、且在生产环境必然失败的请求。
