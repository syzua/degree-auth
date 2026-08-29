# 基于Hyperledger Fabric的学历认证平台

基于Hyperledger Fabric 2.4联盟链的去中心化学历认证系统，支持学生、高校管理员、用人单位三种角色，实现学历信息上链存证、真伪验证、授权访问、历史追溯等完整功能。

## 功能特点

### 核心功能
- 学历信息上链存证：高校管理员将学历信息写入区块链，数据不可篡改
- 学历真伪验证：用人单位通过证书编号+姓名秒级验证学历真伪
- 授权访问控制：学生可授权用人单位查看详细学历信息
- 修改历史追溯：所有学历信息变更记录链上可追溯审计
- 多角色权限体系：高校（写入）、用人单位（验证）、学生（授权）

### 技术特性
- 联盟链架构：Fabric 2.4 + Raft共识
- 状态数据库：CouchDB（支持富查询）
- 容器化部署：Docker + Docker Compose
- 身份认证：Fabric CA + JWT

## 技术栈

| 层 | 技术 |
|----|------|
| 区块链底层 | Hyperledger Fabric 2.4 |
| 智能合约（链码） | Go 1.19 |
| 后端 | Go + Gin + Fabric-Go-SDK |
| 前端 | Vue.js 2 + Element UI |
| 数据库 | CouchDB 3.2 |
| 部署 | Docker + Docker Compose |
| 身份认证 | Fabric CA + JWT |

## 项目结构

```
degree-auth/
├── chaincode/
│   └── education/
│       ├── education.go          # 学历认证智能合约（链码）
│       └── go.mod
├── backend/
│   ├── main.go                   # 应用入口，路由配置
│   ├── go.mod
│   ├── Dockerfile
│   ├── handlers/
│   │   └── education.go          # RESTful API处理器
│   ├── middleware/
│   │   └── jwt.go                # JWT认证与权限中间件
│   ├── models/
│   │   └── education.go          # 数据模型定义
│   └── services/
│       └── fabric.go             # Fabric SDK服务层
├── frontend/
│   ├── index.html
│   ├── package.json
│   └── src/
│       ├── App.vue               # 主界面（三角色视图）
│       └── main.js
├── network/
│   ├── docker-compose.yml        # Fabric网络编排
│   ├── configtx.yaml             # 通道配置
│   ├── crypto-config.yaml        # 证书生成配置
│   └── scripts/
│       └── deploy.sh             # 一键部署脚本
└── README.md
```

## 智能合约核心函数

| 函数 | 说明 | 调用角色 |
|------|------|---------|
| `AddEducation()` | 添加学历信息上链 | 高校管理员 |
| `UpdateEducation()` | 修改学历信息 | 高校管理员 |
| `QueryEducationByID()` | 按证书编号查询学历 | 所有角色 |
| `VerifyByCertNoAndName()` | 验证学历真伪 | 用人单位 |
| `GetHistoryByID()` | 查询修改历史 | 审计 |
| `AuthorizeViewer()` | 授权查看权限 | 学生 |
| `QueryAllEducation()` | 查询全部学历 | 管理员 |
| `EducationExists()` | 检查学历是否存在 | 内部调用 |

## 系统架构

```
用户浏览器
    ↓ HTTP
Go后端服务 (Gin + JWT认证)
    ↓ gRPC
Fabric SDK (Fabric-Go-SDK)
    ↓
Hyperledger Fabric 网络
    ├── Orderer (Raft共识)
    ├── Peer0 (Org1) + CouchDB0
    ├── Peer1 (Org1) + CouchDB1
    └── Fabric CA (证书颁发)
```

## 快速开始

### 前置要求
- Docker 20.10+
- Docker Compose 2.0+
- Go 1.19+
- Node.js 16+
- Hyperledger Fabric binaries (cryptogen, configtxgen)

### 部署步骤

1. 生成证书和通道配置
```bash
cd network
./scripts/deploy.sh
```

2. 启动后端服务
```bash
cd backend
go mod tidy
go run main.go
```

3. 启动前端
```bash
cd frontend
npm install
npm run dev
```

4. 访问 `http://localhost:8080/frontend`

### 测试账号

| 角色 | 用户名 | 密码 |
|------|--------|------|
| 高校管理员 | admin | admin123 |
| 用人单位 | employer | emp123 |
| 学生 | student | stu123 |

## 使用流程

### 高校管理员
1. 登录系统
2. 填写学历信息（证书编号、姓名、学号、学校、专业、学位、毕业日期）
3. 点击"上链存证"，学历信息写入区块链

### 用人单位
1. 登录系统
2. 输入证书编号和学生姓名
3. 点击"验证"，系统自动比对链上数据，返回验证结果

### 学生
1. 登录系统
2. 查询自己的学历信息
3. 授权用人单位查看详细信息

## License

MIT
