-- name: CreateHeadingBlock :exec
INSERT INTO heading_block (block_id, title, description)
VALUES (@block_id, @title, @description);

-- name: UpdateHeadingBlock :exec
UPDATE heading_block
SET title = @title,
    description = @description
WHERE block_id = @block_id;

-- name: GetHeadingBlock :one
SELECT * FROM heading_block WHERE block_id = @block_id;

-- name: DeleteHeadingBlock :exec
DELETE FROM heading_block WHERE block_id = @block_id;
