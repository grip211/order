create table public.order
(
    order_uid          varchar(19) primary key  not null,
    track_number       varchar                  not null,
    entry              varchar                  not null,
    -- "delivery": {
    --     "name": "Test Testov",
    --     "phone": "+9720000000",
    --     "zip": "2639809",
    --     "city": "Kiryat Mozkin",
    --     "address": "Ploshad Mira 15",
    --     "region": "Kraiot",
    --     "email": "test@gmail.com"
    -- },
    delivery           jsonb                             default '{
      "name": "",
      "phone": ""
    }'::jsonb not null,
    locale             varchar                  not null,
    internal_signature varchar                  null,
    customer_id        varchar                  not null,
    delivery_service   varchar                  not null,
    shardkey           varchar                  not null,
    sm_id              int                      not null,
    date_created       timestamp with time zone not null default now(),
    oof_shard          varchar                  not null
);

create extension if not exists pg_trgm;

-- create index if not exists idx_tgrm_track_number
--    on public.order using gin (track_number gin_trgm_ops);

create unique index if not exists uq_track_number
    on public.order (track_number);

create table public.payment
(
    transaction_id varchar(19) primary key not null,
    request_id     varchar                 null,
    currency       varchar                 not null,
    provider       varchar                 not null,
    amount         decimal                 not null,
    payment_dt     bigint                  not null,
    bank           varchar                 not null,
    delivery_cost  decimal                 not null,
    goods_total    decimal                 not null,
    custom_fee     decimal                 not null,
    FOREIGN KEY (transaction_id) REFERENCES public.order (order_uid)
);

create table public.order_item
(
    chrt_id      int primary key not null,
    track_number varchar         not null,
    rid          varchar         not null,
    price        decimal         not null,
    name         varchar         not null,
    brand        varchar         not null,
    sale         decimal         not null,
    size         varchar         null,
    total_price  decimal         not null,
    nm_id        int             not null,
    status       smallint        not null,
    FOREIGN KEY (track_number) REFERENCES public.order (track_number)
);
