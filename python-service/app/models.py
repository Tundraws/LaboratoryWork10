from typing import Literal
from uuid import UUID

from pydantic import BaseModel, ConfigDict, Field


class Address(BaseModel):
    city: str = Field(min_length=2, max_length=64)
    street: str = Field(min_length=3, max_length=128)
    zip_code: str = Field(pattern=r"^\d{6}$")


class Item(BaseModel):
    name: str = Field(min_length=2, max_length=64)
    quantity: int = Field(ge=1, le=100)
    price: float = Field(gt=0)


class Metadata(BaseModel):
    priority: Literal["low", "medium", "high"]
    tags: list[str] = Field(min_length=1, max_length=5)


class ProcessRequest(BaseModel):
    request_id: UUID
    customer: str = Field(min_length=3, max_length=64)
    address: Address
    items: list[Item] = Field(min_length=1, max_length=10)
    metadata: Metadata

    model_config = ConfigDict(str_strip_whitespace=True)


class LoginRequest(BaseModel):
    username: str = Field(min_length=3, max_length=32, pattern=r"^[a-zA-Z0-9]+$")
    password: str = Field(min_length=8, max_length=64)


class ForwardedResponse(BaseModel):
    request_id: UUID
    approved_by: str
    items_count: int
    total_amount: float
    tags: list[str]
    status: str
    verified_subject: str
