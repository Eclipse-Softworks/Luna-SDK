"""
Luna SDK - Python Basic Usage Examples

This file demonstrates common operations with the Luna SDK.
Run with: python basic_usage.py
"""

import asyncio
import os
from luna import LunaClient


# Initialize the client with API key authentication
client = LunaClient(api_key=os.environ["LUNA_API_KEY"])


# ============================================
# Example 1: User Management
# ============================================

async def user_management_example():
    """Demonstrate user CRUD operations."""
    print("=== User Management ===\n")

    # List users with pagination
    user_list = await client.users.list(limit=10)
    print(f"Found {len(user_list.data)} users")

    # Create a new user
    new_user = await client.users.create(
        email="jane.doe@example.com",
        name="Jane Doe",
    )
    print(f"Created user: {new_user.id}")

    # Get user details
    user = await client.users.get(new_user.id)
    print(f"User name: {user.name}, Email: {user.email}")

    # Update the user
    updated_user = await client.users.update(
        new_user.id,
        name="Jane M. Doe",
        avatar_url="https://example.com/avatar.jpg",
    )
    print(f"Updated user name: {updated_user.name}")

    # Delete the user
    await client.users.delete(new_user.id)
    print("User deleted")


# ============================================
# Example 2: Project Management
# ============================================

async def project_management_example():
    """Demonstrate project CRUD operations."""
    print("\n=== Project Management ===\n")

    # Create a project
    project = await client.projects.create(
        name="My Awesome App",
        description="A revolutionary application built with Luna SDK",
    )
    print(f"Created project: {project.id}")

    # List all projects
    projects = await client.projects.list(limit=20)
    print(f"Total projects: {len(projects.data)}")

    # Get project details
    project_details = await client.projects.get(project.id)
    print(f"Project: {project_details.name}")
    print(f"Owner: {project_details.owner_id}")
    print(f"Created: {project_details.created_at}")

    # Update project
    updated = await client.projects.update(
        project.id,
        description="Updated description with new features",
    )
    print(f"Updated project: {updated.description}")

    # Clean up
    await client.projects.delete(project.id)
    print("Project deleted")


# ============================================
# Example 3: Paginating Through Results
# ============================================

async def pagination_example():
    """Demonstrate automatic pagination with iterators."""
    print("\n=== Pagination Example ===\n")

    # Using the async iterator for automatic pagination
    count = 0
    async for user in client.users.iterate(limit=10):
        print(f"User: {user.name} ({user.email})")
        count += 1
        if count >= 50:  # Limit for demo
            break
    
    print(f"Iterated through {count} users")


# ============================================
# Example 4: Error Handling
# ============================================

async def error_handling_example():
    """Demonstrate proper error handling."""
    print("\n=== Error Handling ===\n")

    from luna.errors import NotFoundError, ValidationError, RateLimitError

    try:
        # Try to get a non-existent user
        await client.users.get("usr_nonexistent123")
    except NotFoundError as e:
        print(f"Not found error: {e.message}")
        print(f"Error code: {e.code}")
    except ValidationError as e:
        print(f"Validation error: {e.message}")
    except RateLimitError as e:
        print(f"Rate limited! Retry after: {e.retry_after} seconds")
    except Exception as e:
        print(f"Unexpected error: {e}")


# ============================================
# Example 5: Context Manager Usage
# ============================================

async def context_manager_example():
    """Demonstrate using the client as a context manager."""
    print("\n=== Context Manager Example ===\n")

    async with LunaClient(api_key=os.environ["LUNA_API_KEY"]) as luna:
        users = await luna.users.list(limit=5)
        print(f"Found {len(users.data)} users")
        
        for user in users.data:
            print(f"  - {user.name}")
    
    print("Client automatically closed")


# ============================================
# Run all examples
# ============================================

async def main():
    """Run all examples."""
    try:
        await user_management_example()
        await project_management_example()
        await pagination_example()
        await error_handling_example()
        await context_manager_example()
    except Exception as e:
        print(f"Example failed: {e}")
        raise


if __name__ == "__main__":
    asyncio.run(main())
