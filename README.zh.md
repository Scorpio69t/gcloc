
# gcloc: 源代码行统计工具

**[English](./README.md)** | **中文**

**gcloc** 是一个开源工具，用于统计各种编程语言的源代码文件数量和代码行数。它支持多种语言，易于扩展以包括自定义语言。`gcloc` 简单易用，可以帮助开发者快速了解代码库并跟踪变化。

---

## 致谢

本项目基于 [hhatto 的 gocloc 项目](https://github.com/hhatto/gocloc) 开发。我们诚挚感谢原作者为开源社区的贡献，这为 gcloc 的进一步开发提供了基础。

---

## 功能特点
- **语言支持**: 统计多个编程语言的文件数、空白行数、注释行数和代码行数。
- **可定制化**: 根据需要添加更多语言支持。
- **高性能**: 使用并发处理技术，高效地获取和分析文件，即使是大型代码库也能快速完成。
- **输出格式**: 支持多种输出格式，包括默认表格视图、XML、JSON 和 Sloccount 格式。
- **过滤功能**: 根据文件扩展名、正则表达式模式或特定语言包括或排除文件和目录。
- **详细统计**: 可选按文件报告统计结果。

---

## 未来开发计划

- [X] 支持 Git 仓库
    - 直接分析 Git 仓库，而无需指定目录。
    - 显示整个仓库的统计数据，包括提交历史。

- [X] 支持 Web 界面
    - 创建适配 gcloc 的 Web 界面，提供用户友好的使用体验。
    - 允许用户上传文件或目录以供分析。

- [ ] 支持压缩文件
    - 分析压缩文件（如 `.zip`, `.tar.gz`）而无需手动解压。
    - 动态解压并分析文件，提供准确的统计结果。

---

## 安装

使用以下命令安装最新版本：
```bash
go install github.com/Scorpio69t/gcloc/app/gcloc@latest
```

或者，从 [Releases](https://github.com/Scorpio69t/gcloc/releases) 页面下载预编译的二进制文件。

---

## 使用方法

```bash
gcloc [flags] PATH...
gcloc [command]
```

### 示例

分析当前目录：

```bash
gcloc .
```

**示例输出**:
```
$gcloc .
github.com/Scorpio69t/gcloc T=0.03 s (7318.7 files/s 2219353.8 lines/s)
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
C++ Header                       8           9494           3756          37956
C++                             28           3137           1752          13517
JSON                            43              1              0           9368
...
Total                          249          14710           7060          75508
-------------------------------------------------------------------------------
```

### 命令
#### 通用命令
- `completion`: 为当前 Shell 生成自动补全脚本。
- `help`: 显示命令帮助信息。
- `show-lang`: 列出所有支持的语言及其扩展名。
- `version`: 显示 gcloc 的版本信息。

#### 参数选项
| 参数                | 描述                                                                                |
|---------------------|--------------------------------------------------------------------------------------------|
| `--by-file`         | 报告每个源文件的统计结果。                                                                              |
| `--debug`           | 为开发者输出调试日志。                                                                                |
| `--exclude-ext`     | 排除指定扩展名的文件（逗号分隔）。                                                                          |
| `--exclude-lang`    | 排除指定语言（逗号分隔）。                                                                              |
| `--include-lang`    | 包括指定语言（逗号分隔）。                                                                              |
| `--match`           | 包括符合正则表达式的文件。                                                                              |
| `--not-match`       | 排除符合正则表达式的文件。                                                    |
| `--match-d`         | 包括符合正则表达式的目录。                                              |
| `--not-match-d`     | 排除符合正则表达式的目录。                                              |
| `--output-type`     | 输出格式：[default, gcloc-xml, sloccount, json]。 |
| `--skip-duplicated` | 跳过重复文件。                                                                  |
| `--sort`            | 按 name, files, blanks, comments, codes 排序（默认 codes）。      |

---

## 示例

### 分析特定目录
```bash
gcloc /path/to/code
```

### 仅包含特定语言
```bash
gcloc --include-lang "C++,JSON" .
```

### 排除特定扩展名的文件
```bash
gcloc --exclude-ext "json,xml" .
```

### 生成 JSON 格式的输出
```bash
gcloc --output-type json .
```

### 显示支持的语言
```bash
gcloc show-lang
```

### 启动 Web 服务器
```bash
gcloc web
```

---

## 贡献

欢迎贡献代码！贡献方式：
1. Fork 此仓库。
2. 为你的新功能或修复创建一个分支。
3. 提交一个详细说明的 Pull Request。

---

## 许可证

本项目基于 [MIT License](https://github.com/Scorpio69t/gcloc/blob/main/LICENSE) 开源。

---

## 致谢

特别感谢所有贡献者，帮助 gcloc 成为一个强大而可靠的源代码分析工具。
