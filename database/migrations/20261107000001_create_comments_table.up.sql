CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    store_id UUID NULL,
    product_id UUID NULL,
    content VARCHAR(2000) NOT NULL,
    rating SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT comments_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT comments_store_id_fkey
        FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE,
    CONSTRAINT comments_product_id_fkey
        FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT comments_exactly_one_target_check
        CHECK ((store_id IS NOT NULL AND product_id IS NULL) OR
               (store_id IS NULL AND product_id IS NOT NULL)),
    CONSTRAINT comments_rating_check CHECK (rating BETWEEN 1 AND 5)
);

CREATE INDEX IF NOT EXISTS idx_comments_store_created_at
    ON comments(store_id, created_at DESC) WHERE store_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_comments_product_created_at
    ON comments(product_id, created_at DESC) WHERE product_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
