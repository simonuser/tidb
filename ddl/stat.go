// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package ddl

import (
	"github.com/juju/errors"
	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tidb/meta"
	"github.com/pingcap/tidb/model"
)

func (d *ddl) Stat() (map[string]interface{}, error) {
	var (
		owner     *model.Owner
		job       *model.Job
		schemaVer int64
	)
	err := kv.RunInNewTxn(d.store, false, func(txn kv.Transaction) error {
		t := meta.NewMeta(txn)
		var err error
		owner, err = t.GetDDLOwner()
		if err != nil {
			return errors.Trace(err)
		}

		job, err = t.GetDDLJob(0)
		if err != nil {
			return errors.Trace(err)
		}

		schemaVer, err = t.GetSchemaVersion()
		if err != nil {
			return errors.Trace(err)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Trace(err)
	}

	m := make(map[string]interface{})
	m["DDL_server_id"] = d.uuid

	m["DDL_schema_version"] = schemaVer

	if owner != nil {
		m["DDL_owner_id"] = owner.OwnerID
		// LastUpdateTs uses nanosecond.
		m["DDL_owner_last_update_ts"] = owner.LastUpdateTS / 1e9
	}

	if job != nil {
		m["DDL_job_id"] = job.ID
		m["DDL_job_action"] = job.Type.String()
		m["DDL_job_last_update_ts"] = job.LastUpdateTS
		m["DDL_job_state"] = job.State.String()
		m["DDL_job_error"] = job.Error
		m["DDL_job_schema_state"] = job.SchemaState.String()
		m["DDL_job_schema_id"] = job.SchemaID
		m["DDL_job_table_id"] = job.TableID
		m["DDL_job_snapshot_ver"] = job.SnapshotVer
		m["DDL_job_reorg_handle"] = job.ReorgHandle
		m["DDL_job_args"] = job.Args
	}

	return m, nil
}
