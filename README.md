# Metaphysics - AI 智能八字分析系统

## 简介

Metaphysics 是一个基于人工智能的中国传统八字命理分析系统，结合现代 AI 技术与传统命理学，为用户提供专业、客观的命理分析服务。系统采用前后端分离架构，后端基于 Go 语言开发，前端使用 React 框架，实现高效、准确的八字分析与现代化的用户界面。

## 项目结构

```
metaphysics/
├── backend/                # 后端项目
│   ├── cmd/                # 命令行入口
│   ├── configs/            # 配置文件
│   ├── docs/               # 文档
│   ├── internal/           # 内部包
│   │   ├── ai/             # AI服务集成
│   │   ├── model/          # 数据模型
│   │   └── utils/          # 工具函数
│   ├── pkg/                # 公共包
│   └── main.go             # 主程序入口
│
└── frontend/               # 前端项目 (React)
    ├── public/             # 静态资源
    ├── src/                # 源代码
    │   ├── components/     # UI组件
    │   ├── pages/          # 页面
    │   ├── services/       # API服务
    │   ├── store/          # Redux状态管理
    │   └── App.tsx         # 应用入口
    └── package.json        # 依赖配置
```

## 核心功能

- **精准八字排盘**：基于 lunar-go 库，支持公历/农历转换，精确计算四柱八字
- **AI 命理分析**：接入 Deepseek、Ollama 等大语言模型，提供专业的命理解读
- **多维度分析**：包括性格特征、事业发展、财运结构、健康体质、婚姻感情、未来趋势等
- **流式响应**：支持 AI 分析结果的流式输出，提升用户体验
- **历史记录**：保存用户的分析记录，支持随时查询

## 技术栈

### 后端

- **Go 1.23+**：主要开发语言
- **Gin**：Web 框架
- **GORM**：ORM 框架，支持多种数据库
- **Redis**：go-redis/v9 客户端
- **Logrus**：结构化日志与日志轮转
- **Viper**：配置管理与热加载
- **LangChain**：AI 推理与分析
- **Lunar-Go**：中国传统历法计算
- **中间件**：requestid、CORS、Gzip、Secure、Recovery 等

### 前端

- **React 18**：用户界面库
- **TypeScript**：类型安全的 JavaScript 超集
- **Redux Toolkit**：状态管理
- **React Router**：路由管理
- **Material UI**：UI 组件库
- **Axios**：HTTP 客户端

## 快速开始

### 后端

1. 进入后端目录 `cd backend`
2. 配置 `configs/config.yaml`，设置数据库和 AI 服务参数
3. 启动服务 `go run main.go`

### 前端

1. 进入前端目录 `cd frontend`
2. 安装依赖 `npm install`
3. 启动开发服务器 `npm start`
4. 构建生产版本 `npm run build`

## API 接口

### 分析八字

**请求**：`POST /api/v1/bazi/analyze`

```json
{
  "name": "张三",
  "gender": "男",
  "birth_time": "2000-01-01T12:00:00Z",
  "calendar": "solar"
}
```

**响应**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "张三",
    "gender": "男",
    "calendar": "solar",

    "year_pillar": "庚辰",
    "month_pillar": "己丑",
    "day_pillar": "丙子",
    "hour_pillar": "戊午",

    "year_gan": "庚",
    "year_zhi": "辰",
    "month_gan": "己",
    "month_zhi": "丑",
    "day_gan": "丙",
    "day_zhi": "子",
    "hour_gan": "戊",
    "hour_zhi": "午",

    "yin_yang": "阳",
    "wu_xing": "火",

    "analysis": "AI分析结果..."
  }
}
```

### 流式分析

**请求**：`POST /api/v1/bazi/analyze/stream`

支持 Server-Sent Events (SSE) 格式的流式响应，实时返回 AI 分析结果。

### 查询记录

**请求**：`GET /api/v1/bazi/record/{request_id}`

## AI 模型支持

系统支持多种 AI 模型，可在配置文件中灵活切换：

- **Deepseek**：专业大模型，提供更准确的命理解读
- **Ollama**：支持本地部署的开源模型，降低延迟和成本

## 命理分析维度

AI 分析结果涵盖六大核心维度：

1. **命格结构与性格特征**：分析五行分布、日主旺衰、性格特点等
2. **事业发展方向与职业类型**：分析适合的职业方向、事业发展关键期等
3. **财运结构与财富管理**：分析财运来源、理财倾向、守财能力等
4. **健康体质与调养建议**：分析体质特点、潜在健康风险、养生方向等
5. **情感婚姻与人际关系**：分析婚恋趋势、适合的伴侣特征等
6. **大运流年与未来趋势**：分析未来五年运势变化、关键时间点等

## 依赖库

### 后端

- [github.com/6tail/lunar-go](https://github.com/6tail/lunar-go)：日历、公历(阳历)、农历(阴历、老黄历)、道历、佛历，支持节假日、星座、儒略日、干支、生肖、节气、节日等。
- [github.com/tmc/langchaingo](https://github.com/tmc/langchaingo)：Go 语言的 LangChain 实现，用于 AI 推理分析。

### 前端

- [React](https://reactjs.org/)：用户界面库
- [Material UI](https://mui.com/)：React UI 组件库
- [Redux Toolkit](https://redux-toolkit.js.org/)：状态管理工具
- [Axios](https://axios-http.com/)：基于 Promise 的 HTTP 客户端
