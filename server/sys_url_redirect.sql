--
-- PostgreSQL database dump
--

-- Dumped from database version 14.7
-- Dumped by pg_dump version 14.7

-- Started on 2023-06-01 16:45:14

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;


-- Database Structure For sys_url_redirect
CREATE DATABASE sys_url_redirect ENCODING 'UTF8';

-- Connect to database sys_url_redirect
\c sys_url_redirect

--
-- TOC entry 210 (class 1259 OID 16481)
-- Name: access_logs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.access_logs (
    id integer NOT NULL,
    from_domain character varying(200) NOT NULL,
    to_domain character varying(200),
    access_time timestamp with time zone DEFAULT now() NOT NULL,
    ip character varying(64),
    user_agent character varying(1000),
    rule_id integer,
    referer character varying(1000),
    uv_cookie character varying(200),
    log_uuid character varying(200) NOT NULL
);


ALTER TABLE public.access_logs OWNER TO postgres;

--
-- TOC entry 209 (class 1259 OID 16480)
-- Name: access_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.access_logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.access_logs_id_seq OWNER TO postgres;

--
-- TOC entry 3407 (class 0 OID 0)
-- Dependencies: 209
-- Name: access_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.access_logs_id_seq OWNED BY public.access_logs.id;


--
-- TOC entry 218 (class 1259 OID 16687)
-- Name: rule_domains; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rule_domains (
    id integer NOT NULL,
    rule_id integer NOT NULL,
    from_domain character varying(200) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.rule_domains OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 16686)
-- Name: rule_domains_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.rule_domains_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.rule_domains_id_seq OWNER TO postgres;

--
-- TOC entry 3408 (class 0 OID 0)
-- Dependencies: 217
-- Name: rule_domains_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.rule_domains_id_seq OWNED BY public.rule_domains.id;


--
-- TOC entry 222 (class 1259 OID 16773)
-- Name: rule_shares; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rule_shares (
    id integer NOT NULL,
    rule_id integer NOT NULL,
    password text NOT NULL,
    share_url text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.rule_shares OWNER TO postgres;

--
-- TOC entry 221 (class 1259 OID 16772)
-- Name: rule_shares_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.rule_shares_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.rule_shares_id_seq OWNER TO postgres;

--
-- TOC entry 3409 (class 0 OID 0)
-- Dependencies: 221
-- Name: rule_shares_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.rule_shares_id_seq OWNED BY public.rule_shares.id;


--
-- TOC entry 212 (class 1259 OID 16514)
-- Name: rules; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rules (
    id integer NOT NULL,
    rule_data text NOT NULL,
    status smallint DEFAULT 1 NOT NULL,
    remark character varying(200) DEFAULT ''::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    app_name character varying(100) DEFAULT ''::character varying NOT NULL,
    updated_at timestamp with time zone DEFAULT now(),
    default_url character varying(200) DEFAULT ''::character varying,
    ip_blacks text[]
);


ALTER TABLE public.rules OWNER TO postgres;

--
-- TOC entry 211 (class 1259 OID 16513)
-- Name: rules_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.rules_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.rules_id_seq OWNER TO postgres;

--
-- TOC entry 3410 (class 0 OID 0)
-- Dependencies: 211
-- Name: rules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.rules_id_seq OWNED BY public.rules.id;


--
-- TOC entry 216 (class 1259 OID 16672)
-- Name: stats_days; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stats_days (
    id integer NOT NULL,
    rule_id integer NOT NULL,
    access_time timestamp with time zone NOT NULL,
    pv_num integer DEFAULT 0 NOT NULL,
    ip_num integer DEFAULT 0 NOT NULL,
    uv_num integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.stats_days OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 16671)
-- Name: stats_day_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.stats_day_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.stats_day_id_seq OWNER TO postgres;

--
-- TOC entry 3411 (class 0 OID 0)
-- Dependencies: 215
-- Name: stats_day_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.stats_day_id_seq OWNED BY public.stats_days.id;


--
-- TOC entry 220 (class 1259 OID 16735)
-- Name: stats_minutes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stats_minutes (
    id integer NOT NULL,
    rule_id integer NOT NULL,
    access_time timestamp with time zone NOT NULL,
    pv_num integer DEFAULT 0 NOT NULL,
    ip_num integer DEFAULT 0 NOT NULL,
    uv_num integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.stats_minutes OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 16734)
-- Name: stats_minutes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.stats_minutes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.stats_minutes_id_seq OWNER TO postgres;

--
-- TOC entry 3412 (class 0 OID 0)
-- Dependencies: 219
-- Name: stats_minutes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.stats_minutes_id_seq OWNED BY public.stats_minutes.id;


--
-- TOC entry 214 (class 1259 OID 16537)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    account character varying(200) NOT NULL,
    password text NOT NULL,
    name character varying(200) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 213 (class 1259 OID 16536)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 3413 (class 0 OID 0)
-- Dependencies: 213
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 3194 (class 2604 OID 16484)
-- Name: access_logs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.access_logs ALTER COLUMN id SET DEFAULT nextval('public.access_logs_id_seq'::regclass);


--
-- TOC entry 3209 (class 2604 OID 16690)
-- Name: rule_domains id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rule_domains ALTER COLUMN id SET DEFAULT nextval('public.rule_domains_id_seq'::regclass);


--
-- TOC entry 3215 (class 2604 OID 16776)
-- Name: rule_shares id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rule_shares ALTER COLUMN id SET DEFAULT nextval('public.rule_shares_id_seq'::regclass);


--
-- TOC entry 3196 (class 2604 OID 16517)
-- Name: rules id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rules ALTER COLUMN id SET DEFAULT nextval('public.rules_id_seq'::regclass);


--
-- TOC entry 3205 (class 2604 OID 16675)
-- Name: stats_days id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stats_days ALTER COLUMN id SET DEFAULT nextval('public.stats_day_id_seq'::regclass);


--
-- TOC entry 3212 (class 2604 OID 16738)
-- Name: stats_minutes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stats_minutes ALTER COLUMN id SET DEFAULT nextval('public.stats_minutes_id_seq'::regclass);


--
-- TOC entry 3204 (class 2604 OID 16547)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 3389 (class 0 OID 16481)
-- Dependencies: 210
-- Data for Name: access_logs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.access_logs (id, from_domain, to_domain, access_time, ip, user_agent, rule_id, referer, uv_cookie, log_uuid) FROM stdin;
\.


--
-- TOC entry 3397 (class 0 OID 16687)
-- Dependencies: 218
-- Data for Name: rule_domains; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rule_domains (id, rule_id, from_domain, created_at) FROM stdin;
\.


--
-- TOC entry 3401 (class 0 OID 16773)
-- Dependencies: 222
-- Data for Name: rule_shares; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rule_shares (id, rule_id, password, share_url, created_at) FROM stdin;
\.


--
-- TOC entry 3391 (class 0 OID 16514)
-- Dependencies: 212
-- Data for Name: rules; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rules (id, rule_data, status, remark, created_at, app_name, updated_at, default_url, ip_blacks) FROM stdin;
\.


--
-- TOC entry 3395 (class 0 OID 16672)
-- Dependencies: 216
-- Data for Name: stats_days; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stats_days (id, rule_id, access_time, pv_num, ip_num, uv_num) FROM stdin;
\.


--
-- TOC entry 3399 (class 0 OID 16735)
-- Dependencies: 220
-- Data for Name: stats_minutes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stats_minutes (id, rule_id, access_time, pv_num, ip_num, uv_num) FROM stdin;
\.


--
-- TOC entry 3393 (class 0 OID 16537)
-- Dependencies: 214
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, account, password, name) FROM stdin;
1	admin	AXhhKihAa2LaRwY5mftnngSPKDF4N9JignnQ4skynY8y	管理员
\.


--
-- TOC entry 3414 (class 0 OID 0)
-- Dependencies: 209
-- Name: access_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.access_logs_id_seq', 1, false);


--
-- TOC entry 3415 (class 0 OID 0)
-- Dependencies: 217
-- Name: rule_domains_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.rule_domains_id_seq', 1, false);


--
-- TOC entry 3416 (class 0 OID 0)
-- Dependencies: 221
-- Name: rule_shares_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.rule_shares_id_seq', 1, false);


--
-- TOC entry 3417 (class 0 OID 0)
-- Dependencies: 211
-- Name: rules_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.rules_id_seq', 1, false);


--
-- TOC entry 3418 (class 0 OID 0)
-- Dependencies: 215
-- Name: stats_day_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.stats_day_id_seq', 1, false);


--
-- TOC entry 3419 (class 0 OID 0)
-- Dependencies: 219
-- Name: stats_minutes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.stats_minutes_id_seq', 1, false);


--
-- TOC entry 3420 (class 0 OID 0)
-- Dependencies: 213
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 11, true);


--
-- TOC entry 3221 (class 2606 OID 17151)
-- Name: access_logs access_logs_log_uuid_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.access_logs
    ADD CONSTRAINT access_logs_log_uuid_key UNIQUE (log_uuid);


--
-- TOC entry 3223 (class 2606 OID 16489)
-- Name: access_logs access_logs_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.access_logs
    ADD CONSTRAINT access_logs_pk PRIMARY KEY (id);


--
-- TOC entry 3237 (class 2606 OID 16693)
-- Name: rule_domains rule_domain_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rule_domains
    ADD CONSTRAINT rule_domain_pk PRIMARY KEY (id);


--
-- TOC entry 3246 (class 2606 OID 16781)
-- Name: rule_shares rule_shares_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rule_shares
    ADD CONSTRAINT rule_shares_pk PRIMARY KEY (id);


--
-- TOC entry 3248 (class 2606 OID 16783)
-- Name: rule_shares rule_shares_rule_id_un; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rule_shares
    ADD CONSTRAINT rule_shares_rule_id_un UNIQUE (rule_id);


--
-- TOC entry 3227 (class 2606 OID 16523)
-- Name: rules ruleid_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rules
    ADD CONSTRAINT ruleid_pk PRIMARY KEY (id);


--
-- TOC entry 3234 (class 2606 OID 16680)
-- Name: stats_days stats_day_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stats_days
    ADD CONSTRAINT stats_day_pk PRIMARY KEY (id);


--
-- TOC entry 3244 (class 2606 OID 16743)
-- Name: stats_minutes stats_minutes_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stats_minutes
    ADD CONSTRAINT stats_minutes_pk PRIMARY KEY (id);


--
-- TOC entry 3240 (class 2606 OID 16733)
-- Name: rule_domains uni_rule_fromDomain; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rule_domains
    ADD CONSTRAINT "uni_rule_fromDomain" UNIQUE (from_domain);


--
-- TOC entry 3229 (class 2606 OID 16546)
-- Name: users users_account_un; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_account_un UNIQUE (account);


--
-- TOC entry 3231 (class 2606 OID 16544)
-- Name: users users_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pk PRIMARY KEY (id);


--
-- TOC entry 3217 (class 1259 OID 16492)
-- Name: access_logs_access_time_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX access_logs_access_time_idx ON public.access_logs USING btree (access_time);


--
-- TOC entry 3218 (class 1259 OID 16490)
-- Name: access_logs_from_domain_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX access_logs_from_domain_idx ON public.access_logs USING btree (from_domain);


--
-- TOC entry 3219 (class 1259 OID 16493)
-- Name: access_logs_ip_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX access_logs_ip_idx ON public.access_logs USING btree (ip);


--
-- TOC entry 3224 (class 1259 OID 16491)
-- Name: access_logs_to_domain_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX access_logs_to_domain_idx ON public.access_logs USING btree (to_domain);


--
-- TOC entry 3225 (class 1259 OID 16494)
-- Name: access_logs_ua_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX access_logs_ua_idx ON public.access_logs USING btree (user_agent);


--
-- TOC entry 3238 (class 1259 OID 16694)
-- Name: rule_from_domain_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX rule_from_domain_idx ON public.rule_domains USING btree (from_domain);


--
-- TOC entry 3232 (class 1259 OID 16681)
-- Name: stats_day_access_time_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX stats_day_access_time_idx ON public.stats_days USING btree (access_time);


--
-- TOC entry 3235 (class 1259 OID 16684)
-- Name: stats_day_rule_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX stats_day_rule_id_idx ON public.stats_days USING btree (rule_id);


--
-- TOC entry 3241 (class 1259 OID 16744)
-- Name: stats_minute_access_time_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX stats_minute_access_time_idx ON public.stats_minutes USING btree (access_time);


--
-- TOC entry 3242 (class 1259 OID 16746)
-- Name: stats_minute_rule_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX stats_minute_rule_id_idx ON public.stats_minutes USING btree (rule_id);


-- Completed on 2023-06-01 16:45:15

--
-- PostgreSQL database dump complete
--

