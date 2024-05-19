package tarantoolrepo

const (
	migrateMessagesQuery = `box.schema.space.create('messages', {if_not_exists=true})`

	migrateMessagesPrimaryIndexQuery = `box.space.messages:create_index('primary', {type="HASH", unique=true, parts={1, 'string'}, if_not_exists=true})`

	migrateMessagesIndexQuery = `box.space.messages:create_index('dialog_id', {type="TREE", unique=false, parts={2, 'string'}, if_not_exists=true})`

	migrateGetFuncQuery = `box.schema.func.create('get_dialog', {body=[[function(dialog_id) return box.space.messages.index.dialog_id:select({dialog_id}) end]], if_not_exists=true})`

	migrateInsertFuncQuery = `box.schema.func.create('insert_dialog', {body=[[function(id, dialog_id, from, to, text) return box.space.messages:insert({id, dialog_id, from, to, text}) end]], if_not_exists=true})`
)

const (
	getDialogFuncName = `get_dialog`

	insertDialogFuncName = `insert_dialog`
)
