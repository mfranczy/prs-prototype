CREATE SEQUENCE pkg_id_seq;
CREATE SEQUENCE sub_id_seq;
CREATE SEQUENCE cmp_id_seq;
CREATE SEQUENCE rep_id_seq;
CREATE SEQUENCE pr_id_seq;
CREATE SEQUENCE lb_id_seq;
CREATE SEQUENCE lbf_id_seq;
CREATE SEQUENCE lfv_id_seq;
CREATE SEQUENCE pl_id_seq;

CREATE TABLE subsystems (
    sub_id bigint default nextval('sub_id_seq') PRIMARY KEY,
    sub_name varchar(255) UNIQUE,
    sub_description text,
    sub_modify_date timestamp default now(),
    sub_create_date timestamp default now()
);

CREATE TABLE components (
    cmp_id bigint default nextval('cmp_id_seq') PRIMARY KEY,
    cmp_name varchar(255) UNIQUE,
    cmp_description text,
    cmp_modify_date timestamp default now(),
    cmp_create_date timestamp default now()
);

CREATE TABLE repos (
    rep_id bigint default nextval('rep_id_seq') PRIMARY KEY,
    rep_name varchar(255) UNIQUE,
    rep_description text,
    rep_modify_date timestamp default now(),
    rep_create_date timestamp default now()
);

CREATE TABLE packages (
    pkg_id bigint default nextval('pkg_id_seq') PRIMARY KEY,
    pkg_name varchar(255) not null,
    pkg_version varchar(255) not null,
    pkg_maintainer varchar(255) not null,
    pkg_sub_id bigint REFERENCES subsystems (sub_id),
    pkg_cmp_id bigint REFERENCES components (cmp_id),
    pkg_comments text,
    pkg_summary text,
    UNIQUE (pkg_name, pkg_version)
);

CREATE TABLE packages_repos (
    pr_id bigint default nextval('pr_id_seq'),
    pr_pkg_id bigint REFERENCES packages (pkg_id),
    pr_rep_id bigint REFERENCES repos (rep_id),
    PRIMARY KEY (pr_id, pr_pkg_id, pr_rep_id)
);

CREATE TABLE labels (
    lb_id bigint default nextval('lb_id_seq') PRIMARY KEY,
    lb_name varchar(255) UNIQUE,
    lb_summary text,
    lb_color varchar(6) default 'ff0000',
    lb_modify_date timestamp default now(),
    lb_create_date timestamp default now(),
    lb_active boolean default true
);

CREATE TABLE labels_fields (
    lbf_id bigint default nextval('lbf_id_seq') PRIMARY KEY,
    lbf_lb_id bigint not null REFERENCES labels (lb_id),
    lbf_name varchar(255) not null,
    lbf_modify_date timestamp default now(),
    lbf_create_date timestamp default now(),
    lbf_description text,
    lbf_active boolean default true,
    UNIQUE (lbf_lb_id, lbf_name)
);

CREATE TABLE packages_labels (
    pl_id bigint not null default nextval('pl_id_seq') PRIMARY KEY,
    pl_pkg_id bigint not null REFERENCES packages (pkg_id),
    pl_lb_id bigint not null REFERENCES labels (lb_id),
    pl_active boolean default true,
    UNIQUE (pl_pkg_id, pl_lb_id)
);

CREATE TABLE labels_fields_values (
    lfv_id bigint not null default nextval('lfv_id_seq') PRIMARY KEY,
    lfv_pl_id bigint not null REFERENCES packages_labels (pl_id),
    lfv_lbf_id bigint not null REFERENCES labels_fields (lbf_id),
    lfv_value text,
    UNIQUE (lfv_pl_id, lfv_lbf_id)
);
