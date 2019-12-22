CREATE TABLE links (
       id varchar(20) primary key not null,
       url string not null,
       created_at timestamp not null default current_timestamp
);

CREATE UNIQUE INDEX IF NOT EXISTS links_url_idx ON links (url);
