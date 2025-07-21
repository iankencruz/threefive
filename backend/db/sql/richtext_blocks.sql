
-- name: CreateRichTextBlock :exec
INSERT INTO richtext_block (block_id, html)
VALUES (@block_id, @html);

-- name: UpdateRichTextBlock :exec
UPDATE richtext_block
SET html = @html
WHERE block_id = @block_id;

-- name: DeleteRichTextBlock :exec
DELETE FROM richtext_block
WHERE block_id = @block_id;

-- name: GetRichTextBlock :one
SELECT * FROM richtext_block
WHERE block_id = @block_id;

