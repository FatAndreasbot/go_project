-- +goose Up
SELECT 'up SQL query';

CREATE TABLE public.group_permissions (
    group_id uuid NOT NULL,
    permission_id uuid NOT NULL
);

CREATE TABLE public.groups (
    name character varying NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    idempotency_key uuid NOT NULL,
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY
);

CREATE TABLE public.permissions (
    name character varying NOT NULL UNIQUE,
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY
);

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    email character varying NOT NULL UNIQUE,
    password_hash character varying(128) NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    idempotency_key uuid NOT NULL UNIQUE,
    group_id uuid NOT NULL
);

ALTER TABLE ONLY public.group_permissions
    ADD CONSTRAINT group_permissions_pk PRIMARY KEY (group_id, permission_id);

ALTER TABLE ONLY public.group_permissions
    ADD CONSTRAINT group_permissions_groups_fk FOREIGN KEY (group_id) REFERENCES public.groups(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.group_permissions
    ADD CONSTRAINT group_permissions_permissions_fk FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.groups(id) ON DELETE RESTRICT;

INSERT INTO public."permissions" ("name")
VALUES 
    ('user_create'),
    ('user_read'),
    ('user_update'),
    ('user_delete'),
    ('group_create'),
    ('group_read'),
    ('group_update'),
    ('group_delete');

INSERT INTO public."groups" ("name", "idempotency_key")
	VALUES ('admin', gen_random_uuid());

INSERT INTO public.group_permissions (group_id, permission_id)
SELECT 
    g.id, p.id
FROM
    public.groups g
    CROSS JOIN public.permissions p
WHERE g.name = 'admin'

;


-- +goose Down
SELECT 'down SQL query';

DROP TABLE public.users;

DROP TABLE public.group_permissions;

DROP TABLE public.groups;

DROP TABLE public.permissions;
