---
date: 2020-11-15T20:29:40+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "gomod 深入讲解"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/docs/mod/10_module_go_sum"  # 设置网页永久链接
tags: [ "go", "gomod" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

为了确保一致性构建，Go引入了`go.mod`文件来标记每个依赖包的版本，在构建过程中`go`命令会下载`go.mod`中的依赖包，下载的依赖包会缓存在本地，以便下次构建。
考虑到下载的依赖包有可能是被黑客恶意篡改的，以及缓存在本地的依赖包也有被篡改的可能，单单一个`go.mod`文件并不能保证一致性构建。

为了解决Go module的这一安全隐患，Go开发团队在引入`go.mod`的同时也引入了`go.sum`文件，用于记录每个依赖包的哈希值，在构建时，如果本地的依赖包hash值与`go.sum`文件中记录得不一致，则会拒绝构建。

本节暂不对模块校验细节展开介绍，只从日常应用层面介绍：

- go.sum 文件记录含义
- go.sum文件内容是如何生成的
- go.sum是如何保证一致性构建的

## go.sum文件记录

`go.sum`文件中每行记录由`module`名、版本和哈希组成，并由空格分开：
```
<module> <version>[/go.mod] <hash>
```
比如，某个`go.sum`文件中记录了`github.com/google/uuid` 这个依赖包的`v1.1.1`版本的哈希值：
```
github.com/google/uuid v1.1.1 h1:Gkbcsh/GbpXz7lPftLA3P6TYMwjCLYm83jiFQZF/3gY=  
github.com/google/uuid v1.1.1/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
```
在Go module机制下，我们需要同时使用依赖包的名称和版本才可以准确的描述一个依赖，为了方便叙述，下面我们使用`依赖包版本`来指代依赖包名称和版本。

正常情况下，每个`依赖包版本`会包含两条记录，第一条记录为该`依赖包版本`整体（所有文件）的哈希值，第二条记录仅表示该`依赖包版本`中`go.mod`文件的哈希值，如果该`依赖包版本`没有`go.mod`文件，则只有第一条记录。如上面的例子中，`v1.1.1`表示该`依赖包版本`整体，而`v1.1.1/go.mod`表示该`依赖包版本`中`go.mod`文件。

`依赖包版本`中任何一个文件（包括`go.mod`）改动，都会改变其整体哈希值，此处再额外记录`依赖包版本`的`go.mod`文件主要用于计算依赖树时不必下载完整的`依赖包版本`，只根据`go.mod`即可计算依赖树。

每条记录中的哈希值前均有一个表示哈希算法的`h1:`，表示后面的哈希值是由算法`SHA-256`计算出来的，自Go module从v1.11版本初次实验性引入，直至v1.14 ，只有这一个算法。

此外，细心的读者或许会发现`go.sum`文件中记录的`依赖包版本`数量往往比`go.mod`文件中要多，这是因为二者记录的粒度不同导致的。`go.mod`只需要记录直接依赖的`依赖包版本`，只在`依赖包版本`不包含`go.mod`文件时候才会记录间接`依赖包版本`，而`go.sum`则是要记录构建用到的所有`依赖包版本`。

## 生成

假设我们在开发某个项目，当我们在GOMODULE模式下引入一个新的依赖时，通常会使用`go get`命令获取该依赖，比如：
```
go get github.com/google/uuid@v1.0.0
```
`go get`命令首先会将该依赖包下载到本地缓存目录`$GOPATH/pkg/mod/cache/download`，该依赖包为一个后缀为`.zip`的压缩包，如`v1.0.0.zip`。`go get`下载完成后会对该`.zip`包做哈希运算，并将结果存放在后缀为`.ziphash`的文件中，如`v1.0.0.ziphash`。如果在项目的根目录中执行`go get`命令的话，`go get`会同步更新`go.mod`和`go.sum`文件，`go.mod`中记录的是依赖名及其版本，如：
```
require (
	github.com/google/uuid v1.0.0
)
```
`go.sum`文件中则会记录依赖包的哈希值（同时还有依赖包中go.mod的哈希值），如：
```
github.com/google/uuid v1.0.0 h1:b4Gk+7WdP/d3HZH8EJsZpvV7EtDOgaZLtnaNGIu1adA=
github.com/google/uuid v1.0.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
```

值得一提的是，在更新`go.sum`之前，为了确保下载的依赖包是真实可靠的，`go`命令在下载完依赖包后还会查询`GOSUMDB`环境变量所指示的服务器，以得到一个权威的`依赖包版本`哈希值。如果`go`命令计算出的`依赖包版本`哈希值与`GOSUMDB`服务器给出的哈希值不一致，`go`命令将拒绝向下执行，也不会更新`go.sum`文件。

`go.sum`存在的意义在于，我们希望别人或者在别的环境中构建当前项目时所使用依赖包跟`go.sum`中记录的是完全一致的，从而达到一致构建的目的。

## 校验

假设我们拿到某项目的源代码并尝试在本地构建，`go`命令会从本地缓存中查找所有`go.mod`中记录的依赖包，并计算本地依赖包的哈希值，然后与`go.sum`中的记录进行对比，即检测本地缓存中使用的`依赖包版本`是否满足项目`go.sum`文件的期望。

如果校验失败，说明本地缓存目录中`依赖包版本`的哈希值和项目中`go.sum`中记录的哈希值不一致，`go`命令将拒绝构建。
这就是`go.sum`存在的意义，即如果不使用我期望的版本，就不能构建。

当校验失败时，有必要确认到底是本地缓存错了，还是`go.sum`记录错了。
需要说明的是，二者都可能出错，本地缓存目录中的`依赖包版本`有可能被有意或无意地修改过，项目中`go.sum`中记录的哈希值也可能被篡改过。

当校验失败时，`go`命令倾向于相信`go.sum`，因为一个新的`依赖包版本`在被添加到`go.sum`前是经过`GOSUMDB`（校验和数据库）验证过的。此时即便系统中配置了`GOSUMDB`（校验和数据库），`go`命令也不会查询该数据库。

## 校验和数据库

环境变量`GOSUMDB`标识一个`checksum database`，即校验和数据库，实际上是一个web服务器，该服务器提供查询`依赖包版本`哈希值的服务。

该数据库中记录了很多`依赖包版本`的哈希值，比如Google官方的`sum.golang.org`则记录了所有的可公开获得的`依赖包版本`。除了使用官方的数据库，还可以指定自行搭建的数据库，甚至干脆禁用它（`export GOSUMDB=off`）。

如果系统配置了`GOSUMDB`，在`依赖包版本`被写入`go.sum`之前会向该数据库查询该`依赖包版本`的哈希值进行二次校验，校验无误后再写入`go.sum`。

如果系统禁用了`GOSUMDB`，在`依赖包版本`被写入`go.sum`之前则不会进行二次校验，`go`命令会相信所有下载到的依赖包，并把其哈希值记录到`go.sum`中。
