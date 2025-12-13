"""Luna SDK - South African ID Validation."""

from __future__ import annotations

import re
from dataclasses import dataclass
from datetime import date, datetime
from typing import Literal

from .types import SAIDInfo


class IDValidation:
    """
    South African ID number validation and parsing.

    SA ID format: YYMMDD SSSS C A Z
    - YYMMDD: Date of birth
    - SSSS: Gender (0000-4999 = female, 5000-9999 = male)
    - C: Citizenship (0 = SA citizen, 1 = permanent resident)
    - A: Usually 8 (was used for race classification, now obsolete)
    - Z: Luhn checksum digit
    """

    def validate(self, id_number: str) -> SAIDInfo:
        """Validate and parse a South African ID number."""
        cleaned = re.sub(r"[\s-]", "", id_number)

        if len(cleaned) != 13:
            return self._invalid_result(id_number, "ID number must be 13 digits")

        if not cleaned.isdigit():
            return self._invalid_result(id_number, "ID number must contain only digits")

        # Extract components
        year = int(cleaned[0:2])
        month = int(cleaned[2:4])
        day = int(cleaned[4:6])
        gender_digits = int(cleaned[6:10])
        citizenship_digit = int(cleaned[10:11])

        # Validate date
        full_year = 2000 + year if year <= 30 else 1900 + year
        try:
            date_of_birth = date(full_year, month, day)
        except ValueError:
            return self._invalid_result(id_number, "Invalid date of birth")

        # Validate checksum
        checksum_valid = self._validate_luhn(cleaned)

        return SAIDInfo(
            id_number=cleaned,
            is_valid=checksum_valid,
            date_of_birth=date_of_birth,
            gender="male" if gender_digits >= 5000 else "female",
            is_sa_citizen=citizenship_digit == 0,
            checksum_valid=checksum_valid,
        )

    def is_valid(self, id_number: str) -> bool:
        """Quick validation - returns boolean."""
        return self.validate(id_number).is_valid

    def get_date_of_birth(self, id_number: str) -> date | None:
        """Extract date of birth from ID number."""
        result = self.validate(id_number)
        return result.date_of_birth if result.is_valid else None

    def get_age(self, id_number: str) -> int | None:
        """Get age from ID number."""
        dob = self.get_date_of_birth(id_number)
        if not dob:
            return None

        today = date.today()
        age = today.year - dob.year
        if (today.month, today.day) < (dob.month, dob.day):
            age -= 1
        return age

    def is_over(self, id_number: str, age: int) -> bool:
        """Check if person is over a specific age."""
        current_age = self.get_age(id_number)
        return current_age is not None and current_age >= age

    def generate(
        self,
        date_of_birth: date | None = None,
        gender: Literal["male", "female"] | None = None,
        is_citizen: bool = True,
    ) -> str:
        """Generate a valid SA ID number (for testing)."""
        import random

        dob = date_of_birth or date(1990, 1, 1)
        year = str(dob.year)[-2:]
        month = str(dob.month).zfill(2)
        day = str(dob.day).zfill(2)

        # Gender sequence
        gender_min = 0 if gender == "female" else 5000
        gender_max = 4999 if gender == "female" else 9999
        gender_seq = str(random.randint(gender_min, gender_max)).zfill(4)

        # Citizenship
        citizen = "0" if is_citizen else "1"

        # A digit
        a_digit = "8"

        # Calculate checksum
        partial = f"{year}{month}{day}{gender_seq}{citizen}{a_digit}"
        checksum = self._calculate_luhn_checksum(partial)

        return f"{partial}{checksum}"

    def _invalid_result(self, id_number: str, reason: str) -> SAIDInfo:
        """Return invalid result."""
        return SAIDInfo(
            id_number=id_number,
            is_valid=False,
            date_of_birth=date(1, 1, 1),
            gender="male",
            is_sa_citizen=False,
            checksum_valid=False,
        )

    def _validate_luhn(self, number: str) -> bool:
        """Validate using Luhn algorithm."""
        total = 0
        is_even = False

        for i in range(len(number) - 1, -1, -1):
            digit = int(number[i])

            if is_even:
                digit *= 2
                if digit > 9:
                    digit -= 9

            total += digit
            is_even = not is_even

        return total % 10 == 0

    def _calculate_luhn_checksum(self, partial: str) -> str:
        """Calculate Luhn checksum digit."""
        for i in range(10):
            if self._validate_luhn(partial + str(i)):
                return str(i)
        return "0"
