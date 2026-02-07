package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		_, err := app.DB().NewQuery(`
drop table if exists message_box;
create table message_box (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	biz_id TEXT NOT NULL,
    status INTEGER NOT NULL,
    message TEXT NOT NULL,
    source_request TEXT NOT NULL,
    source_type INT NOT NULL,
    destination_type INT NOT NULL,
    created_at INT NOT NULL,
    lasted_sent_at INT,
    lasted_error TEXT
);
`).Execute()

		return err
	}, func(app core.App) error {
		// add down queries...

		return nil
	})
}
