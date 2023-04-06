CREATE TABLE histories (
    history_uid uuid NOT NULL PRIMARY KEY,
    created_at timestamp DEFAULT now(),
    updated_at timestamp,
    deleted_at timestamp
);

CREATE TABLE statuses (
    status_uid uuid NOT NULL PRIMARY KEY,
    status varchar(30) NOT NULL,
    UNIQUE(status)
);

CREATE TABLE measure_units (
    mu_uid uuid NOT NULL PRIMARY KEY,
    unit varchar(30) NOT NULL,
    UNIQUE(unit)
);

CREATE TABLE countries (
    country_uid uuid NOT NULL PRIMARY KEY,
    country varchar(30) NOT NULL,
    UNIQUE(country)
);

CREATE TABLE users (
    user_uid uuid NOT NULL PRIMARY KEY,
    username varchar(60) NOT NULL,
    pw_hash varchar(255) NOT NULL,
    email varchar(60) NOT NULL,
    phone varchar(30),
    country_uid uuid NOT NULL REFERENCES countries(country_uid),
    history_uid uuid NOT NULL REFERENCES histories(history_uid),
    UNIQUE(username),
    UNIQUE(email)
);

CREATE TABLE shops (
    shop_uid uuid NOT NULL PRIMARY KEY,
    name varchar(60) NOT NULL,
    description text,
    history_uid uuid REFERENCES histories(history_uid),
    UNIQUE(name)
);

CREATE TABLE items (
    item_uid uuid NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    vendor varchar(60) NOT NULL,
    price numeric(19, 2) NOT NULL CHECK (price > 0),
    description text,
    in_stock int NOT NULL CHECK (in_stock >= 0),
    shop_uid uuid REFERENCES shops(shop_uid),
    history_uid uuid NOT NULL REFERENCES histories(history_uid)
);

CREATE TABLE shop_ratings (
    rating_uid uuid NOT NULL PRIMARY KEY,
    user_uid uuid REFERENCES users(user_uid),
    shop_uid uuid REFERENCES shops(shop_uid),
    mark smallint NOT NULL CHECK (mark >= 0),
    commentary text,
    history_uid uuid NOT NULL REFERENCES histories(history_uid)
);

CREATE TABLE item_ratings (
    rating_uid uuid NOT NULL PRIMARY KEY,
    user_uid uuid REFERENCES users(user_uid),
    item_uid uuid REFERENCES items(item_uid),
    mark smallint NOT NULL CHECK (mark >= 0),
    commentary text,
    history_uid uuid NOT NULL REFERENCES histories(history_uid)
);

CREATE TABLE orders (
    order_uid uuid NOT NULL PRIMARY KEY,
    user_uid uuid REFERENCES users(user_uid),
    status_uid uuid REFERENCES statuses(status_uid),
    history_uid uuid REFERENCES histories(history_uid)
);

CREATE TABLE order_to_item (
    oti_uid uuid NOT NULL PRIMARY KEY,
    order_uid uuid REFERENCES orders(order_uid),
    item_uid uuid REFERENCES items(item_uid),
    quantity numeric(10, 2) NOT NULL CHECK (quantity > 0),
    measure_unit_id uuid REFERENCES measure_units(mu_uid)
);
