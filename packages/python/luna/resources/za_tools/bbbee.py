"""Luna SDK - B-BBEE Compliance Verification for Python."""

from __future__ import annotations

from luna.http import HttpClient
from .types import BBBEEConfig, BBBEECertificate, BBBEELevel, BBBEEScorecard


class BBBEE:
    """Broad-Based Black Economic Empowerment compliance verification."""

    def __init__(self, client: HttpClient, config: BBBEEConfig) -> None:
        self._client = client
        self._config = config

    async def verify_certificate(self, certificate_number: str) -> BBBEECertificate | None:
        """Verify a B-BBEE certificate by number."""
        try:
            response = await self._client.request(
                method="GET",
                path=f"/v1/za/bbbee/certificates/{certificate_number}",
            )
            return BBBEECertificate(**response.data)
        except Exception:
            return None

    async def lookup_by_company(self, registration_number: str) -> BBBEECertificate | None:
        """Look up B-BBEE status by company registration number."""
        try:
            response = await self._client.request(
                method="GET",
                path="/v1/za/bbbee/lookup",
                query={"registrationNumber": registration_number},
            )
            return BBBEECertificate(**response.data)
        except Exception:
            return None

    async def meets_requirement(
        self,
        registration_number: str,
        minimum_level: BBBEELevel,
    ) -> dict:
        """Check if a company meets minimum B-BBEE level requirement."""
        certificate = await self.lookup_by_company(registration_number)

        if not certificate or not certificate.is_valid:
            return {"meets": False}

        level_number = certificate.level if isinstance(certificate.level, int) else 9
        required_number = minimum_level if isinstance(minimum_level, int) else 9

        return {
            "meets": level_number <= required_number,
            "actual_level": certificate.level,
            "certificate": certificate,
        }

    def get_recognition_level(self, level: BBBEELevel) -> int:
        """Get procurement recognition percentage for a B-BBEE level."""
        recognition_map = {
            1: 135,
            2: 125,
            3: 110,
            4: 100,
            5: 80,
            6: 60,
            7: 50,
            8: 10,
            "non-compliant": 0,
        }
        return recognition_map.get(level, 0)

    def determine_enterprise_category(self, annual_revenue: float) -> dict:
        """Calculate EME or QSE status."""
        if annual_revenue <= 10_000_000:
            return {
                "category": "EME",
                "automatic_level": 4,
                "description": "Exempted Micro Enterprise - automatic Level 4 (or Level 1 with 51%+ black ownership)",
            }
        elif annual_revenue <= 50_000_000:
            return {
                "category": "QSE",
                "automatic_level": None,
                "description": "Qualifying Small Enterprise - simplified scorecard applies",
            }
        else:
            return {
                "category": "Generic",
                "automatic_level": None,
                "description": "Generic Enterprise - full scorecard applies",
            }

    def calculate_score(self, scorecard: BBBEEScorecard) -> dict:
        """Calculate scorecard points."""
        scores = {
            "ownership": min(scorecard.ownership, 25),
            "management_control": min(scorecard.management_control, 19),
            "skills_development": min(scorecard.skills_development, 20),
            "enterprise_supplier_development": min(scorecard.enterprise_supplier_development, 40),
            "socio_economic_development": min(scorecard.socio_economic_development, 5),
        }

        total = sum(scores.values())
        level = self._points_to_level(total)

        return {"total": total, "level": level}

    def _points_to_level(self, points: float) -> BBBEELevel:
        """Convert total points to B-BBEE level."""
        if points >= 100:
            return 1
        if points >= 95:
            return 2
        if points >= 90:
            return 3
        if points >= 80:
            return 4
        if points >= 75:
            return 5
        if points >= 70:
            return 6
        if points >= 55:
            return 7
        if points >= 40:
            return 8
        return "non-compliant"

    def is_certificate_valid(self, certificate: BBBEECertificate) -> bool:
        """Check if a certificate is still valid."""
        from datetime import datetime

        expiry_date = datetime.fromisoformat(certificate.expiry_date.replace("Z", "+00:00"))
        return certificate.is_valid and expiry_date > datetime.now(expiry_date.tzinfo)

    def get_days_until_expiry(self, certificate: BBBEECertificate) -> int:
        """Get days until certificate expiry."""
        from datetime import datetime

        expiry_date = datetime.fromisoformat(certificate.expiry_date.replace("Z", "+00:00"))
        diff = expiry_date - datetime.now(expiry_date.tzinfo)
        return diff.days
