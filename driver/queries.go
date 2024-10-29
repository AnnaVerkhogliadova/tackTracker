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
    tasks.task_id, tasks.title, tasks.description, tasks.status, tasks.create_date FROM tasks.tasks
	WHERE task_id = $1;
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
		 tasks.task_id, tasks.title, tasks.description, tasks.status, tasks.create_date FROM tasks.tasks
	WHERE
		status = coalesce($1, status) 
`

	queryCreateSubTask = `
	INSERT 
	INTO tasks.sub_tasks
	    (task_id, title, description, status)
	VALUES 
	    ($1, $2, $3, $4)
	returning sub_task_id;`

	queryExistTaskId = `
	select exists(
    select 1 from tasks.tasks where tasks.task_id = $1
           )
`
)
