
-- name: CreateImageBlock :exec
INSERT INTO image_block (block_id, media_id, alt_text )
VALUES (@block_id, @media_id, @alt_text);

-- name: UpdateImageBlock :exec
UPDATE image_block
SET media_id = @media_id,
    alt_text = @alt_text,
    align = @align,
    object_fit = @object_fit
WHERE block_id = @block_id;

-- name: DeleteImageBlock :exec
DELETE FROM image_block
WHERE block_id = @block_id;

-- name: GetImageBlock :one
SELECT * FROM image_block
WHERE block_id = @block_id;
