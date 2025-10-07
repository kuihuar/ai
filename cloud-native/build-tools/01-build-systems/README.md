# 构建系统详解

## 📚 学习目标

通过本模块学习，您将掌握：
- 各种构建系统的工作原理和特点
- 构建缓存和增量构建策略
- 构建性能优化技巧
- 构建系统选择和迁移策略
- 企业级构建最佳实践

## 🎯 构建系统概览

### 1. 构建系统分类

```
构建系统分类
├── 语言特定构建工具
│   ├── Java: Maven, Gradle, Ant
│   ├── JavaScript: npm, yarn, pnpm, webpack
│   ├── Python: pip, poetry, setuptools
│   ├── Go: go build, go mod
│   ├── Rust: cargo
│   └── C/C++: make, cmake, ninja
├── 通用构建工具
│   ├── Bazel
│   ├── Buck
│   ├── Pants
│   └── Please
└── 容器构建工具
    ├── Docker
    ├── Buildah
    ├── Podman
    └── BuildKit
```

### 2. 构建系统特性对比

| 特性 | Maven | Gradle | Bazel | Webpack | Docker |
|------|-------|--------|-------|---------|--------|
| 增量构建 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 并行构建 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 缓存支持 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 依赖管理 | ✅ | ✅ | ✅ | ✅ | ❌ |
| 多语言支持 | ❌ | ✅ | ✅ | ❌ | ✅ |
| 学习曲线 | 简单 | 中等 | 陡峭 | 中等 | 简单 |

## 🏗️ 语言构建系统

### 1. Java 构建系统

#### Maven 详解
```xml
<!-- pom.xml 完整示例 -->
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>com.example</groupId>
    <artifactId>myapp</artifactId>
    <version>1.0.0</version>
    <packaging>jar</packaging>

    <name>My Application</name>
    <description>My Spring Boot Application</description>

    <properties>
        <maven.compiler.source>11</maven.compiler.source>
        <maven.compiler.target>11</maven.compiler.target>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
        <spring.boot.version>2.7.0</spring.boot.version>
        <junit.version>5.8.2</junit.version>
    </properties>

    <dependencyManagement>
        <dependencies>
            <dependency>
                <groupId>org.springframework.boot</groupId>
                <artifactId>spring-boot-dependencies</artifactId>
                <version>${spring.boot.version}</version>
                <type>pom</type>
                <scope>import</scope>
            </dependency>
        </dependencies>
    </dependencyManagement>

    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-data-jpa</artifactId>
        </dependency>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-test</artifactId>
            <scope>test</scope>
        </dependency>
    </dependencies>

    <build>
        <plugins>
            <plugin>
                <groupId>org.springframework.boot</groupId>
                <artifactId>spring-boot-maven-plugin</artifactId>
                <version>${spring.boot.version}</version>
                <executions>
                    <execution>
                        <goals>
                            <goal>repackage</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-compiler-plugin</artifactId>
                <version>3.10.1</version>
                <configuration>
                    <source>${maven.compiler.source}</source>
                    <target>${maven.compiler.target}</target>
                </configuration>
            </plugin>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-surefire-plugin</artifactId>
                <version>3.0.0-M7</version>
                <configuration>
                    <includes>
                        <include>**/*Test.java</include>
                    </includes>
                </configuration>
            </plugin>
        </plugins>
    </build>

    <profiles>
        <profile>
            <id>dev</id>
            <activation>
                <activeByDefault>true</activeByDefault>
            </activation>
            <properties>
                <spring.profiles.active>dev</spring.profiles.active>
            </properties>
        </profile>
        <profile>
            <id>prod</id>
            <properties>
                <spring.profiles.active>prod</spring.profiles.active>
            </properties>
        </profile>
    </profiles>
</project>
```

## 🚀 通用构建系统

### 1. Bazel 详解

#### WORKSPACE 文件
```python
# WORKSPACE
workspace(name = "myapp")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# 下载规则
http_archive(
    name = "rules_nodejs",
    sha256 = "5aef09ed3279aa01d5c928e3beb248f9ad32dde6aafe6753e8a74e8aec2c4005",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/5.5.3/rules_nodejs-5.5.3.tar.gz"],
)

load("@rules_nodejs//nodejs:repositories.bzl", "nodejs_register_toolchains")
nodejs_register_toolchains(
    name = "nodejs",
    node_version = "16.9.0",
)
```

## 🛠️ 实践练习

### 练习1: 多模块项目构建

```gradle
// settings.gradle
rootProject.name = 'myapp'
include 'core'
include 'web'
include 'api'

// build.gradle (root)
allprojects {
    group = 'com.example'
    version = '1.0.0'
}

subprojects {
    apply plugin: 'java'
    
    repositories {
        mavenCentral()
    }
    
    dependencies {
        testImplementation 'junit:junit:4.13.2'
    }
}
```

## 📚 相关资源

### 官方文档
- [Maven 官方文档](https://maven.apache.org/guides/)
- [Gradle 官方文档](https://docs.gradle.org/)
- [Bazel 官方文档](https://bazel.build/docs)

### 学习资源
- [构建系统最佳实践](https://martinfowler.com/articles/ci.html)
- [构建性能优化指南](https://bazel.build/configure/performance)

### 工具推荐
- **Maven**: Java 项目构建
- **Gradle**: 多语言项目构建
- **Bazel**: 大规模项目构建
- **Webpack**: JavaScript 模块打包
- **Vite**: 现代前端构建工具

---

**掌握构建系统，提升开发效率！** 🚀
