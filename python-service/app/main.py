from __future__ import annotations

from contextlib import asynccontextmanager
from typing import AsyncIterator

from fastapi import Depends, FastAPI, Header, HTTPException, Request, status
from httpx import HTTPError

from app.auth import validate_jwt_token
from app.client import GoServiceClient
from app.config import Settings, get_settings
from app.models import ForwardedResponse, LoginRequest, ProcessRequest


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncIterator[None]:
    settings = get_settings()
    app.state.settings = settings
    app.state.go_client = GoServiceClient(settings)
    yield
    await app.state.go_client.close()


app = FastAPI(
    title="python-service",
    version="1.0.0",
    lifespan=lifespan,
)


def get_go_client(request: Request) -> GoServiceClient:
    return request.app.state.go_client


def get_settings_dep(request: Request) -> Settings:
    return request.app.state.settings


@app.get("/health")
async def health(go_client: GoServiceClient = Depends(get_go_client)) -> dict[str, str]:
    go_health = await go_client.health()
    return {"status": "ok", "go_status": go_health["status"]}


@app.post("/auth/token")
async def auth_token(
    credentials: LoginRequest,
    go_client: GoServiceClient = Depends(get_go_client),
) -> dict:
    try:
        return await go_client.request_token(credentials.username, credentials.password)
    except HTTPError as exc:
        raise HTTPException(
            status_code=status.HTTP_502_BAD_GATEWAY,
            detail=f"go-service auth failed: {exc}",
        ) from exc


@app.post("/api/forward", response_model=ForwardedResponse)
async def forward_payload(
    payload: ProcessRequest,
    authorization: str = Header(..., alias="Authorization"),
    go_client: GoServiceClient = Depends(get_go_client),
    settings: Settings = Depends(get_settings_dep),
) -> ForwardedResponse:
    if not authorization.startswith("Bearer "):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="invalid authorization header",
        )

    token = authorization.split(" ", 1)[1].strip()
    claims = validate_jwt_token(token, settings.jwt_secret)

    try:
        response = await go_client.forward_payload(payload.model_dump(mode="json"), token)
    except HTTPError as exc:
        raise HTTPException(
            status_code=status.HTTP_502_BAD_GATEWAY,
            detail=f"go-service request failed: {exc}",
        ) from exc

    return ForwardedResponse(**response, verified_subject=claims["sub"])
