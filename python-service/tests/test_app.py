from __future__ import annotations

from typing import Any

import jwt
import pytest
from fastapi import FastAPI
from httpx import ASGITransport, AsyncClient

from app.client import GoServiceClient
from app.config import Settings
from app.main import app, get_go_client, get_settings_dep


class FakeGoClient(GoServiceClient):
    def __init__(self) -> None:
        self.calls: list[dict[str, Any]] = []

    async def close(self) -> None:
        return None

    async def health(self) -> dict[str, Any]:
        return {"status": "ok"}

    async def request_token(self, username: str, password: str) -> dict[str, Any]:
        return {"token": f"{username}:{password}", "type": "Bearer"}

    async def forward_payload(self, payload: dict[str, Any], token: str) -> dict[str, Any]:
        self.calls.append({"payload": payload, "token": token})
        total_amount = sum(item["quantity"] * item["price"] for item in payload["items"])
        return {
            "request_id": payload["request_id"],
            "approved_by": "student",
            "items_count": len(payload["items"]),
            "total_amount": total_amount,
            "tags": payload["metadata"]["tags"],
            "status": "accepted",
        }


@pytest.fixture
def test_app() -> FastAPI:
    fake_client = FakeGoClient()
    settings = Settings(
        go_service_url="http://go-service:8080",
        jwt_secret="super-secret-key",
        request_timeout=5.0,
    )

    app.dependency_overrides[get_go_client] = lambda: fake_client
    app.dependency_overrides[get_settings_dep] = lambda: settings
    yield app
    app.dependency_overrides.clear()


@pytest.fixture
async def async_client(test_app: FastAPI) -> AsyncClient:
    transport = ASGITransport(app=test_app)
    async with AsyncClient(transport=transport, base_url="http://testserver") as client:
        yield client


def build_token() -> str:
    return jwt.encode(
        {"sub": "student", "username": "student", "role": "integration-client"},
        "super-secret-key",
        algorithm="HS256",
    )


def build_payload() -> dict[str, Any]:
    return {
        "request_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
        "customer": "Ivan Petrov",
        "address": {
            "city": "Moscow",
            "street": "Tverskaya 1",
            "zip_code": "123456",
        },
        "items": [
            {"name": "keyboard", "quantity": 2, "price": 1500.50},
            {"name": "mouse", "quantity": 1, "price": 999.90},
        ],
        "metadata": {"priority": "high", "tags": ["study", "api"]},
    }


@pytest.mark.asyncio
async def test_health(async_client: AsyncClient) -> None:
    response = await async_client.get("/health")

    assert response.status_code == 200
    assert response.json() == {"status": "ok", "go_status": "ok"}


@pytest.mark.asyncio
async def test_forward_payload_success(async_client: AsyncClient) -> None:
    response = await async_client.post(
        "/api/forward",
        json=build_payload(),
        headers={"Authorization": f"Bearer {build_token()}"},
    )

    assert response.status_code == 200
    body = response.json()
    assert body["approved_by"] == "student"
    assert body["items_count"] == 2
    assert body["verified_subject"] == "student"
    assert body["total_amount"] == pytest.approx(4000.9)


@pytest.mark.asyncio
async def test_forward_payload_rejects_bad_token(async_client: AsyncClient) -> None:
    response = await async_client.post(
        "/api/forward",
        json=build_payload(),
        headers={"Authorization": "Bearer invalid-token"},
    )

    assert response.status_code == 401
    assert response.json()["detail"] == "token validation failed"


@pytest.mark.asyncio
async def test_forward_payload_validation_error(async_client: AsyncClient) -> None:
    payload = build_payload()
    payload["items"] = []

    response = await async_client.post(
        "/api/forward",
        json=payload,
        headers={"Authorization": f"Bearer {build_token()}"},
    )

    assert response.status_code == 422


@pytest.mark.asyncio
async def test_auth_token_proxy(async_client: AsyncClient) -> None:
    response = await async_client.post(
        "/auth/token",
        json={"username": "student", "password": "securepass123"},
    )

    assert response.status_code == 200
    assert response.json()["type"] == "Bearer"
