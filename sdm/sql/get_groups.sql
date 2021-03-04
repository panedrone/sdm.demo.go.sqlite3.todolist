select g.*, 
(select count(*) from tasks where g_id=g.g_id) as tasks_count
from groups g