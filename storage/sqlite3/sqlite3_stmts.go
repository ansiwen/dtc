package sqlite3

const CreateTokenTable = `
    CREATE TABLE IF NOT EXISTS token (
        label		PRIMARY KEY,
        pin			TEXT,
        so_pin		TEXT
    )`

const InsertTokenQuery = `	
    INSERT INTO token (label, pin, so_pin) VALUES (?, ?, ?)
`

const GetTokenQuery = `
        SELECT pin, so_pin
        FROM token
        WHERE label = ?
`

const CreateCryptoObjectTable = `
    CREATE TABLE IF NOT EXISTS crypto_object (
        token_label		TEXT,
        handle			INTEGER AUTOINCREMENT, // Auto handle management
        PRIMARY KEY (token_label, handle)
    )`

const InsertCryptoObjectQuery = `
	INSERT OR IGNORE INTO crypto_object (token_label)
	VALUES (?)
`

const CleanCryptoObjectsQuery = `
	DELETE FROM crypto_object WHERE token_label = ?
`

const GetCryptoObjectAttrsQuery = `
        SELECT co.handle, att.type, att.value
		FROM crypto_object as co
        LEFT JOIN attribute as att
		ON att.token_label = co.token_label
		AND att.crypto_object_handle = co.handle
        WHERE co.token_label = ?
`

const CreateAttributeTable = `
    CREATE TABLE IF NOT EXISTS attribute (
        token_label				TEXT,
        crypto_object_handle	INTEGER,
        type					INTEGER,
        value					BLOB,
        PRIMARY KEY (token_label, crypto_object_handle, type)
    )`

const InsertAttributeQuery = `
	INSERT OR REPLACE INTO attribute (token_label, crypto_object_handle, type, value)
	VALUES (?, ?, ?, ?)
`

const CleanAttributesQuery = `
	DELETE FROM attribute WHERE token_label = ?
`

const GetMaxHandleQuery = `
	SELECT COALESCE(MAX(handle), 0) FROM crypto_object
`

var CreateStmts = []string{CreateTokenTable, CreateCryptoObjectTable, CreateAttributeTable}
