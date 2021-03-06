---
date: 2020-11-16T21:52:09+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 图片处理之合成、拼接"  # 文章标题
url:  "posts/go/io/image/merge"  # 设置网页永久链接
tags: [ "go", "图片处理"]  # 标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 标准库实现合成

```go
//
// Created by Rustle Karl on 2020.11.18 16:07.
//

package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

func main() {
	err := OverlayImage("storage/bg.jpg", "storage/fg.jpg")
	fmt.Println(os.Getwd())
	if err != nil {
		panic(err)
	}
}

func OverlayImage(dst, src string) error {
	imgDst, err := os.Open(dst)
	if err != nil {
		return err
	}
	defer imgDst.Close()

	imgDstDec, err := jpeg.Decode(imgDst)
	if err != nil {
		return err
	}

	imgSrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer imgSrc.Close()

	imgSrcDec, err := jpeg.Decode(imgSrc)
	if err != nil {
		return err
	}

	bound := imgDstDec.Bounds()
	rgba := image.NewRGBA(bound)
	draw.Draw(rgba, rgba.Bounds(), imgDstDec, image.Point{}, draw.Src)

	overlayImage(rgba, imgSrcDec.Bounds(), imgSrcDec)

	out, err := os.Create("storage/out.jpg")

	return jpeg.Encode(out, rgba, &jpeg.Options{Quality: 100})
}

func overlayImage(dst draw.Image, r image.Rectangle, src image.Image) image.Image {
	draw.Draw(dst, r, src, image.Pt(-100, -100), draw.Src)
	return dst
}
```

## 图片合成

```shell
go get -u github.com/disintegration/imaging
```

```go
package main

import (
    "fmt"
    "image"

    "github.com/disintegration/imaging"
)

func HandleUserImage(fileName string) (string, error) {
    m, err := imaging.Open("target.jpg")
    if err != nil {
        fmt.Printf("open file failed")
    }

    bm, err := imaging.Open("bg.jpg")
    if err != nil {
        fmt.Printf("open file failed")
    }

    // 图片按比例缩放
    dst := imaging.Resize(m, 200, 200, imaging.Lanczos)
    // 将图片粘贴到背景图的固定位置
    result := imaging.Overlay(bm, dst, image.Pt(120, 140), 1)

    fileName := fmt.Sprintf("%d.jpg", fileName)
    err = imaging.Save(result, fileName)
    if err != nil {
        return "", err
    }

    return fileName, nil
}
```

## 添加文字

```go
package main

import (
    "fmt"
    "image"
    "image/color"
    "io/ioutil"

    "github.com/disintegration/imaging"
    "github.com/golang/freetype"
    "github.com/golang/freetype/truetype"
    "golang.org/x/image/font"
)

func main() {
    HandleUserImage()
}

// HandleUserImage paste user image onto background
func HandleUserImage() (string, error) {
    m, err := imaging.Open("target.png")
    if err != nil {
        fmt.Printf("open file failed")
    }

    bm, err := imaging.Open("bg.jpg")
    if err != nil {
        fmt.Printf("open file failed")
    }

    // 图片按比例缩放
    dst := imaging.Resize(m, 200, 200, imaging.Lanczos)
    // 将图片粘贴到背景图的固定位置
    result := imaging.Overlay(bm, dst, image.Pt(120, 140), 1)
    writeOnImage(result)

    fileName := fmt.Sprintf("%d.jpg", 1234)
    err = imaging.Save(result, fileName)
    if err != nil {
        return "", err
    }

    return fileName, nil
}

var dpi = flag.Float64("dpi", 256, "screen resolution")

func writeOnImage(target *image.NRGBA) {
    c := freetype.NewContext()

    c.SetDPI(*dpi)
    c.SetClip(target.Bounds())
    c.SetDst(target)
    c.SetHinting(font.HintingFull)

        // 设置文字颜色、字体、字大小
    c.SetSrc(image.NewUniform(color.RGBA{R: 240, G: 240, B: 245, A: 180}))
    c.SetFontSize(16)
    fontFam, err := getFontFamily()
    if err != nil {
        fmt.Println("get font family error")
    }
    c.SetFont(fontFam)

    pt := freetype.Pt(500, 400)

    _, err = c.DrawString("我是水印", pt)
    if err != nil {
        fmt.Printf("draw error: %v \n", err)
    }

}

func getFontFamily() (*truetype.Font, error) {
        // 这里需要读取中文字体，否则中文文字会变成方格
    fontBytes, err := ioutil.ReadFile("Hei.ttc")
    if err != nil {
        fmt.Println("read file error:", err)
        return &truetype.Font{}, err
    }

    f, err := freetype.ParseFont(fontBytes)
    if err != nil {
        fmt.Println("parse font error:", err)
        return &truetype.Font{}, err
    }

    return f, err
```
