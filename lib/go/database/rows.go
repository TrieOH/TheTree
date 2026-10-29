package database

import "lib/xslices"

// One finishes a single-row query: the row mapped to its domain type, or
// the query's error turned into the Backend's named error through dbe.
//
//	row, err := repo.db(ctx).GetEditionByID(ctx, id)
//	return database.One(row, err, repo.dbe, mapEdition)
func One[R, M any](row R, err error, dbe ErrorHandler, m func(R) M) (*M, error) {
	if err != nil {
		return nil, dbe(err)
	}
	return new(m(row)), nil
}

// Many is One for a multi-row query.
func Many[R, M any](rows []R, err error, dbe ErrorHandler, m func(R) M) ([]M, error) {
	if err != nil {
		return nil, dbe(err)
	}
	return xslices.MapSlice(rows, m), nil
}
