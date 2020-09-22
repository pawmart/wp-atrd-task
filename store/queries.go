package store

const (
	queryInitSecretsTable = `
		CREATE TABLE IF NOT EXISTS "secrets" (
			"id" 		 UUID PRIMARY KEY, 
			"content" 	 text DEFAULT NULL, 
			"remaining_views" text DEFAULT NULL, 
			"created_at" timestamp NULL DEFAULT NULL,
			"expires_at" timestamp NULL DEFAULT NULL
		)
	`

	queryGetSecretByID = `
		SELECT 
			"id", "content", "remaining_views", "created_at", "expires_at" 
		FROM 
			"secrets" 
		WHERE 
			"id" = $1
		LIMIT 1
	`

	queryInsertSecret = `
		INSERT INTO "secrets" 
			("id", "content", "remaining_views", "created_at", "expires_at") 
		VALUES
			(:id, :content, :remaining_views, :created_at, :expires_at)
		RETURNING
			"id"
	`

	queryUpdateSecret = `
	UPDATE 
		"secrets" 
	SET 
		"content" = :content,
		"remaining_views" = :remaining_views
	WHERE
		"id" = :id
	`

	queryDeleteSecretByID = `
		DELETE FROM 
			"secrets"
		WHERE 
			"id" = $1
	`

	querySeedTestData = `
	INSERT INTO "secrets" 
		("id", "content", "remaining_views") 
	VALUES
		('b75ce598-f349-4c61-9246-2053e230187d', 'test', 999)
	ON CONFLICT ("id")
	DO NOTHING
	`
)
