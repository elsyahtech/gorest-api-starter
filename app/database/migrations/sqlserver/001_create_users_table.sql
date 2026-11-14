IF OBJECT_ID(N'dbo.users', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.users (
        id VARCHAR(36) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        is_active BIT NOT NULL CONSTRAINT DF_users_is_active DEFAULT (1),
        created_at DATETIME2 NOT NULL CONSTRAINT DF_users_created_at DEFAULT SYSUTCDATETIME(),
        updated_at DATETIME2 NULL,
        deleted_at DATETIME2 NULL,
        CONSTRAINT PK_users PRIMARY KEY (id),
        CONSTRAINT UQ_users_email UNIQUE (email)
    );
END;

IF NOT EXISTS (
    SELECT 1
    FROM sys.indexes
    WHERE name = N'IX_users_is_active'
      AND object_id = OBJECT_ID(N'dbo.users')
)
BEGIN
    CREATE INDEX IX_users_is_active ON dbo.users (is_active);
END;
