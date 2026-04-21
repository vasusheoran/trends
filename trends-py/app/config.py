from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file=".env", env_file_encoding="utf-8", extra="ignore")

    database_url: str = "postgresql://trends:trends@localhost:5432/trends"
    port: int = 5001
    zerodha_access_token: str = ""
    zerodha_api_key: str = ""
    # Instrument token for Nifty 50 on Zerodha
    zerodha_nifty_token: int = 256265
    excel_seed_path: str = "../data/Nifty-17-04-2026.xlsx"
    # Minimum bars needed before futures are calculated
    futures_min_bars: int = 100


settings = Settings()
