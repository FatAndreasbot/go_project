-- +goose Up
SELECT 'up SQL query';

CREATE TABLE public.group_permissions (
    group_id uuid NOT NULL,
    permission_id uuid NOT NULL
);

CREATE TABLE public.groups (
    name character varying NOT NULL,
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY
);

CREATE TABLE public.permissions (
    name character varying NOT NULL UNIQUE,
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY
);

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    username character varying NOT NULL UNIQUE,
    password_hash character varying(128) NOT NULL,
    group_id uuid NOT NULL
);

ALTER TABLE ONLY public.group_permissions
    ADD CONSTRAINT group_permissions_pk PRIMARY KEY (group_id, permission_id);

ALTER TABLE ONLY public.group_permissions
    ADD CONSTRAINT group_permissions_groups_fk FOREIGN KEY (group_id) REFERENCES public.groups(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.group_permissions
    ADD CONSTRAINT group_permissions_permissions_fk FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.groups(id);

-- +goose Down
SELECT 'down SQL query';

DROP TABLE public.users;

DROP TABLE public.group_permissions;

DROP TABLE public.groups;

DROP TABLE public.permissions;
