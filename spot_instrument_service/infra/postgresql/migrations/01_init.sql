-- +goose Up
select "up SQL query";

CREATE TABLE public.Market(
    market_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name CHARACTER VARYING NOT NULL UNIQUE,
    enabled BOOLEAN NOT NULL DEFAULT true,
    deleted_at timestamp,
    created_at timestamp DEFAULT now() not null
)