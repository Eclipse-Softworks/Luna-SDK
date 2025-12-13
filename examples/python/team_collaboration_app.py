"""
Luna SDK - Full-Stack Application Example

Build a complete application using multiple Luna SDK features.
This example demonstrates a team collaboration tool.
"""

import asyncio
import os
from dataclasses import dataclass
from datetime import datetime
from typing import Optional

from luna import LunaClient
from luna.errors import NotFoundError, ValidationError


client = LunaClient(api_key=os.environ["LUNA_API_KEY"])


# ============================================
# Team Collaboration Application
# ============================================

@dataclass
class TeamMember:
    """Represents a team member."""
    user_id: str
    name: str
    email: str
    role: str


@dataclass
class Team:
    """Represents a team with members and projects."""
    project_id: str
    name: str
    members: list[TeamMember]


class TeamCollaborationApp:
    """A team collaboration application built with Luna SDK."""

    def __init__(self):
        self.current_team: Optional[Team] = None

    # ============================================
    # Team Management
    # ============================================

    async def create_team(self, name: str, description: str) -> Team:
        """Create a new team (project)."""
        print(f"\nCreating team: {name}")

        # Create project as the team container
        project = await client.projects.create(
            name=name,
            description=description,
        )

        team = Team(
            project_id=project.id,
            name=project.name,
            members=[],
        )

        print(f"Team created with ID: {project.id}")
        return team

    async def add_member(self, team: Team, email: str, name: str, role: str) -> TeamMember:
        """Add a member to the team."""
        print(f"\nAdding member: {name}")

        # Create or find user
        try:
            user = await client.users.create(
                email=email,
                name=name,
            )
            print(f"   Created new user: {user.id}")
        except ValidationError:
            # User might already exist
            users = await client.users.list(limit=100)
            user = next((u for u in users.data if u.email == email), None)
            if not user:
                raise ValueError(f"Could not find or create user with email: {email}")
            print(f"   Found existing user: {user.id}")

        member = TeamMember(
            user_id=user.id,
            name=user.name,
            email=user.email,
            role=role,
        )

        team.members.append(member)
        print(f"{name} added as {role}")
        return member

    async def list_members(self, team: Team):
        """List all team members."""
        print(f"\nTeam: {team.name}")
        print("-" * 40)

        if not team.members:
            print("   No members yet")
            return

        for member in team.members:
            print(f"   • {member.name} ({member.email})")
            print(f"     Role: {member.role}")
            print(f"     ID: {member.user_id}")

    # ============================================
    # AI-Powered Features
    # ============================================

    async def generate_project_summary(self, team: Team) -> str:
        """Use AI to generate a project summary."""
        print(f"\nGenerating project summary...")

        # Get project details
        project = await client.projects.get(team.project_id)

        # Generate summary using AI
        response = await client.ai.chat_completions(
            model="luna-gpt-4",
            messages=[
                {
                    "role": "system",
                    "content": "You are a project manager assistant. Generate concise, professional summaries.",
                },
                {
                    "role": "user",
                    "content": f"""Generate a brief project status summary for:
                    
Project: {project.name}
Description: {project.description}
Team Size: {len(team.members)} members
Created: {project.created_at}

Include:
1. Project overview
2. Team composition
3. Suggested next steps""",
                },
            ],
            temperature=0.5,
        )

        summary = response.choices[0].message.content
        print("\nProject Summary:\n")
        print(summary)
        return summary

    async def suggest_tasks(self, team: Team, context: str) -> str:
        """Use AI to suggest tasks for the team."""
        print(f"\nGenerating task suggestions...")

        response = await client.ai.chat_completions(
            model="luna-gpt-4",
            messages=[
                {
                    "role": "system",
                    "content": "You are a project planning assistant. Suggest actionable tasks.",
                },
                {
                    "role": "user",
                    "content": f"""Based on this context, suggest 5 actionable tasks:
                    
Team: {team.name}
Team Members: {', '.join(m.name for m in team.members)}
Context: {context}

Format as a numbered list with assigned team member (if applicable).""",
                },
            ],
            temperature=0.7,
        )

        tasks = response.choices[0].message.content
        print("\nSuggested Tasks:\n")
        print(tasks)
        return tasks

    # ============================================
    # File Management
    # ============================================

    async def setup_team_storage(self, team: Team):
        """Set up file storage for the team."""
        print(f"\nSetting up team storage...")

        buckets = await client.storage.buckets.list()
        print(f"   Available buckets: {len(buckets.data)}")

        if buckets.data:
            bucket = buckets.data[0]
            print(f"   Using bucket: {bucket.name}")
            print(f"   Region: {bucket.region}")

            # List existing files
            files = await client.storage.files.list(bucket.id)
            print(f"   Existing files: {len(files.data)}")

    # ============================================
    # Workflow Automation
    # ============================================

    async def setup_team_workflows(self, team: Team):
        """Set up automation workflows for the team."""
        print(f"\nSetting up team workflows...")

        workflows = await client.automation.workflows.list()
        active_workflows = [w for w in workflows.data if w.is_active]

        print(f"   Available workflows: {len(workflows.data)}")
        print(f"   Active workflows: {len(active_workflows)}")

        for workflow in active_workflows[:5]:
            print(f"\n   {workflow.name}")
            print(f"      Trigger: {workflow.trigger_type}")
            print(f"      ID: {workflow.id}")

    # ============================================
    # Team Dashboard
    # ============================================

    async def show_dashboard(self, team: Team):
        """Display a team dashboard."""
        print("\n" + "=" * 60)
        print(f"TEAM DASHBOARD: {team.name}")
        print("=" * 60)

        # Team info
        print(f"\nOverview")
        print(f"   Project ID: {team.project_id}")
        print(f"   Members: {len(team.members)}")

        # Members
        await self.list_members(team)

        # Recent activity (simulated)
        print(f"\nRecent Activity")
        print(f"   - Team created")
        for member in team.members[-3:]:
            print(f"   - {member.name} joined the team")

        # Quick actions
        print(f"\nQuick Actions")
        print("   1. Add new member")
        print("   2. Generate project summary")
        print("   3. Suggest tasks")
        print("   4. Upload file")
        print("   5. Trigger workflow")


# ============================================
# Main Application
# ============================================

async def main():
    """Run the team collaboration application."""
    app = TeamCollaborationApp()

    try:
        # Create a new team
        team = await app.create_team(
            name="Project Phoenix",
            description="A revolutionary new product development initiative",
        )

        # Add team members
        await app.add_member(team, "alice@example.com", "Alice Johnson", "Project Lead")
        await app.add_member(team, "bob@example.com", "Bob Smith", "Developer")
        await app.add_member(team, "carol@example.com", "Carol Williams", "Designer")

        # Show dashboard
        await app.show_dashboard(team)

        # Generate AI-powered project summary
        await app.generate_project_summary(team)

        # Get task suggestions
        await app.suggest_tasks(
            team,
            context="We're starting a new mobile app project. We need to set up the development environment and create initial designs.",
        )

        # Set up storage and workflows
        await app.setup_team_storage(team)
        await app.setup_team_workflows(team)

        # Cleanup (comment out to keep the team)
        print("\nCleaning up...")
        for member in team.members:
            try:
                await client.users.delete(member.user_id)
                print(f"   Deleted user: {member.name}")
            except NotFoundError:
                pass

        await client.projects.delete(team.project_id)
        print(f"   Deleted project: {team.name}")
        print("Cleanup complete")

    except Exception as e:
        print(f"Application error: {e}")
        raise


if __name__ == "__main__":
    asyncio.run(main())
