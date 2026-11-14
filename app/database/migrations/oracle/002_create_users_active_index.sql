DECLARE
    index_count NUMBER;
BEGIN
    SELECT COUNT(*)
    INTO index_count
    FROM user_indexes
    WHERE index_name = 'IX_USERS_IS_ACTIVE';

    IF index_count = 0 THEN
        EXECUTE IMMEDIATE 'CREATE INDEX ix_users_is_active ON users (is_active)';
    END IF;
END;
