package repo

// selectAllContactsWithPhonesSql contains SQL request that returns all rows where single row is result of merging
// contact row with phone row
const selectAllContactsWithPhonesSql =
/*language=sql*/ `
SELECT
    c.id AS id, c.first_name AS first_name, c.last_name as last_name,
    p.type AS phone_type, p.phone_number AS phone_number
FROM contacts c 
LEFT JOIN phones p on c.id = p.contact_id
ORDER BY c.last_name, c.first_name
`

const insertContactSql =
/*language=sql*/ `
INSERT INTO contacts(first_name, last_name)
VALUES (:firstName, :lastName)
`

const insertPhoneSql =
/*language=sql*/ `
INSERT INTO phones(type, phone_number, contact_id)
VALUES (:type, :phoneNumber, :contactId)
`

const selectContactsWithPhonesByIdSql =
/*language=sql*/ `
SELECT
    c.id AS id, c.first_name AS first_name, c.last_name as last_name,
    p.type AS phone_type, p.phone_number AS phone_number
FROM contacts c 
LEFT JOIN phones p on c.id = p.contact_id
WHERE c.id = :id
`

const deleteContactByIdSql =
/*language=sql*/ `
DELETE FROM contacts WHERE id = :id
`

const deletePhonesByContactIdSql =
/*language=sql*/ `
DELETE FROM phones WHERE contact_id = :contactId
`

const updateContactByIdSql =
/*language=sql*/ `
UPDATE contacts
SET first_name = :firstName,
    last_name = :lastName
WHERE id = :contactId
`
