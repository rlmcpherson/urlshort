CREATE TABLE urls
(
        short character varying(20) NOT NULL PRIMARY KEY,
        url character varying(3000) NOT NULL,
        created date not null default CURRENT_DATE
);
