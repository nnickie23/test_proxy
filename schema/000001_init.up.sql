-- Table: public.tasks

CREATE SEQUENCE tasks_id_seq;

CREATE TABLE IF NOT EXISTS public.tasks
(
  id                        integer NOT NULL DEFAULT nextval('tasks_id_seq'::regclass)::integer,
  uuid                      text    NOT NULL DEFAULT gen_random_uuid()::text,
  status                    text    NOT NULL DEFAULT 'new'::text,
  request_method            text    NOT NULL DEFAULT ''::text,
  request_url               text    NOT NULL DEFAULT ''::text,
  request_headers           json    NOT NULL DEFAULT '{}'::json,
  response_http_status_code integer NOT NULL DEFAULT 0::integer,
  response_headers          json    NOT NULL DEFAULT '{}'::json,
  response_content_length   bigint  NOT NULL DEFAULT 0::bigint,
  CONSTRAINT tasks_pkey PRIMARY KEY (id),
  CONSTRAINT unique_task_id UNIQUE (uuid)
)

TABLESPACE pg_default;

ALTER TABLE public.tasks
  OWNER to postgres;