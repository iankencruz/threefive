

-- name: GetPageGalleries :many
SELECT g.*
FROM galleries g
INNER JOIN gallery_page gp ON gp.gallery_id = g.id
WHERE gp.page_id = @page_id
ORDER BY gp.sort_order ASC;

-- name: LinkGalleryToPage :exec
INSERT INTO gallery_page (gallery_id, page_id, sort_order)
VALUES (@gallery_id, @page_id, COALESCE(
  (SELECT MAX(sort_order) + 1 FROM gallery_page WHERE page_id = @page_id), 0
));

-- name: UnlinkGalleryFromPage :exec
DELETE FROM gallery_page
WHERE gallery_id = @gallery_id AND page_id = @page_id;

