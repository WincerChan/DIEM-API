# Hitokoto

[![Build Status](https://travis-ci.org/WincerChan/Hitokoto.svg?branch=master)](https://travis-ci.org/WincerChan/Hitokoto)

~~接口地址： `https://api.itswincer.com/hitokoto/get?encode=`~~ 

新的接口地址：https://api.itswincer.com/hitokoto/v2/，旧版本的地址会 301 至新版本。

> 注意：所有请求的类型**仅支持** "GET" 

|   参数   |           含义            |    默认     |
| :------: | :-----------------------: | :---------: |
|  length  | `int` 限制查询一言的长度  |     无      |
|  encode  |  `string` 返回的数据格式  | 主体 + 出处 |
| charset  |    `string` 指定字符集    |    utf-8    |
| callback | `string` 指定回调函数名称 |     无      |

## 编码格式（encode）：

- js：JavaScript 脚本，将一言插入 HTML 中第一次出现 `class = 'hitokoto'` 的标签中
- json：JSON 格式的字符串，包含主体（hitokoto），出处（source）
- text：一言句子的主体
- 默认为：`×××××——「×××」`

## 字符集（charset）：

- utf-8（默认）：在 Header 中添加 `content-type: utf-8`
- gbk：在 Header 中添加 `content-type: gbk`

## 长度（length）：

值为数字，`length=40` 表示将查询一言句子的长度限制在 `40` 字符内

## 回调（callback）：

值为 JavaScript 的合法函数名，返回的函数参数是一个字典（dict），值是 `hitokoto` 和 `source`

> **注意：callback 参数会覆盖掉 encode 参数**

## 例子

```bash
curl 'https://api.itswincer.com/hitokoto/v2/'
幸福是生生不息，却难以触及的远。——「樱桃之远」

curl 'https://api.itswincer.com/hitokoto/v2/?encode=js&length=5'
var hitokoto="哦~\n——「袴田日向」";var dom=document.querySelector('.hitokoto');Array.isArray(dom)?dom[0].innerText=hitokoto:dom.innerText=hitokoto;

curl 'https://api.itswincer.com/hitokoto/v2/?encode=js&length=10&callback=hak'
hak({"hitokoto":"我在。","source":"凤囚凰"});
```



