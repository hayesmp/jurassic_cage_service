CREATE TABLE IF NOT EXISTS cage (
    id     UUID NOT NULL DEFAULT uuid_generate_v4(),
    name   TEXT,
    status INT,

    CONSTRAINT pkey_id_cage PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS dinosaur (
    id           UUID NOT NULL DEFAULT uuid_generate_v4(),
    name         TEXT,
    eating_habit INT,
    species      INT,
    cage_id      UUID,

    CONSTRAINT pkey_id_dinosaur PRIMARY KEY (id),
    CONSTRAINT FK_dinosaur_cage FOREIGN KEY (cage_id)
        REFERENCES cage(id)
);
