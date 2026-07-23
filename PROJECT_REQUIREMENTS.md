# 商品库存管理系统需求说明

## 1. 项目定位

本项目第一版定位为一个轻量级商品库存管理系统，主要用于个人或小团队管理同一套库存数据。

系统不开放用户注册，由后端初始化脚本或超级管理员创建账号。用户登录后根据角色权限访问不同功能。

第一版优先实现：

- 用户登录
- 用户管理
- 角色权限管理
- 商品库存管理
- 基于角色的权限控制

暂不引入公司、空间、多租户、仓库、入库出库流水等复杂业务。

## 2. 技术方向

当前计划技术栈：

- 前端：React
- 后端：Go + Gin
- 数据库：PostgreSQL
- ORM：GORM
- 进程管理：Supervisor，后续也可以容器化部署

当前项目结构计划：

- frontend：前端项目，使用 Vite 创建 React 项目
- backend：后端项目，使用 Go + Gin
- init_db：数据库启动与初始化相关文件
- Dockerfile.frontend：前端镜像构建文件，构建 React 静态资源并用 Nginx 提供服务
- Dockerfile.backend：后端镜像构建文件，构建 Go API 服务和初始化脚本
- docker-compose.frontend.yml：前端服务 Compose
- docker-compose.backend.yml：后端服务 Compose

数据库 Compose 文件放在 init_db 目录：

- init_db/docker-compose.yml
- 数据库名：inventory_db
- 数据库用户：inventory_user

后续部署方式按分布式部署考虑：

- 数据库有独立 Compose
- 后端可以有独立 Compose
- 前端可以有独立 Compose
- 前端和后端分别构建为独立镜像
- 数据库、后端、前端通过同一个 Docker 网络 inventory_network 通信

技术栈和部署方式后续可以根据项目发展调整。当前文档先聚焦第一版业务范围、模型和功能边界。

## 3. 账号与登录

### 3.1 登录方式

第一版使用用户名和密码登录。

不使用：

- 邮箱登录
- 邮箱验证码
- 手机号登录
- 第三方登录
- 找回密码

### 3.2 用户来源

系统不开放注册。

用户来源有两种：

- 后端初始化脚本创建默认账号
- 超级管理员在用户管理页面创建账号

### 3.3 初始化账号

后端提供初始化脚本，用于创建：

- 默认权限
- 默认角色
- 默认用户

示例账号可以是：

| 用户名 | 角色 | 说明 |
| --- | --- | --- |
| root | super_admin | 超级管理员 |
| manager | stock_manager | 库存管理员 |
| viewer | stock_viewer | 库存查看员 |

密码不建议在代码中写死，开发阶段可以使用默认密码，正式部署时应通过环境变量或初始化参数设置。

## 4. 库存数据范围

第一版采用单空间库存模型。

所有登录用户访问的是同一套库存数据，不按用户隔离库存。

也就是说：

- 超级管理员可以看到全部库存
- 库存管理员可以看到全部库存，并可维护库存
- 库存查看员可以看到全部库存，但不能修改

后续如果需要服务多个公司或团队，可以扩展为公司/空间模型：

- 公司/空间
- 用户属于某个公司/空间
- 库存属于某个公司/空间
- 角色在公司/空间内生效

第一版暂不实现该扩展。

## 5. 角色设计

第一版内置三个角色。

| 角色标识 | 角色名称 | 说明 |
| --- | --- | --- |
| super_admin | 超级管理员 | 拥有系统全部权限 |
| stock_manager | 库存管理员 | 可以查看、新增、编辑、删除库存 |
| stock_viewer | 库存查看员 | 只能查看库存 |

## 6. 权限设计

权限采用字符串标识，用户通过角色获得权限。

### 6.1 用户权限

| 权限标识 | 说明 |
| --- | --- |
| user:view | 查看用户列表 |
| user:create | 创建用户 |
| user:update | 编辑用户 |
| user:disable | 启用或禁用用户 |

### 6.2 角色权限

| 权限标识 | 说明 |
| --- | --- |
| role:view | 查看角色列表 |
| role:create | 创建角色 |
| role:update | 编辑角色 |
| role:delete | 删除角色 |

### 6.3 库存权限

| 权限标识 | 说明 |
| --- | --- |
| stock:view | 查看库存 |
| stock:create | 新增商品 |
| stock:update | 编辑商品 |
| stock:delete | 删除商品 |

### 6.4 默认角色权限

| 角色 | 权限 |
| --- | --- |
| super_admin | 全部权限 |
| stock_manager | stock:view, stock:create, stock:update, stock:delete |
| stock_viewer | stock:view |

## 7. 功能模块

### 7.1 登录页

用户输入：

- 用户名
- 密码

登录成功后进入库存管理页面。

登录失败时提示用户名或密码错误。

禁用用户不能登录。

### 7.2 库存管理

所有登录用户可访问库存管理页面。

库存查看员只能查看列表和详情，不能新增、编辑、删除。

库存管理员和超级管理员可以进行库存维护。

第一版商品字段：

| 字段 | 说明 | 是否必填 |
| --- | --- | --- |
| name | 商品名称 | 是 |
| sku | 商品编码/SKU | 是 |
| category | 分类 | 否 |
| quantity | 当前库存数量 | 是 |
| remark | 备注 | 否 |
| created_at | 创建时间 | 系统生成 |
| updated_at | 更新时间 | 系统生成 |

库存管理第一版功能：

- 查看商品列表
- 按商品名称、SKU 搜索
- 新增商品
- 编辑商品
- 删除商品

第一版库存数量允许直接编辑，但每次库存数量发生变化时，后端需要写入库存变动记录。

第一版暂不区分入库、出库、盘点等业务流程，只记录库存数量从多少变成多少、由谁操作、什么时候操作。

### 7.3 用户管理

只有超级管理员可以访问用户管理页面。

第一版功能：

- 查看用户列表
- 新建用户
- 编辑用户信息
- 启用用户
- 禁用用户
- 给用户分配角色

用户字段建议：

| 字段 | 说明 |
| --- | --- |
| username | 用户名 |
| password_hash | 密码哈希 |
| display_name | 显示名称 |
| status | 用户状态，enabled/disabled |
| role_id | 角色 ID |
| created_at | 创建时间 |
| updated_at | 更新时间 |

第一版不做用户删除，避免误删历史数据。

### 7.4 角色权限管理

只有超级管理员可以访问角色权限管理页面。

第一版功能：

- 查看角色列表
- 新建角色
- 编辑角色名称
- 配置角色权限
- 删除自定义角色

内置角色建议不允许删除：

- super_admin
- stock_manager
- stock_viewer

## 8. 数据模型

### 8.1 用户表 users

| 字段 | 说明 |
| --- | --- |
| id | 用户 ID |
| username | 用户名，唯一 |
| password_hash | 密码哈希 |
| display_name | 显示名称 |
| status | 用户状态，enabled/disabled |
| role_id | 角色 ID |
| created_at | 创建时间 |
| updated_at | 更新时间 |

说明：

- 用户名用于登录。
- 密码只保存哈希，不保存明文。
- 禁用用户不能登录。
- 第一版一个用户只绑定一个角色。

### 8.2 角色表 roles

| 字段 | 说明 |
| --- | --- |
| id | 角色 ID |
| code | 角色标识，唯一，例如 super_admin |
| name | 角色名称，例如 超级管理员 |
| description | 角色说明 |
| is_system | 是否系统内置角色 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

说明：

- 系统内置角色不允许删除。
- 角色标识用于后端逻辑判断和初始化数据。
- 角色名称用于前端展示。

### 8.3 权限表 permissions

| 字段 | 说明 |
| --- | --- |
| id | 权限 ID |
| code | 权限标识，唯一，例如 stock:view |
| name | 权限名称，例如 查看库存 |
| module | 所属模块，例如 stock/user/role |
| description | 权限说明 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

说明：

- 权限由系统初始化脚本创建。
- 第一版权限可以不在页面新增或删除，只在角色权限管理中分配。

### 8.4 角色权限关联表 role_permissions

| 字段 | 说明 |
| --- | --- |
| id | 关联 ID |
| role_id | 角色 ID |
| permission_id | 权限 ID |
| created_at | 创建时间 |

说明：

- 一个角色可以拥有多个权限。
- 一个权限可以分配给多个角色。
- role_id 和 permission_id 建议建立唯一约束，避免重复分配。

### 8.5 商品库存表 products

| 字段 | 说明 |
| --- | --- |
| id | 商品 ID |
| name | 商品名称 |
| sku | 商品编码/SKU |
| category | 分类 |
| quantity | 当前库存数量 |
| remark | 备注 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

说明：

- sku 建议唯一。
- quantity 为当前库存数量。
- 新增商品时，如果初始库存数量不为 0，也需要写入库存变动记录。

### 8.6 库存变动记录表 stock_movements

| 字段 | 说明 |
| --- | --- |
| id | 记录 ID |
| product_id | 商品 ID |
| before_quantity | 变动前库存数量 |
| after_quantity | 变动后库存数量 |
| change_quantity | 变动数量，after_quantity - before_quantity |
| operator_id | 操作用户 ID |
| remark | 备注 |
| created_at | 创建时间 |

说明：

- 每次商品库存数量发生变化时，都必须新增一条库存变动记录。
- 只修改商品名称、SKU、分类、备注时，如果库存数量没有变化，不需要写入库存变动记录。
- 删除商品时第一版可以不写库存变动记录，后续如果需要可以扩展操作类型字段。

## 9. 页面与菜单

第一版页面结构：

- 登录页
- 主布局
  - 库存管理
  - 用户管理
  - 角色权限管理

菜单根据权限显示：

| 菜单 | 显示条件 |
| --- | --- |
| 库存管理 | 拥有 stock:view |
| 用户管理 | 拥有 user:view |
| 角色权限管理 | 拥有 role:view |

按钮根据权限显示：

| 按钮 | 显示条件 |
| --- | --- |
| 新增商品 | 拥有 stock:create |
| 编辑商品 | 拥有 stock:update |
| 删除商品 | 拥有 stock:delete |
| 新建用户 | 拥有 user:create |
| 编辑用户 | 拥有 user:update |
| 启用/禁用用户 | 拥有 user:disable |
| 新建角色 | 拥有 role:create |
| 编辑角色 | 拥有 role:update |
| 删除角色 | 拥有 role:delete |

后端接口必须继续做权限校验，不能只依赖前端隐藏菜单或按钮。

## 10. 开发与部署约定

### 10.1 前端项目

前端使用 Vite 创建 React 项目，目录为 frontend。

前端主要负责：

- 登录页
- 主布局
- 库存管理页面
- 用户管理页面
- 角色权限管理页面
- 根据权限显示菜单和按钮

### 10.2 后端项目

后端使用 Go + Gin，目录为 backend。

建议目录结构：

```text
backend/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── seed/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── middleware/
│   ├── models/
│   ├── routes/
│   ├── handlers/
│   ├── services/
│   └── repositories/
├── migrations/
├── go.mod
└── go.sum
```

后端主要负责：

- 登录认证
- JWT 或类似 token 签发与校验
- 用户 CRUD
- 角色权限 CRUD
- 库存 CRUD
- 库存变动记录写入
- 接口级权限校验

### 10.3 数据库

数据库使用 PostgreSQL。

数据库启动文件：

```text
init_db/docker-compose.yml
```

开发阶段可以只启动数据库容器，前端和后端在本机运行。

启动命令：

```bash
cd init_db
docker compose up -d
```

### 10.4 Docker

项目后续使用一个统一 Dockerfile。

Dockerfile 计划采用多阶段构建：

- frontend build stage：构建前端静态资源
- backend build stage：编译 Go 后端二进制
- frontend runtime stage：运行或托管前端产物
- backend runtime stage：运行 Go 后端服务

前端、后端、数据库可以分别通过不同 Compose 文件部署，便于分布式部署和单独升级。

当前 Compose 文件：

```text
init_db/docker-compose.yml
docker-compose.backend.yml
docker-compose.frontend.yml
```

前端容器监听 8080 端口，并通过 Nginx 将 /api 请求代理到后端 3000 端口。

后端容器监听 3000 端口，通过环境变量连接 PostgreSQL。

## 11. 第一版不做的内容

以下内容暂不进入第一版：

- 用户自行注册
- 邮箱验证码
- 找回密码
- 公司/空间/多租户
- 多仓库
- 入库/出库业务流程
- 库存盘点
- 商品图片
- 条码扫码
- Excel 导入导出
- 供应商管理
- 客户管理
- 操作日志详情页
- 审批流程

## 12. 后续可扩展方向

第一版完成后，可以逐步扩展：

- 公司/空间模型
- 多仓库管理
- 入库/出库记录
- 更完整的库存流水
- 库存预警
- Excel 导入导出
- 操作日志
- 数据看板
- 商品分类管理
- 供应商管理

## 13. 第一版目标总结

第一版目标是打通一个简单、清晰、可用的库存管理闭环：

1. 后端初始化用户、角色、权限
2. 用户通过用户名和密码登录
3. 所有用户共享一套库存数据
4. 不同角色拥有不同操作权限
5. 超级管理员管理用户和角色
6. 库存管理员维护商品库存
7. 库存查看员只能查看库存
8. 每次库存数量变化都记录库存变动
