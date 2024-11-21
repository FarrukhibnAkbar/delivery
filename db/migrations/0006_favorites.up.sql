CREATE SEQUENCE seq_favorites_id;
CREATE TABLE favorites (
    id BIGINT NOT NULL DEFAULT nextval('seq_favorites_id') PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    product_id UUID NOT NULL,
    is_favorited BOOLEAN DEFAULT true,
    added_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, product_id)
);