	CREATE TABLE public.payments
	(
		id bigserial NOT NULL,
		"from" bigint NOT NULL,
		"to" bigint NOT NULL,
		amount numeric(22,4) NOT NULL DEFAULT 0,
		date timestamp with time zone NOT NULL DEFAULT now(),
		CONSTRAINT payments_pk PRIMARY KEY (id)
	)
	
	TABLESPACE pg_default;
	
	ALTER TABLE public.payments
		OWNER to coins;
	
	CREATE INDEX payments_from_to_idx
		ON public.payments USING btree
		("from" ASC NULLS LAST, "to" ASC NULLS LAST)
		TABLESPACE pg_default;


	CREATE TABLE public.accounts
	(
		id bigserial NOT NULL,
		name character varying(32) COLLATE pg_catalog."default" NOT NULL,
		balance numeric(22,4) NOT NULL DEFAULT 0,
		currency character varying COLLATE pg_catalog."default" NOT NULL,
		CONSTRAINT accounts_pk PRIMARY KEY (id),
		CONSTRAINT accounts_name UNIQUE (name)
	)
	
	TABLESPACE pg_default;
	
	ALTER TABLE public.accounts
		OWNER to coins;