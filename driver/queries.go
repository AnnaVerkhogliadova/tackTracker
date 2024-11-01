package driver

const (
	queryCreateTask = `
	INSERT 
	INTO tasks.tasks
	    (title, description, status)
	VALUES 
	    ($1, $2, $3)
	returning task_id;`

	queryGet = `
	SELECT
    	t.task_id, t.title, t.description, t.status, t.create_date,
    	COALESCE(json_agg(to_jsonb(st) - '{task_id}'::text[]) FILTER (WHERE st.sub_task_id IS NOT NULL), '[]'::json) AS sub_tasks
	FROM tasks.tasks t
		LEFT JOIN tasks.sub_tasks st ON t.task_id = st.task_id
	WHERE t.task_id = $1
	GROUP BY t.task_id;
`
	querySetStatus = `
	UPDATE tasks.tasks
		SET status = $2 WHERE task_id = $1
`

	queryDelete = `
	DELETE 
		FROM tasks.tasks
	WHERE
		task_id = $1
`
	queryGetList = `
	SELECT
    	t.task_id, t.title, t.description, t.status, t.create_date,
    	COALESCE(json_agg(to_jsonb(st) - '{task_id}'::text[]) FILTER (WHERE st.sub_task_id IS NOT NULL), '[]'::json) AS sub_tasks
	FROM tasks.tasks t
		LEFT JOIN tasks.sub_tasks st ON t.task_id = st.task_id
	WHERE t.status = coalesce($1, t.status)
	GROUP BY t.task_id;
`

	queryCreateSubTask = `
	INSERT 
	INTO tasks.sub_tasks
	    (task_id, title, description, status)
	VALUES 
	    ($1, $2, $3, $4)
	returning sub_task_id;`

	queryExistTaskId = `
	SELECT exists(
    	SELECT 1 FROM tasks.tasks WHERE tasks.task_id = $1)
`
)
