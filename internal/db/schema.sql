CREATE TABLE entities (
    id INTEGER PRIMARY KEY,
    content TEXT UNIQUE,
    source TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE embeddings (
    entity_id INTEGER,
    model_name TEXT,
    vector BLOB,
    FOREIGN KEY(entity_id) REFERENCES entities(id)
);