"""
Luna SDK - AI Chatbot Example

Build a conversational AI assistant using Luna's AI API.
This example demonstrates chat completions with context management.
"""

import asyncio
import os
from dataclasses import dataclass, field
from typing import Literal, Callable, Optional

from luna import LunaClient


client = LunaClient(api_key=os.environ["LUNA_API_KEY"])

Role = Literal["system", "user", "assistant"]


@dataclass
class ChatMessage:
    """A chat message."""
    role: Role
    content: str


@dataclass
class LunaChatbot:
    """A chatbot powered by Luna's AI API."""
    
    model: str = "luna-gpt-4"
    temperature: float = 0.7
    system_prompt: str = "You are a helpful assistant powered by Luna SDK. Be concise and helpful."
    conversation_history: list[ChatMessage] = field(default_factory=list)

    def __post_init__(self):
        """Initialize with system prompt."""
        self.conversation_history.append(
            ChatMessage(role="system", content=self.system_prompt)
        )

    async def chat(self, user_message: str) -> str:
        """Send a message and get a response."""
        # Add user message to history
        self.conversation_history.append(
            ChatMessage(role="user", content=user_message)
        )

        # Build messages for API
        messages = [
            {"role": msg.role, "content": msg.content}
            for msg in self.conversation_history
        ]

        # Call the AI API
        response = await client.ai.chat_completions(
            model=self.model,
            messages=messages,
            temperature=self.temperature,
        )

        # Extract assistant response
        assistant_message = response.choices[0].message.content

        # Add to history for context
        self.conversation_history.append(
            ChatMessage(role="assistant", content=assistant_message)
        )

        return assistant_message

    def clear_history(self):
        """Clear conversation history, keeping system prompt."""
        self.conversation_history = [
            ChatMessage(role="system", content=self.system_prompt)
        ]

    def get_history(self) -> list[ChatMessage]:
        """Get a copy of the conversation history."""
        return self.conversation_history.copy()


# ============================================
# Simple Q&A Example
# ============================================

async def simple_qa_example():
    """Demonstrate a simple question and answer."""
    print("Simple Q&A Example\n")

    response = await client.ai.chat_completions(
        model="luna-gpt-4",
        messages=[
            {"role": "system", "content": "You are a helpful coding assistant."},
            {"role": "user", "content": "What is a Python decorator?"},
        ],
        temperature=0.5,
    )

    print("Question: What is a Python decorator?")
    print("\nAnswer:", response.choices[0].message.content)


# ============================================
# Multi-turn Conversation Example
# ============================================

async def conversation_example():
    """Demonstrate a multi-turn conversation."""
    print("\nMulti-turn Conversation Example\n")

    chatbot = LunaChatbot(
        system_prompt="You are a friendly coding tutor. Explain concepts simply.",
        temperature=0.6,
    )

    # First message
    print("User: Explain what an API is")
    response = await chatbot.chat("Explain what an API is")
    print(f"Assistant: {response}")

    # Follow-up question (chatbot remembers context)
    print("\nUser: Can you give me an example?")
    response = await chatbot.chat("Can you give me an example?")
    print(f"Assistant: {response}")

    # Another follow-up
    print("\nUser: How do REST APIs differ?")
    response = await chatbot.chat("How do REST APIs differ?")
    print(f"Assistant: {response}")


# ============================================
# Specialized Assistant Example
# ============================================

async def specialized_assistant_example():
    """Demonstrate a specialized code review assistant."""
    print("\nCode Review Assistant Example\n")

    code_reviewer = LunaChatbot(
        model="luna-gpt-4",
        temperature=0.3,  # Lower temperature for more focused responses
        system_prompt="""You are an expert code reviewer. 
            Analyze code for:
            - Bugs and potential issues
            - Performance improvements
            - Best practices
            - Security concerns
            Be specific and constructive.""",
    )

    code_to_review = '''
def fetch_user_data(user_id):
    data = None
    import requests
    response = requests.get(f'/api/users/{user_id}')
    data = response.json()
    return data
'''

    print("Code to review:")
    print(code_to_review)

    review = await code_reviewer.chat(
        f"Please review this Python function:\n{code_to_review}"
    )
    print("\nCode Review:\n", review)


# ============================================
# Text Analysis Assistant
# ============================================

async def text_analysis_example():
    """Demonstrate text analysis capabilities."""
    print("\nText Analysis Example\n")

    analyzer = LunaChatbot(
        model="luna-gpt-4",
        temperature=0.2,
        system_prompt="""You are a text analysis assistant. 
            When given text, provide:
            1. A brief summary
            2. Key themes/topics
            3. Sentiment analysis
            4. Notable entities mentioned
            Be concise and structured.""",
    )

    sample_text = """
    Luna SDK has revolutionized how developers build applications on the Eclipse 
    platform. With its intuitive API design and comprehensive documentation, 
    teams have reported a 40% reduction in integration time. The SDK supports 
    TypeScript, Python, and Go, making it accessible to developers across 
    different technology stacks. CEO Jane Smith announced that version 2.0 
    will include AI-powered features and enhanced security capabilities.
    """

    print("Analyzing text...")
    analysis = await analyzer.chat(f"Analyze this text:\n{sample_text}")
    print("\nAnalysis:\n", analysis)


# ============================================
# Interactive Chat
# ============================================

async def interactive_chat():
    """Run an interactive chat session."""
    print("\nInteractive Chat Mode\n")
    print("Type your messages below. Type 'quit' to exit, 'clear' to reset.\n")

    chatbot = LunaChatbot(
        system_prompt="You are a helpful AI assistant. Be friendly and informative.",
    )

    while True:
        try:
            user_input = input("You: ").strip()
        except (EOFError, KeyboardInterrupt):
            print("\nGoodbye! 👋")
            break

        if user_input.lower() == "quit":
            print("Goodbye! 👋")
            break

        if user_input.lower() == "clear":
            chatbot.clear_history()
            print("Conversation cleared.\n")
            continue

        if not user_input:
            continue

        try:
            response = await chatbot.chat(user_input)
            print(f"\nAssistant: {response}\n")
        except Exception as e:
            print(f"Error: {e}")


# ============================================
# Streaming Responses (Conceptual)
# ============================================

async def streaming_example():
    """Demonstrate streaming responses (conceptual)."""
    print("\nStreaming Response Example\n")
    
    print("""
# If the SDK supports streaming, you could use:

async for chunk in client.ai.chat_completions_stream(
    model="luna-gpt-4",
    messages=[
        {"role": "user", "content": "Write a poem about coding"}
    ],
):
    print(chunk.delta.content, end="", flush=True)

print()  # Newline after streaming completes
""")


# ============================================
# Main
# ============================================

async def main():
    """Run all examples."""
    try:
        await simple_qa_example()
        await conversation_example()
        await specialized_assistant_example()
        await text_analysis_example()
        await streaming_example()

        # Uncomment to start interactive mode
        # await interactive_chat()

    except Exception as e:
        print(f"Error: {e}")
        raise


if __name__ == "__main__":
    asyncio.run(main())
