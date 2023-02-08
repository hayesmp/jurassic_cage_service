CREATE TABLE IF NOT EXISTS cage (
                                    id     UUID NOT NULL DEFAULT uuid_generate_v4(),
                                    name   TEXT UNIQUE,
                                    status INT,
                                    dinosaur_count INT,
                                    predominate_eating_habit INT,

                                    CONSTRAINT pkey_id_cage PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS dinosaur (
                                        id           UUID NOT NULL DEFAULT uuid_generate_v4(),
                                        name         TEXT UNIQUE,
                                        eating_habit INT,
                                        species      INT,
                                        cage_id      UUID,

                                        CONSTRAINT pkey_id_dinosaur PRIMARY KEY (id)
);
