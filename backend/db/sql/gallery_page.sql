




select p.* from gallery_page gp
join pages p on p.id = gp.page_id
where gp.gallery_id = 'a690ba02-aceb-45a3-852b-b22ea6e07d5f'
order by  gp.sort_order ASC;
