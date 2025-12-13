# Luna SDK Examples

This directory contains practical code examples demonstrating how to build applications using the Luna SDK. Examples are provided for all three supported languages: **TypeScript**, **Python**, and **Go**.

## Structure

```
examples/
├── typescript/
│   ├── basic-usage.ts          # Core SDK operations
│   ├── resmate-finder.ts       # Student housing search app
│   ├── ai-chatbot.ts           # Conversational AI assistant
│   ├── file-storage.ts         # File management system
│   └── workflow-automation.ts  # Workflow automation
├── python/
│   ├── basic_usage.py          # Core SDK operations
│   ├── resmate_finder.py       # Student housing search app
│   ├── ai_chatbot.py           # Conversational AI assistant
│   └── team_collaboration_app.py # Full-stack team app
└── go/
    ├── basic_usage.go          # Core SDK operations
    ├── resmate_finder.go       # Student housing search app
    ├── ai_chatbot.go           # Conversational AI assistant
    └── project_manager.go      # Project management CLI
```

## Getting Started

### Prerequisites

1. Get your Luna API key from the [Eclipse Developer Portal](https://developer.eclipse.dev)
2. Set the environment variable:

```bash
# Linux/macOS
export LUNA_API_KEY="lk_prod_your_api_key_here"

# Windows (PowerShell)
$env:LUNA_API_KEY = "lk_prod_your_api_key_here"
```

### Running TypeScript Examples

```bash
cd examples/typescript

# Install dependencies (if needed)
npm install @eclipse-softworks/luna-sdk

# Run examples
npx tsx basic-usage.ts
npx tsx resmate-finder.ts
npx tsx ai-chatbot.ts
npx tsx file-storage.ts
npx tsx workflow-automation.ts
```

### Running Python Examples

```bash
cd examples/python

# Install the SDK
pip install luna-sdk

# Run examples
python basic_usage.py
python resmate_finder.py
python ai_chatbot.py
python team_collaboration_app.py
```

### Running Go Examples

```bash
cd examples/go

# Initialize module and get dependencies
go mod init examples
go get github.com/eclipse-softworks/luna-sdk-go

# Run examples
go run basic_usage.go
go run resmate_finder.go
go run ai_chatbot.go
go run project_manager.go
```

## Example Descriptions

### Basic Usage
Demonstrates core SDK functionality:
- Client initialization with API key
- User CRUD operations (Create, Read, Update, Delete)
- Project management
- Automatic pagination with iterators
- Error handling patterns

### ResMate Finder
A student accommodation search application:
- Listing available campuses
- Searching residences with filters (price, NSFAS, gender, rating)
- Detailed residence information
- Comparing multiple residences
- Automatic pagination through results

### AI Chatbot
Build conversational AI assistants:
- Simple question and answer
- Multi-turn conversations with context
- Specialized assistants (code reviewer, text analyzer)
- Interactive chat mode
- Temperature and model configuration

### File Storage (TypeScript)
File management system:
- Listing storage buckets
- File upload and download
- Generating download URLs
- Storage usage reports
- Batch file operations

### Workflow Automation (TypeScript)
Automate business processes:
- Listing available workflows
- Triggering workflows with parameters
- Event-driven workflows
- Webhook integrations
- Workflow orchestration patterns

### Team Collaboration App (Python)
Full-stack application combining multiple SDK features:
- Team/project creation
- Member management
- AI-powered project summaries
- Task suggestions
- Storage and workflow integration

### Project Manager (Go)
CLI project management tool:
- Team creation and management
- AI-generated summaries and tasks
- Storage information display
- Workflow listing
- Interactive dashboard

## Key Concepts

### Authentication

```typescript
// API Key authentication (recommended for server-side)
const client = new LunaClient({
    apiKey: process.env.LUNA_API_KEY,
});

// Token authentication (for client-side with OAuth)
const client = new LunaClient({
    accessToken: session.accessToken,
    refreshToken: session.refreshToken,
    onTokenRefresh: async (tokens) => {
        await saveTokens(tokens);
    },
});
```

### Pagination

All list operations support automatic pagination:

```typescript
// Manual pagination
const page1 = await client.users.list({ limit: 10 });
const page2 = await client.users.list({ limit: 10, cursor: page1.next_cursor });

// Automatic iteration
for await (const user of client.users.iterate({ limit: 10 })) {
    console.log(user.name);
}
```

### Error Handling

```typescript
import { NotFoundError, ValidationError, RateLimitError } from '@eclipse-softworks/luna-sdk';

try {
    const user = await client.users.get('usr_invalid');
} catch (error) {
    if (error instanceof NotFoundError) {
        console.log('User not found');
    } else if (error instanceof RateLimitError) {
        console.log(`Rate limited. Retry after ${error.retryAfter}s`);
    }
}
```

## Resources

- [Luna SDK Documentation](https://docs.eclipse.dev/luna-sdk)
- [API Reference](https://api.eclipse.dev/docs)
- [Developer Portal](https://developer.eclipse.dev)
- [GitHub Repository](https://github.com/eclipse-softworks/luna-sdk)

## License

These examples are provided under the MIT License. See the main repository for full license details.
