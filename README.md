# 项目名称

cryptLocal

## 项目简介

cryptLocal 是一款用于加密/解密(配合 Alist 的 Crypt 驱动)本地文件的命令行工具。

## 功能特性

- **基于 cobra**：支持交互式命令行工具，可通过命令行参数或交互式 prompt 输入参数。
- **自动获取加密配置**：手动输入或从 Alist 获取加密配置。
- **加密操作**：新建加密文件、新建加密文件夹
- **解密操作**：解密文件、解密文件夹
- **文件加密列表**：列出当前路径下文件列表。

## 使用说明

参考

```
cryptLocal --help
```

## 技术栈

- **编程语言**：Go 1.23.2
- **框架**：Cobra
- **交互**: promptui
- **加解密**: rclone

## 开发环境

### 克隆项目到本地

```bash
git clone https://github.com/cybenc/cryptLocal.git
```

### 下载依赖

```bash
go mod tidy
```

## 项目说明

项目仍然处于开发阶段,欢迎提出宝贵建议。项目目前存在诸多 bug。暂不发布 release 版本。有需要可自行下载编译。

## 许可证

本项目遵循 [MIT](https://opensource.org/licenses/MIT) 许可证。请查看 [LICENSE](https://github.com/cybenc/cryptLocal/blob/master/LICENSE) 文件以获取更多信息。
