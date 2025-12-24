底层存储选择minio
镜像和charts都存放在minio中，通过helm部署到k8s中
代码仓库使用gitlab，存储可以选择磁盘或者minio(迁移，备份，容灾，可靠)

minio扩容，支持单机扩容和分布式扩容

gitlab cicd 就可以实现自动构建镜像并推送到私有仓库，然后通过helm部署到k8s中

GitLab CI/CD
简介：这是 GitLab 自带的 CI/CD 工具，与 GitLab 紧密集成，无需额外配置复杂的连接。只要在项目仓库中添加 .gitlab-ci.yml 文件，就能定义 CI/CD 流程。
优点
无缝集成：和 GitLab 的代码仓库无缝结合，可直接基于代码仓库的操作（如提交、合并请求）触发 CI/CD 任务。
易于使用：配置文件采用 YAML 格式，简洁明了，易于编写和维护。
分布式执行：支持使用 GitLab Runner 进行分布式构建，可根据需求灵活扩展构建能力。
缺点：如果不使用 GitLab 作为代码托管平台，该工具的优势就难以发挥。
适用场景：适合已经使用 GitLab 进行代码管理的团队，能快速搭建起 CI/CD 流程。