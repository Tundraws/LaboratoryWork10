from __future__ import annotations

from typing import Any

import httpx

from app.config import Settings


class GoServiceClient:
    def __init__(self, settings: Settings) -> None:
        self._settings = settings
        self._client = httpx.AsyncClient(
            base_url=settings.go_service_url,
            timeout=settings.request_timeout,
        )

    async def close(self) -> None:
        await self._client.aclose()

    async def health(self) -> dict[str, Any]:
        response = await self._client.get("/health")
        response.raise_for_status()
        return response.json()

    async def request_token(self, username: str, password: str) -> dict[str, Any]:
        response = await self._client.post(
            "/auth/token",
            json={"username": username, "password": password},
        )
        response.raise_for_status()
        return response.json()

    async def forward_payload(
        self,
        payload: dict[str, Any],
        token: str,
    ) -> dict[str, Any]:
        response = await self._client.post(
            "/api/process",
            json=payload,
            headers={"Authorization": f"Bearer {token}"},
        )
        response.raise_for_status()
        return response.json()
