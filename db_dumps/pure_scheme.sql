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

--ALTER TABLE ONLY public.transactions DROP CONSTRAINT "transactions_toBook_fkey";
--ALTER TABLE ONLY public.transactions DROP CONSTRAINT "transactions_fromBook_fkey";
--ALTER TABLE ONLY public.personal_accounts DROP CONSTRAINT "bankBooks_userId_fkey";
--ALTER TABLE ONLY public.users DROP CONSTRAINT users_pkey;
--ALTER TABLE ONLY public.transactions DROP CONSTRAINT transactions_pkey;
--ALTER TABLE ONLY public.personal_accounts DROP CONSTRAINT "none";
--DROP TABLE public.users;
--DROP TABLE public.transactions;
--DROP TABLE public.personal_accounts;
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
    from_account uuid,
    to_account uuid,
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

