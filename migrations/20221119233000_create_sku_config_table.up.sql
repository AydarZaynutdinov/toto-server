CREATE TABLE IF NOT EXISTS sku_configs
(
    id             uuid         NOT NULL PRIMARY KEY,
    package        VARCHAR(255) NOT NULL,
    country_code   VARCHAR(2)   NOT NULL,
    percentile_min INT          NOT NULL,
    percentile_max INT          NOT NULL,
    main_sku       VARCHAR(255)   NOT NULL
);