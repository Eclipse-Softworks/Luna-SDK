"""Type definitions for Luna SDK."""

from __future__ import annotations

from typing import TypeVar
from pydantic import BaseModel, Field


class PaginationParams(BaseModel):
    """Pagination parameters for list requests."""

    limit: int | None = Field(default=20, ge=1, le=100)
    cursor: str | None = None


T = TypeVar("T")


class ListResponse(BaseModel):
    """Generic list response with pagination."""

    has_more: bool
    next_cursor: str | None = None


class User(BaseModel):
    """User resource."""

    id: str
    email: str
    name: str
    avatar_url: str | None = None
    created_at: str
    updated_at: str


class UserCreate(BaseModel):
    """Parameters for creating a user."""

    email: str
    name: str
    avatar_url: str | None = None


class UserUpdate(BaseModel):
    """Parameters for updating a user."""

    name: str | None = None
    avatar_url: str | None = None


class UserList(ListResponse):
    """Paginated list of users."""

    data: list[User]


class Project(BaseModel):
    """Project resource."""

    id: str
    name: str
    description: str | None = None
    owner_id: str
    created_at: str
    updated_at: str


class ProjectCreate(BaseModel):
    """Parameters for creating a project."""

    name: str
    description: str | None = None


class ProjectUpdate(BaseModel):
    """Parameters for updating a project."""

    name: str | None = None
    description: str | None = None


class ProjectList(ListResponse):
    """Paginated list of projects."""

    data: list[Project]
