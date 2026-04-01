from typing import Any

import jwt
from fastapi import HTTPException, status


def validate_jwt_token(token: str, secret: str) -> dict[str, Any]:
    try:
        payload = jwt.decode(token, secret, algorithms=["HS256"])
    except jwt.InvalidTokenError as exc:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="token validation failed",
        ) from exc

    return payload
