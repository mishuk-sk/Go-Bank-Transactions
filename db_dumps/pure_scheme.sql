--
-- PostgreSQL database dump
--

-- Dumped from database version 11.2 (Debian 11.2-1.pgdg90+1)
-- Dumped by pg_dump version 11.2 (Debian 11.2-1.pgdg90+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE ONLY public.transactions DROP CONSTRAINT "transactions_toBook_fkey";
ALTER TABLE ONLY public.transactions DROP CONSTRAINT "transactions_fromBook_fkey";
ALTER TABLE ONLY public.personal_accounts DROP CONSTRAINT "bankBooks_userId_fkey";
ALTER TABLE ONLY public.users DROP CONSTRAINT users_pkey;
ALTER TABLE ONLY public.transactions DROP CONSTRAINT transactions_pkey;
ALTER TABLE ONLY public.personal_accounts DROP CONSTRAINT "none";
DROP TABLE public.users;
DROP TABLE public.transactions;
DROP TABLE public.personal_accounts;
SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: personal_accounts; Type: TABLE; Schema: public; Owner: postgres-dev
--

CREATE TABLE public.personal_accounts (
    id uuid NOT NULL,
    balance money NOT NULL,
    user_id uuid NOT NULL,
    name text
);


ALTER TABLE public.personal_accounts OWNER TO "postgres-dev";

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: postgres-dev
--

CREATE TABLE public.transactions (
    id uuid NOT NULL,
    from_account uuid NOT NULL,
    to_account uuid NOT NULL,
    date date NOT NULL,
    money money NOT NULL
);


ALTER TABLE public.transactions OWNER TO "postgres-dev";

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres-dev
--

CREATE TABLE public.users (
    first_name text NOT NULL,
    id uuid NOT NULL,
    phone text,
    second_name text NOT NULL,
    email text,
    CONSTRAINT "OneIsNotNULL" CHECK (((phone IS NOT NULL) OR (email IS NOT NULL))),
    CONSTRAINT users_email_check CHECK ((email <> ''::text)),
    CONSTRAINT users_first_name_check CHECK ((first_name <> ''::text)),
    CONSTRAINT users_phone_check CHECK ((phone <> ''::text)),
    CONSTRAINT users_second_name_check CHECK ((second_name <> ''::text))
);


ALTER TABLE public.users OWNER TO "postgres-dev";

--
-- Data for Name: personal_accounts; Type: TABLE DATA; Schema: public; Owner: postgres-dev
--

COPY public.personal_accounts (id, balance, user_id, name) FROM stdin;
a144d7b4-69d3-4c37-b400-4584548cd78c	$0.00	86ebaad9-c2ff-4740-a53c-91736146a869	\N
a25a9644-e740-4cf0-ade9-a2c0448d383c	$0.00	86ebaad9-c2ff-4740-a53c-91736146a869	\N
0daa575a-2385-44f0-bc50-c62b8e0c5f0e	$0.00	86ebaad9-c2ff-4740-a53c-91736146a869	\N
a6d3e736-bdac-45a3-a0b5-56e1dee1b6d3	$12,345.00	86ebaad9-c2ff-4740-a53c-91736146a869	\N
456def90-dcdf-4efa-8762-19de0b07cb54	$509.90	86ebaad9-c2ff-4740-a53c-91736146a869	\N
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres-dev
--

COPY public.transactions (id, from_account, to_account, date, money) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres-dev
--

COPY public.users (first_name, id, phone, second_name, email) FROM stdin;
Mikhail	45729667-08e2-4ce2-aad9-fe5540820cf1	234555	Skuratovich	\N
Mikhail	0aa445cc-62aa-47af-bc41-c835f42d400a	234555	Skuratovich	\N
lol	54831cee-95d8-401e-be92-19695c1b6e4e	234555	lol	\N
lol	9cde7f19-54cd-4ee7-8cc2-3b4d0186490e	234555	lol	\N
lol	2080dc41-6f8c-43f5-a09f-e2903938ac68	234555	lol	\N
lol	de0ffa3b-9b60-4515-b9c6-35b8d253abfc	234555	lol	\N
lol	1f71f0d6-078f-4e70-8b27-f45bd85e6fe7	234555	lol	\N
l	a7a566de-e53c-4981-8e38-2abcdd19eea8	234555	lol	\N
l	9ccd9d9b-5ecc-4fa8-86b8-721b35b8baf1	234555	lol	mishuk
l	8ce42428-069f-4a8b-a8f7-c2d7d97fe767	234555	lol	mishuk
Someone	86ebaad9-c2ff-4740-a53c-91736146a869	234555	Skuratovich	THIS ONE IS TEST
\.


--
-- Name: personal_accounts none; Type: CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.personal_accounts
    ADD CONSTRAINT "none" PRIMARY KEY (id);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: personal_accounts bankBooks_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.personal_accounts
    ADD CONSTRAINT "bankBooks_userId_fkey" FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- Name: transactions transactions_fromBook_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT "transactions_fromBook_fkey" FOREIGN KEY (from_account) REFERENCES public.personal_accounts(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- Name: transactions transactions_toBook_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT "transactions_toBook_fkey" FOREIGN KEY (to_account) REFERENCES public.personal_accounts(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- PostgreSQL database dump complete
--

