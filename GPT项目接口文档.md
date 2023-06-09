# GPT项目接口文档

## 一些流程说明

```
从整体上我们将网页的内容分为两部分，一部分是实验，一部分是与chatgpt自由聊天的部分（与实验本身的聊天环节区分开来），用户登录后都先通过后端接口判断是否有实验，有则进入实验流程，无则直接进入聊天，实验流程结束也同样进入聊天
```

![流程图](https://i.ibb.co/4F1Y3V7/image.png "实验流程")

## 域名

测试：http://119.91.192.119:8089

## 接口列表

- [注册登录](#注册登录)
- [获取实验](#获取实验)
- [ChatGPT聊天](#ChatGPT聊天)
- [ChatGPT聊天事件流](#ChatGPT聊天事件流)
- [实验提交答案](#实验提交答案)
- [问卷提交答案](#问卷提交答案)
- [获取页面基础信息](#获取页面基础信息)

## 注册登录

**接口描述：** 登录注册

**请求URL：** /api/account/login

**请求方式：** POST

### 请求参数：

| 参数名      | 必选  | 类型     | 说明       |
| -------- | --- | ------ | -------- |
| email    | 是   | string | 邮箱       |
| password | 是   | string | password |

### 返回示例：

```json
{
    "code": 10000,
    "info": {
        "token": "D01B0E972A264FA08EA8F12EF1C0DAB6",
    },
    "message": "成功"
}
```

### 返回参数说明：

| 参数名   | 类型     | 说明         |
| ----- | ------ | ---------- |
| token | string | 生成的登录token |

### 特殊说明：

```markdown
当返回code为-10001时表示token过期需求重新跳转到登录注册页面
除登录接口外，其余接口均需带上登录token进行访问，携带方式为放置在请求header中的X-Token字段中字段：
X-Token: D01B0E972A264FA08EA8F12EF1C0DAB6
```

## 获取实验

**接口描述：** 获取实验相关数据，其中实验包括提示信息页，实验的主体问题（实验中的聊天部分），问卷三部分，无需传递参数，是否有实验以及返回什么实验完全由后端根据规则判定，前端只需要在有返回的实验数据的情况下进行展示即可，**调用时间：登录后以及后台每隔一段时间定时调用，比如每10分钟轮询调用一次**

**请求URL：** /api/experiment

**请求方式：** GET

### 请求参数：

| 参数名 | 必选 | 类型 | 说明                                       |
| ------ | ---- | ---- | ------------------------------------------ |
| id     | 否   | int  | 实验id，用于测试，有传则返回相应的id的实验 |

### 返回示例：

```json
{
    "code": 10000,
    "experiment": {
        "id": 1,
        "experiment_topic": "实验主题或者说实验要回答的创新问题，在实验的聊天页面显示",
        "guide_pages": [
            {
                "id": 1,
                "content": "欢迎光临CHAT小屋xxxxxx",
                "next_button": 1,
                "yes_or_no": 0,
                "chat": 0,
                "chat_times": 0,
                "next_page": 2,
                "yes_page": 0,
                "no_page": 0,
                "word_limit_min": 0,
                "word_limit_max": 0,
                "countdown": 20,
                "answer_page": 0,
                "answer_time_countdown": 0,
                "get_chat_from": 0,
                "save_chat_info": 0,
                "chat_tips":""
            },
            {
                "id": 2,
                "content": "page2xxxxx",
                "next_button": 0,
                "yes_or_no": 1,
                "chat": 0,
                "chat_times": 0,
                "next_page": 0,
                "yes_page": 3,
                "no_page": 4,
                "word_limit_min": 0,
                "word_limit_max": 0,
                "countdown": 20,
                "answer_page": 0,
                "answer_time_countdown": 0,
                "get_chat_from": 0,
                "save_chat_info": 0,
                "chat_tips":""
            },
            {
                "id": 3,
                "content": "xxxx",
                "next_button": 1,
                "yes_or_no": 0,
                "chat": 0,
                "chat_times": 0,
                "next_page": 5,
                "yes_page": 0,
                "no_page": 0,
                "word_limit_min": 0,
                "word_limit_max": 0,
                "countdown": 20,
                "answer_page": 0,
                "answer_time_countdown": 0,
                "get_chat_from": 0,
                "save_chat_info": 0,
                "chat_tips":""
            },
            {
                "id": 4,
                "content": "xxxx",
                "next_button": 1,
                "yes_or_no": 0,
                "chat": 0,
                "chat_times": 0,
                "next_page": 5,
                "yes_page": 0,
                "no_page": 0,
                "word_limit_min": 0,
                "word_limit_max": 0,
                "countdown": 20,
                "answer_page": 0,
                "answer_time_countdown": 0,
                "get_chat_from": 0,
                "save_chat_info": 0,
                "chat_tips":""
            },
            {
                "id": 5,
                "content": "xxxx",
                "next_button": 1,
                "yes_or_no": 0,
                "chat": 1,
                "chat_times": 1,
                "next_page": 6,
                "yes_page": 0,
                "no_page": 0,
                "word_limit_min": 0,
                "word_limit_max": 0
                "countdown": 0,
                "answer_page": 1,
                "answer_time_countdown": 120,
                "get_chat_from": 0,
                "save_chat_info": 1,
                "chat_tips":"xxxxx"
            },
            {
                "id": 6,
                "content": "xxxx",
                "next_button": 0,
                "yes_or_no": 0,
                "chat": 1,
                "chat_times": 0,
                "next_page": 0,
                "yes_page": 0,
                "no_page": 0,
                "word_limit_min": 100,
                "word_limit_max": 500,
                "countdown": 0,
                "answer_page": 1,
                "answer_time_countdown": 240,
                "get_chat_from": 5,
                "save_chat_info": 0,
                "chat_tips":"xxx"
            }
        ],
        "questionnaire": {
            "id": 1,
            "title": "问卷标题",
            "questions": [
                {
                    "id": 1,
                    "content": "题目内容",
                    "choice": [
                        "安踏",
                        "李宁",
                        "特步"
                    ],
                    "type": 1,
                    "score_type": 0,
                    "score_text": [
                        "非常不同意",
                        "非常同意"
                    ],
                    "is_required": 1
                }
            ]
        }
    },
    "message": "成功"
}
```

### 返回参数说明：

| 参数名                       | 类型     | 说明                        |
| ------------------------- | ------ | ------------------------- |
| [experiment](#experiment) | object | 实验内容，包括提示信息页面，实验聊天页面，问卷页面 |

### experiment

| 参数名                          | 类型         | 说明                                                 |
| ------------------------------- | ------------ | ---------------------------------------------------- |
| id                              | int          | 实验id,当id为0代表当前无实验                         |
| experiment_topic                | string       | 实验主题，实验中的聊天页的问题部分可以取这个字段的值 |
| [guide_pages](#guide_pages)     | object array | 提示信息页                                           |
| [questionnaire](#questionnaire) | object       | 问卷页内容                                           |

### guide_pages

| 参数名                | 类型   | 说明                                                         |
| --------------------- | ------ | :----------------------------------------------------------- |
| id                    | int    | 提示信息页id                                                 |
| content               | string | 页面文本内容，最好支持html标签                               |
| next_button           | int    | 是否有下一页的按钮，1为是，0为否                             |
| yes_or_no             | int    | 是否有"是"和“否”的按钮，1为是，0为否                         |
| chat                  | int    | 是否是实验的聊天页面，1为是，0为否                           |
| chat_times            | int    | 限制问问题次数，0为不限制，3为限制3次，以此类推              |
| next_page             | int    | 表示下一页按钮跳转到哪个页面(用页面id来定位页面)             |
| yes_page              | int    | 表示"是"按钮跳转到哪个页面(用页面id来定位页面)               |
| no_page               | int    | 表示"否"按钮跳转到哪个页面(用页面id来定位页面)               |
| countdown             | int    | 页面最短阅读时间，0表示不限制，60表示最少阅读60s才能点击按钮（即最短阅读时间结束之前按钮可以采取置灰操作，同时显示相应的倒计时时间） |
| answer_page           | int    | 是否是回答页面，1为是，0为否，如果是聊天页面则在完成聊天后显示答案提交框，供提交实验答案 |
| answer_time_countdown | int    | 回答问题的时间限制，0代表不做限制，60代表限制60s回答时间，建议也显示个倒计时，超过限定时间则不允许提交答案（提交按钮置灰之类的操作） |
| get_chat_from         | int    | 该页面从哪个页面id获取聊天的记录copy显示                     |
| save_chat_info        | int    | 是否保存聊天记录，1位是，0为否                               |
| chat_tips             | string | 聊天框提示标语                                               |
| word_limit_min        | int    | 答案提交页面，答案允许提交的长度最短需多少字，为0代表不限制最小值，100代表最少100字 |
| word_limit_max        | int    | 答案提交页面，答案允许提交的长度最长不超过多少字，为0代表不限制最大值，500代表最多不能超过500字 |

### 特殊说明：

```
当next=1和next_page=0时代表引导页已经结束，next的下一跳就是问卷页
```

### questionnaire

| 参数名                     | 类型           | 说明     |
| ----------------------- | ------------ | ------ |
| id                      | int          | 问卷id   |
| title                   | string       | 问卷标题   |
| [questions](#questions) | object array | 问卷问题列表 |

### questions

| 参数名      | 类型         | 说明                                              |
| ----------- | ------------ | ------------------------------------------------- |
| id          | int          | 实验id                                            |
| content     | string       | 问题内容                                          |
| choice      | string array | 选项内容，选择题和评分题的内容都用在这个选项      |
| type        | int          | 题目类型，1为单选，2为多选，3为填空，4为评分题    |
| score_type  | int          | 评分题的分数量级，1代表1-5分，2为1-7分，3为1-4分  |
| score_text  | string array | 评分题的首行两端的词，如["非常不同意","非常同意"] |
| is_required | int          | 是否必做题，1为是，0为否                          |

## ChatGPT聊天

**接口描述：** 用于调用chatgpt聊天相关操作

**请求URL：** /api/chat

**请求方式：** POST

### 请求参数：

```json
{
	"id": 1,
	"message": "message1"
}
```

| 参数名  | 必选 | 类型   | 说明                                                         |
| ------- | ---- | ------ | ------------------------------------------------------------ |
| id      | 否   | int    | 实验id，用于实验时传id代表是属于实验的聊天记录，而非自由聊天部分的聊天记录，自由聊天部分id传0 |
| message | 是   | string | 问题                                                         |

### 返回示例：

```json
{
    "code": 10000,
    "reply": "gpt回复内容",
    "message": "成功"
}
```

### 返回参数说明：

| 参数名 | 类型   | 说明          |
| ------ | ------ | ------------- |
| reply  | string | gpt回复的内容 |

## ChatGPT聊天事件流

**接口描述：** 用于调用chatgpt聊天相关操作

**请求URL：** ws://127.0.0.1:8089/api/chat/event_stream?message=哪个城市是中国的首都&id=1

**请求方式：** GET

### 请求参数：

| 参数名  | 必选 | 类型   | 说明                                                         |
| ------- | ---- | ------ | ------------------------------------------------------------ |
| id      | 否   | int    | 实验id，用于实验时传id代表是属于实验的聊天记录，而非自由聊天部分的聊天记录，自由聊天部分id传0 |
| message | 是   | string | 问题                                                         |

### 返回示例：

```json
TEXT：北
TEXT：京
TEXT：是
TEXT：中
TEXT：国
TEXT：首
TEXT：都
```

## 实验提交答案

**接口描述：** 用于提交实验的答案

**请求URL：** /api/experiment/answer

**请求方式：** POST

### 请求参数：
```json
{
	"id": 1,
	"answer": "answer1"
}
```

| 参数名  | 必选 | 类型   | 说明 |
| ------- | ---- | ------ | ---- |
| id | 是   | int | 实验id |
| answer | 是   | string | 回答 |

### 返回示例：

```json
{
    "code": 10000,
    "message": "成功"
}
```

### 返回参数说明：

无

## 问卷提交答案

**接口描述：** 用于提交问卷的答案

**请求URL：** /api/questionnaire/answer

**请求方式：** POST

### 请求参数：
```json
{
    "id": 1,
    "answers": [
        {
            "id": 1,
            "answer": [
                1,
                2
            ],
            "text":"",
            "scores": []
        },
        {
            "id": 2,
            "answer": [],
            "text":"填空题答案",
            "scores": []
        },
        {
            "id": 3,
            "answer": [],
            "text":"",
            "scores": [
               7,5,4 
            ]
        }
    ]
}
```
| 参数名  | 必选 | 类型   | 说明 |
| ------- | ---- | ------ | ---- |
| id | 是   | int | 问卷id |
| answers | 是   | int array | 问卷选择题答案 |
| scores | 是 | int array | 评分题答案，如有3个评分栏目，则[7,5,4]代表1-3栏分别为7分，5分，4分 |

### answers

| 参数名  | 必选 | 类型      | 说明                                                        |
| ------- | ---- | --------- | ----------------------------------------------------------- |
| id      | 是   | int       | 问卷问题id                                                  |
| answers | 是   | int array | 问卷问题答案，每个数字代表问题选择的选项,单选多选均使用列表 |
| text    | 是   | string    | 填空题答案，如果是填空题，则使用此字段存储答案              |

### 返回示例：

```json
{
    "code": 10000,
    "message": "成功"
}
```

### 返回参数说明：

无

## 获取页面基础信息

**接口描述：** 获取页面基础信息

**请求URL：** /api/basic_info

**请求方式：** GET

### 请求参数：
无

### 返回示例：

```json
{
    "code": 10000,
    "index_bottom_tips": "<p style=\"text-align: left;\"><span style=\"color: #e03e2d;\">注：</span>本站仅做科研用途，我们郑重承诺您使用中涉及的所有数据和回答只用于学术研究，对您提供的信息我们将匿名处理并严格保密，其他个人、组织或企业并不会得知您个人回答的情况。答案并无对错之分，请您依据自身实际情况填写即可，您的回答对我们的研究将会有非常大的帮助。</p>",
    "login_title": "<p style=\"text-align: center;\">任意输入邮箱和密码即可登录，登录既注册。<br data-v-7d9ffbf3=\"\" />注册限时开放，请牢记您的账号和密码</p>",
    "message": "成功"
}
```

### 返回参数说明：

| 参数名   | 类型     | 说明         |
| ----- | ------ | ---------- |
| login_title | string | 登录页title,希望支持html |
| index_bottom_tips | string | 引导信息页首页底部提示信息，希望支持html |

