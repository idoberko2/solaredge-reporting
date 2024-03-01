CREATE TABLE se_data (
    t       TIMESTAMPTZ NOT NULL PRIMARY KEY,
    value   INT NOT NULL
);
SELECT create_hypertable('se_data', 't');
