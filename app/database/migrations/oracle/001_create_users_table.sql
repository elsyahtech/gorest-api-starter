DECLARE
    table_count NUMBER;
BEGIN
    SELECT COUNT(*)
    INTO table_count
    FROM user_tables
    WHERE table_name = 'USERS';

    IF table_count = 0 THEN
        EXECUTE IMMEDIATE '
            CREATE TABLE users (
                id VARCHAR2(36 CHAR) NOT NULL,
                email VARCHAR2(255 CHAR) NOT NULL,
                password VARCHAR2(255 CHAR) NOT NULL,
                is_active NUMBER(1) DEFAULT 1 NOT NULL,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                updated_at TIMESTAMP NULL,
                deleted_at TIMESTAMP NULL,
                CONSTRAINT pk_users PRIMARY KEY (id),
                CONSTRAINT uq_users_email UNIQUE (email),
                CONSTRAINT chk_users_is_active CHECK (is_active IN (0, 1))
            )
        ';
    END IF;
END;
