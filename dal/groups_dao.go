package dal

// This code was generated by a tool. Don't modify it manually.
// http://sqldalmaker.sourceforge.net

type GroupsDao struct {
	Ds *DataStore
}

// (C)RUD: groups
// Generated values are passed to DTO.

func (dao *GroupsDao) CreateGroup(p *Group) (err error) {
	sql := `insert into groups (g_name) values (?)`
	res, err := dao.Ds.Insert(sql, "g_id", p.GName)
	if err != nil {
		return
	}
	dao.Ds.Assign(&p.GId, res)
	return
}

// C(R)UD: groups

func (dao *GroupsDao) ReadGroup(gId int64) (res Group, err error) {
	sql := `select * from groups where g_id=?`
	rd, err := dao.Ds.QueryRow(sql, gId)
	if err != nil {
		return
	}
	res = Group{}
	dao.Ds.Assign(&res.GId, rd["g_id"] /* q(g_id) <- t(g_id) */)
	dao.Ds.Assign(&res.GName, rd["g_name"] /* q(g_name) <- t(g_name) */)
	return
}

// CR(U)D: groups
// Returns the number of affected rows or -1 on error.

func (dao *GroupsDao) UpdateGroup(p *Group) (res int64, err error) {
	sql := `update groups set g_name=? where g_id=?`
	res, err = dao.Ds.Exec(sql, p.GName, p.GId)
	return
}

// CRU(D): groups
// Returns the number of affected rows or -1 on error.

func (dao *GroupsDao) DeleteGroup(gId int64) (res int64, err error) {
	sql := `delete from groups where g_id=?`
	res, err = dao.Ds.Exec(sql, gId)
	return
}

func (dao *GroupsDao) GetGroups() (res []*Group, err error) {
	sql := `select g.*,  
		(select count(*) from tasks where g_id=g.g_id) as tasks_count 
		from groups g`
	onDto := func(rd map[string]interface{}) {
		obj := Group{}
		dao.Ds.Assign(&obj.GId, rd["g_id"] /* q(g_id) <- q(g_id) */)
		dao.Ds.Assign(&obj.GName, rd["g_name"] /* q(g_name) <- q(g_name) */)
		dao.Ds.Assign(&obj.TasksCount, rd["tasks_count"] /* q(tasks_count) <- q(tasks_count) */)
		res = append(res, &obj)
	}
	err = dao.Ds.QueryAllRows(sql, onDto)
	return
}

func (dao *GroupsDao) GetGroupsIds() (res []int64, err error) {
	sql := `select g.*,  
		(select count(*) from tasks where g_id=g.g_id) as tasks_count 
		from groups g`
	onRow := func(rd interface{}) {
		var data int64
		dao.Ds.Assign(&data, rd)
		res = append(res, data)
	}
	err = dao.Ds.QueryAll(sql, onRow)
	return
}

func (dao *GroupsDao) GetGroupsId() (res int64, err error) {
	sql := `select g.*,  
		(select count(*) from tasks where g_id=g.g_id) as tasks_count 
		from groups g`
	r, err := dao.Ds.Query(sql)
	if err != nil {
		return
	}
	var v int64
	dao.Ds.Assign(&v, r)
	res = v
	return
}
