package driver

const (
	queryCreateTask = `
	insert 
	into tasks.tasks
	    (title, description, status)
	values
	    ($1, $2, $3)
	returning task_id;`

	queryGet = `
select 
    tasks.task_id, tasks.title, tasks.description, tasks.status, tasks.create_date from tasks.tasks
	where task_id = $1;
`
	queryDelete = `
	delete 
		from tasks.tasks
	where
		task_id = $1
`
)
