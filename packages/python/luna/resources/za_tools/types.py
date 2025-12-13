"""Luna SDK - South African Business Tools Types."""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Literal, Optional
from datetime import date


# CIPC Types
CompanyType = Literal["PTY_LTD", "LTD", "NPC", "CC", "INC"]
CompanyStatus = Literal["active", "in_business_rescue", "deregistered", "liquidated", "dissolved"]


@dataclass
class CIPCConfig:
    """CIPC configuration."""
    api_key: Optional[str] = None
    sandbox: bool = False


@dataclass
class Director:
    """Company director information."""
    name: str
    id_number: Optional[str] = None
    role: str = "Director"
    appointed_date: Optional[str] = None
    resigned_date: Optional[str] = None


@dataclass
class Company:
    """Company information."""
    registration_number: str
    name: str
    type: CompanyType
    status: CompanyStatus
    registration_date: Optional[str] = None
    financial_year_end: Optional[str] = None
    directors: list[Director] = field(default_factory=list)
    physical_address: Optional[str] = None
    postal_address: Optional[str] = None


# B-BBEE Types
BBBEELevel = int | Literal["non-compliant"]


@dataclass
class BBBEEConfig:
    """B-BBEE configuration."""
    api_key: Optional[str] = None
    sandbox: bool = False


@dataclass
class BBBEECertificate:
    """B-BBEE certificate information."""
    certificate_number: str
    company_name: str
    registration_number: str
    level: BBBEELevel
    verification_agency: str
    issue_date: str
    expiry_date: str
    is_valid: bool = True
    scorecard: Optional[dict] = None


@dataclass
class BBBEEScorecard:
    """B-BBEE scorecard breakdown."""
    ownership: float = 0
    management_control: float = 0
    skills_development: float = 0
    enterprise_supplier_development: float = 0
    socio_economic_development: float = 0


# ID Validation Types
@dataclass
class SAIDInfo:
    """South African ID information."""
    id_number: str
    is_valid: bool
    date_of_birth: date
    gender: Literal["male", "female"]
    is_sa_citizen: bool
    checksum_valid: bool


# Address Types
SAProvince = Literal["EC", "FS", "GP", "KZN", "LP", "MP", "NC", "NW", "WC"]


@dataclass
class SAAddress:
    """South African address."""
    street: Optional[str] = None
    suburb: Optional[str] = None
    city: Optional[str] = None
    province: Optional[SAProvince] = None
    postal_code: Optional[str] = None
    country: str = "ZA"


@dataclass
class ZAToolsConfig:
    """Combined ZA tools configuration."""
    cipc: Optional[CIPCConfig] = None
    bbbee: Optional[BBBEEConfig] = None
    strict: bool = False
