/**
 * Luna SDK - AI Chatbot Example
 *
 * Build a conversational AI assistant using Luna's AI API.
 * This example demonstrates chat completions with context management.
 */

import { LunaClient } from '@eclipse-softworks/luna-sdk';
import * as readline from 'readline';

const client = new LunaClient({
    apiKey: process.env.LUNA_API_KEY!,
});

// ============================================
// Chat Message Types
// ============================================

type Role = 'system' | 'user' | 'assistant';

interface ChatMessage {
    role: Role;
    content: string;
}

// ============================================
// Chatbot Class
// ============================================

class LunaChatbot {
    private conversationHistory: ChatMessage[] = [];
    private model: string;
    private temperature: number;
    private systemPrompt: string;

    constructor(options: {
        model?: string;
        temperature?: number;
        systemPrompt?: string;
    } = {}) {
        this.model = options.model || 'luna-gpt-4';
        this.temperature = options.temperature || 0.7;
        this.systemPrompt =
            options.systemPrompt ||
            'You are a helpful assistant powered by Luna SDK. Be concise and helpful.';

        // Initialize with system prompt
        this.conversationHistory.push({
            role: 'system',
            content: this.systemPrompt,
        });
    }

    async chat(userMessage: string): Promise<string> {
        // Add user message to history
        this.conversationHistory.push({
            role: 'user',
            content: userMessage,
        });

        // Call the AI API
        const response = await client.ai.chatCompletions({
            model: this.model,
            messages: this.conversationHistory,
            temperature: this.temperature,
        });

        // Extract assistant response
        const assistantMessage = response.choices[0]?.message?.content || '';

        // Add to history for context
        this.conversationHistory.push({
            role: 'assistant',
            content: assistantMessage,
        });

        return assistantMessage;
    }

    clearHistory() {
        this.conversationHistory = [
            {
                role: 'system',
                content: this.systemPrompt,
            },
        ];
    }

    getHistory(): ChatMessage[] {
        return [...this.conversationHistory];
    }
}

// ============================================
// Simple Q&A Example
// ============================================

async function simpleQAExample() {
    console.log('Simple Q&A Example\n');

    const response = await client.ai.chatCompletions({
        model: 'luna-gpt-4',
        messages: [
            { role: 'system', content: 'You are a helpful coding assistant.' },
            { role: 'user', content: 'What is a TypeScript interface?' },
        ],
        temperature: 0.5,
    });

    console.log('Question: What is a TypeScript interface?');
    console.log('\nAnswer:', response.choices[0]?.message?.content);
}

// ============================================
// Multi-turn Conversation Example
// ============================================

async function conversationExample() {
    console.log('\nMulti-turn Conversation Example\n');

    const chatbot = new LunaChatbot({
        systemPrompt: 'You are a friendly coding tutor. Explain concepts simply.',
        temperature: 0.6,
    });

    // First message
    console.log('User: Explain what an API is');
    let response = await chatbot.chat('Explain what an API is');
    console.log('Assistant:', response);

    // Follow-up question (chatbot remembers context)
    console.log('\nUser: Can you give me an example?');
    response = await chatbot.chat('Can you give me an example?');
    console.log('Assistant:', response);

    // Another follow-up
    console.log('\nUser: How do REST APIs differ?');
    response = await chatbot.chat('How do REST APIs differ?');
    console.log('Assistant:', response);
}

// ============================================
// Specialized Assistant Example
// ============================================

async function specializedAssistantExample() {
    console.log('\nCode Review Assistant Example\n');

    const codeReviewer = new LunaChatbot({
        model: 'luna-gpt-4',
        temperature: 0.3, // Lower temperature for more focused responses
        systemPrompt: `You are an expert code reviewer. 
            Analyze code for:
            - Bugs and potential issues
            - Performance improvements
            - Best practices
            - Security concerns
            Be specific and constructive.`,
    });

    const codeToReview = `
function fetchUserData(userId) {
    var data = null;
    fetch('/api/users/' + userId)
        .then(response => response.json())
        .then(json => {
            data = json;
        });
    return data;
}
`;

    console.log('Code to review:');
    console.log(codeToReview);

    const review = await codeReviewer.chat(
        `Please review this JavaScript function:\n${codeToReview}`
    );
    console.log('\nCode Review:\n', review);
}

// ============================================
// Interactive Chat (CLI)
// ============================================

async function interactiveChat() {
    console.log('\nInteractive Chat Mode\n');
    console.log('Type your messages below. Type "quit" to exit, "clear" to reset.\n');

    const chatbot = new LunaChatbot({
        systemPrompt: 'You are a helpful AI assistant. Be friendly and informative.',
    });

    const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout,
    });

    const askQuestion = () => {
        rl.question('You: ', async (input) => {
            const trimmedInput = input.trim();

            if (trimmedInput.toLowerCase() === 'quit') {
                console.log('Goodbye!');
                rl.close();
                return;
            }

            if (trimmedInput.toLowerCase() === 'clear') {
                chatbot.clearHistory();
                console.log('Conversation cleared.\n');
                askQuestion();
                return;
            }

            if (!trimmedInput) {
                askQuestion();
                return;
            }

            try {
                const response = await chatbot.chat(trimmedInput);
                console.log(`\nAssistant: ${response}\n`);
            } catch (error) {
                console.error('Error:', error);
            }

            askQuestion();
        });
    };

    askQuestion();
}

// ============================================
// Main
// ============================================

async function main() {
    try {
        await simpleQAExample();
        await conversationExample();
        await specializedAssistantExample();

        // Uncomment to start interactive mode
        // await interactiveChat();
    } catch (error) {
        console.error('Error:', error);
        process.exit(1);
    }
}

main();
