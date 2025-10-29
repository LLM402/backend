# LLM 402

🍥 Next-Generation Large Model Gateway and AI Asset Management System

## ✨ Key Features

LLM402 offers a wide range of features:

1. 💵 The LLM402 supports the X402 payment protocol
2. ⚖️ Support for weighted random channel selection
3. 📈 Data dashboard (console)
4. 🔒 Token grouping and model restrictions
5. 🤖 Support for more authorization login methods (LinuxDO, Telegram, OIDC)
6. 🔄 Support for Rerank models (Cohere and Jina)
7. ⚡ Support for OpenAI Realtime API (including Azure channels)
8. ⚡ Support for **OpenAI Responses** format
9. ⚡ Support for **Claude Messages** format
10. ⚡ Support for **Google Gemini** format
11. 🧠 Support for setting reasoning effort through model name suffixes:
    1. OpenAI o-series models
        - Add `-high` suffix for high reasoning effort (e.g.: `o3-mini-high`)
        - Add `-medium` suffix for medium reasoning effort (e.g.: `o3-mini-medium`)
        - Add `-low` suffix for low reasoning effort (e.g.: `o3-mini-low`)
    2. Claude thinking models
        - Add `-thinking` suffix to enable thinking mode (e.g.: `claude-3-7-sonnet-20250219-thinking`)
12. 🔄 Thinking-to-content functionality
13. 🔄 Model rate limiting for users
14. 🔄 Request format conversion functionality, supporting the following three format conversions:
    1. OpenAI Chat Completions => Claude Messages
    2. Claude Messages => OpenAI Chat Completions (can be used for Claude Code to call third-party models)
    3. OpenAI Chat Completions => Gemini Chat
25. 💰 Cache billing support, which allows billing at a set ratio when cache is hit:
    1. Set the `Prompt Cache Ratio` option in `System Settings-Operation Settings`
    2. Set `Prompt Cache Ratio` in the channel, range 0-1, e.g., setting to 0.5 means billing at 50% when cache is hit
    3. Supported channels:
        - [x] OpenAI
        - [x] Azure
        - [x] DeepSeek
        - [x] Claude

## Deployment

### Multi-machine Deployment Considerations
- Environment variable `SESSION_SECRET` must be set, otherwise login status will be inconsistent across multiple machines
- If sharing Redis, `CRYPTO_SECRET` must be set, otherwise Redis content cannot be accessed across multiple machines

### Deployment Requirements
- Local database (default): SQLite (Docker deployment must mount the `/data` directory)
- Remote database: MySQL version >= 5.7.8, PgSQL version >= 9.6

```shell
# Download the project
git clone git@github.com:LLM402/backend.git
cd backend
# Edit docker-compose.yml as needed
# Start
brew install go
go download
go run main.go
go build -o llm402
```