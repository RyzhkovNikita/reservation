package input

type BarIdInPathInput struct {
	BarId int `input:":bar_id;in_path"`
}

type BarIdInQueryInput struct {
	BarId int `input:":bar_id;in_query"`
}

type TableIdInPathInput struct {
	TableId int `input:":table_id;in_path"`
}
