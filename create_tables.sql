CREATE TABLE IF NOT EXISTS url_to_short (
    real_url varchar(100) NOT NULL,
    short varchar(50) NOT NULL,
    PRIMARY KEY (real_url)
);

CREATE TABLE IF NOT EXISTS short_to_url (
    short varchar(50) NOT NULL,
    real_url varchar(100) NOT NULL,
    PRIMARY KEY (short)
);
