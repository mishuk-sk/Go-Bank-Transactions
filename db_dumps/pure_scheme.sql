--
-- PostgreSQL database cluster dump
--

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Drop databases (except postgres and template1)
--

DROP DATABASE db;




--
-- Drop roles
--

DROP ROLE "postgres-dev";


--
-- Roles
--

CREATE ROLE "postgres-dev";
ALTER ROLE "postgres-dev" WITH SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS PASSWORD 'md53e792b8eb5d9fc742b7c797f798e811e';






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

UPDATE pg_catalog.pg_database SET datistemplate = false WHERE datname = 'template1';
DROP DATABASE template1;
--
-- Name: template1; Type: DATABASE; Schema: -; Owner: postgres-dev
--

CREATE DATABASE template1 WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';


ALTER DATABASE template1 OWNER TO "postgres-dev";

\connect template1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: DATABASE template1; Type: COMMENT; Schema: -; Owner: postgres-dev
--

COMMENT ON DATABASE template1 IS 'default template for new databases';


--
-- Name: template1; Type: DATABASE PROPERTIES; Schema: -; Owner: postgres-dev
--

ALTER DATABASE template1 IS_TEMPLATE = true;


\connect template1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: DATABASE template1; Type: ACL; Schema: -; Owner: postgres-dev
--

REVOKE CONNECT,TEMPORARY ON DATABASE template1 FROM PUBLIC;
GRANT CONNECT ON DATABASE template1 TO PUBLIC;


--
-- PostgreSQL database dump complete
--

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

--
-- Name: db; Type: DATABASE; Schema: -; Owner: postgres-dev
--

CREATE DATABASE db WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';


ALTER DATABASE db OWNER TO "postgres-dev";

\connect db

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: bankBooks; Type: TABLE; Schema: public; Owner: postgres-dev
--

CREATE TABLE public."bankBooks" (
    id uuid NOT NULL,
    balance money NOT NULL,
    "userId" uuid NOT NULL
);


ALTER TABLE public."bankBooks" OWNER TO "postgres-dev";

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: postgres-dev
--

CREATE TABLE public.transactions (
    id uuid NOT NULL,
    "fromBook" uuid NOT NULL,
    "toBook" uuid NOT NULL,
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
    phone character varying(13)[],
    second_name text NOT NULL,
    email character varying(40)[],
    CONSTRAINT "OneIsNotNULL" CHECK (((phone IS NOT NULL) OR (email IS NOT NULL)))
);


ALTER TABLE public.users OWNER TO "postgres-dev";

--
-- Data for Name: bankBooks; Type: TABLE DATA; Schema: public; Owner: postgres-dev
--

COPY public."bankBooks" (id, balance, "userId") FROM stdin;
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres-dev
--

COPY public.transactions (id, "fromBook", "toBook", date, money) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres-dev
--

COPY public.users (first_name, id, phone, second_name, email) FROM stdin;
\.


--
-- Name: bankBooks none; Type: CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public."bankBooks"
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
-- Name: bankBooks bankBooks_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public."bankBooks"
    ADD CONSTRAINT "bankBooks_userId_fkey" FOREIGN KEY ("userId") REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- Name: transactions transactions_fromBook_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT "transactions_fromBook_fkey" FOREIGN KEY ("fromBook") REFERENCES public."bankBooks"(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- Name: transactions transactions_toBook_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres-dev
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT "transactions_toBook_fkey" FOREIGN KEY ("toBook") REFERENCES public."bankBooks"(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- PostgreSQL database dump complete
--

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

DROP DATABASE postgres;
--
-- Name: postgres; Type: DATABASE; Schema: -; Owner: postgres-dev
--

CREATE DATABASE postgres WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';


ALTER DATABASE postgres OWNER TO "postgres-dev";

\connect postgres

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: DATABASE postgres; Type: COMMENT; Schema: -; Owner: postgres-dev
--

COMMENT ON DATABASE postgres IS 'default administrative connection database';


--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database cluster dump complete
--

