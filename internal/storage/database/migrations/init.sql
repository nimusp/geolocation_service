CREATE TABLE IF NOT EXISTS geolocation (
    ip_address      VARCHAR(39) NOT NULL, -- ip v6 max length
    country_code    VARCHAR(2) NOT NULL,  -- Alpha 2 format
    country         VARCHAR(200) NOT NULL,
    city            VARCHAR(200) NOT NULL,
    latitude        NUMERIC(7,5) NOT NULL,
    longitude       NUMERIC(8,5) NOT NULL,
    mystery_value   BIGINT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ,
    PRIMARY KEY (ip_address)
);