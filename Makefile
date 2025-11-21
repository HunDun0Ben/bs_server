# ==============================================================================
# Variables
# ==============================================================================

# 输出目录
# 用于存放主应用
APP_BIN_DIR     := bin
# 用于存放脚本工具
BIN_DIR         := scripts/bin

# 主应用
MAIN_APP_SRC    := app/main.go
MAIN_APP_TARGET := $(BIN_DIR)/bs_server

# JWT 工具
JWT_TOOL_SRC    := app/scripts/jwtscr/generate_jwt_tokens.go
JWT_TOOL_TARGET := $(BIN_DIR)/generate_jwt_tokens

# Swagger
SWAGGER_MAIN_FILE  := app/main.go
SWAGGER_OUTPUT_DIR := app/docs/swagger

# ==============================================================================
# Main Targets
# ==============================================================================

# .PHONY 告诉 make, 这些目标不是真实的文件名
.PHONY: all build tools swagger clean help format

# 默认目标：构建所有内容
all: build tools swagger ## Build main app, tools, and generate docs

# 构建主应用
build: $(MAIN_APP_TARGET) ## Build the main application

# 构建所有的 Go 脚本工具
tools: $(JWT_TOOL_TARGET) ## Build all go scripts tools

# 生成 Swagger/OpenAPI 文档
swagger: ## Generate Swagger/OpenAPI documentation
	@echo "📜 Generating Swagger docs..."
	swag init -g $(SWAGGER_MAIN_FILE) --output $(SWAGGER_OUTPUT_DIR)

format: ## Format Go files using gci
	@echo "🎨 Formatting Go files..."
	gci write --section standard --section default --section "prefix(github.com/HunDun0Ben/bs_server)" --section alias --section blank --section dot .

# 清理所有生成的文件
clean: ## Clean up all generated files
	@echo "🧹 Cleaning up..."
	rm -rf $(BIN_DIR) $(APP_BIN_DIR)
	rm -f $(SWAGGER_OUTPUT_DIR)/swagger.* $(SWAGGER_OUTPUT_DIR)/docs.go

# 显示帮助信息
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) \
		| sed -e 's/:.*## /|/' \
		| column -t -s '|'

# ==============================================================================
# Build Rules
# ==============================================================================

# 构建主应用的规则
$(MAIN_APP_TARGET): $(MAIN_APP_SRC)
	@mkdir -p $(APP_BIN_DIR)
	@echo "🚀 Building main application..."
	go build -o $(MAIN_APP_TARGET) $(MAIN_APP_SRC)

# 构建 JWT 工具的规则
$(JWT_TOOL_TARGET): $(JWT_TOOL_SRC)
	@mkdir -p $(BIN_DIR)
	@echo "🔨 Building JWT tool..."
	go build -o $(JWT_TOOL_TARGET) $(JWT_TOOL_SRC)
