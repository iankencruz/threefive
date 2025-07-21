
-- name: GetBlocksForPage :many
SELECT * FROM blocks
WHERE parent_type = 'page' AND parent_id = @page_id
ORDER BY sort_order;

-- name: GetBlocksForBlog :many
SELECT * FROM blocks
WHERE parent_type = 'blog' AND parent_id = @blog_id
ORDER BY sort_order;

-- name: GetBlockByID :one
SELECT * FROM blocks WHERE id = @id;



-- name: CreateBlock :one
INSERT INTO blocks (id, parent_type, parent_id, type, sort_order)
VALUES (@id, @parent_type, @parent_id, @type, @sort_order)
RETURNING *;


-- name: UpdateBlock :exec
UPDATE blocks
SET
  parent_type = @parent_type,
  parent_id   = @parent_id,
  type        = @type,
  sort_order  = @sort_order,
  updated_at  = now()
WHERE id = @id;



-- name: DeleteBlocksByParent :exec
DELETE FROM blocks
WHERE parent_type = @parent_type AND parent_id = @parent_id;


-- -- name: UpdateBlockSortOrder :exec
-- UPDATE blocks
-- SET sort_order = @sort_order,
--     updated_at = NOW()
-- WHERE id = @id;
--




-- name: UpdateBlockSortOrder :exec
UPDATE blocks
SET sort_order = $2,
    updated_at = now()
WHERE id = $1;


-- name: DeleteBlock :exec
DELETE FROM blocks
WHERE id = @id;




-- name: UpdateBlock :exec


