create table router
(
    id          serial
        constraint router_pk
            primary key,
    name        varchar(255) default ''::character varying not null,
    description varchar(255) default ''::character varying not null,
    address     varchar(20)  default ''::character varying not null,
    login       varchar(50)  default ''::character varying not null,
    password    varchar(255) default ''::character varying not null,
    created_at  timestamp                                  not null,
    updated_at  timestamp,
    status      boolean      default false,
    work_time   boolean      default false                 not null
);

create unique index router_address_uindex
    on router (address);

create unique index router_name_uindex
    on router (name);

insert into public.router (id, name, description, address, login, password, created_at, updated_at, status, work_time)
values  (1, 'Router1', 'Коридор', '95.84.134.115:8728', 'admin', 'Vtlcgjgek1', '2021-04-07 19:29:35.000000', null, true, true);