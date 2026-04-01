from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    go_service_url: str = Field(default="http://localhost:8080")
    jwt_secret: str = Field(default="super-secret-key")
    request_timeout: float = Field(default=5.0, gt=0)

    model_config = SettingsConfigDict(env_prefix="PYTHON_SERVICE_", extra="ignore")


def get_settings() -> Settings:
    return Settings()
