package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyma-incubator/compass/components/director/pkg/log"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/pkg/persistence"

	"github.com/pkg/errors"
)

// Deleter missing godoc
type Deleter interface {
	DeleteOne(ctx context.Context, resourceType resource.Type, tenant string, conditions Conditions) error
	DeleteMany(ctx context.Context, resourceType resource.Type, tenant string, conditions Conditions) error
}

// DeleterGlobal missing godoc
type DeleterGlobal interface {
	DeleteOneGlobal(ctx context.Context, conditions Conditions) error
	DeleteManyGlobal(ctx context.Context, conditions Conditions) error
}

type universalDeleter struct {
	tableName    string
	resourceType resource.Type
	tenantColumn *string
}

// NewDeleter missing godoc
func NewDeleter(tableName string) Deleter {
	return &universalDeleter{tableName: tableName}
}

// NewDeleterWithEmbeddedTenant missing godoc
func NewDeleterWithEmbeddedTenant(tableName string, tenantColumn string) Deleter {
	return &universalDeleter{tableName: tableName, tenantColumn: &tenantColumn}
}

// NewDeleterGlobal missing godoc
func NewDeleterGlobal(resourceType resource.Type, tableName string) DeleterGlobal {
	return &universalDeleter{tableName: tableName, resourceType: resourceType}
}

// DeleteOne missing godoc
func (g *universalDeleter) DeleteOne(ctx context.Context, resourceType resource.Type, tenant string, conditions Conditions) error {
	if tenant == "" {
		return apperrors.NewTenantRequiredError()
	}

	if g.tenantColumn != nil {
		conditions = append(Conditions{NewEqualCondition(*g.tenantColumn, tenant)}, conditions...)
		return g.unsafeDelete(ctx, resourceType, conditions, true)
	}

	if resourceType.IsTopLevel() {
		return g.unsafeDeleteTenantAccess(ctx, resourceType, tenant, conditions, true)
	}

	return g.unsafeDeleteChildEntity(ctx, resourceType, tenant, conditions, true)
}

// DeleteMany missing godoc
func (g *universalDeleter) DeleteMany(ctx context.Context, resourceType resource.Type, tenant string, conditions Conditions) error {
	if tenant == "" {
		return apperrors.NewTenantRequiredError()
	}

	if g.tenantColumn != nil {
		conditions = append(Conditions{NewEqualCondition(*g.tenantColumn, tenant)}, conditions...)
		return g.unsafeDelete(ctx, resourceType, conditions, false)
	}

	if resourceType.IsTopLevel() {
		return g.unsafeDeleteTenantAccess(ctx, resourceType, tenant, conditions, false)
	}

	return g.unsafeDeleteChildEntity(ctx, resourceType, tenant, conditions, false)
}

// DeleteOneGlobal missing godoc
func (g *universalDeleter) DeleteOneGlobal(ctx context.Context, conditions Conditions) error {
	return g.unsafeDelete(ctx, g.resourceType, conditions, true)
}

// DeleteManyGlobal missing godoc
func (g *universalDeleter) DeleteManyGlobal(ctx context.Context, conditions Conditions) error {
	return g.unsafeDelete(ctx, g.resourceType, conditions, false)
}

func (g *universalDeleter) unsafeDelete(ctx context.Context, resourceType resource.Type, conditions Conditions, requireSingleRemoval bool) error {
	persist, err := persistence.FromCtx(ctx)
	if err != nil {
		return err
	}

	var stmtBuilder strings.Builder

	stmtBuilder.WriteString(fmt.Sprintf("DELETE FROM %s", g.tableName))

	if len(conditions) > 0 {
		stmtBuilder.WriteString(" WHERE")
	}

	err = writeEnumeratedConditions(&stmtBuilder, conditions)
	if err != nil {
		return errors.Wrap(err, "while writing enumerated conditions")
	}
	allArgs := getAllArgs(conditions)

	query := getQueryFromBuilder(stmtBuilder)
	log.C(ctx).Debugf("Executing DB query: %s", query)
	res, err := persist.ExecContext(ctx, query, allArgs...)
	if err = persistence.MapSQLError(ctx, err, resourceType, resource.Delete, "while deleting object from '%s' table", g.tableName); err != nil {
		return err
	}

	if requireSingleRemoval {
		affected, err := res.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "while checking affected rows")
		}
		if affected != 1 {
			return apperrors.NewInternalError("delete should remove single row, but removed %d rows", affected)
		}
	}

	return nil
}

func (g *universalDeleter) unsafeDeleteTenantAccess(ctx context.Context, resourceType resource.Type, tenant string, conditions Conditions, requireSingleRemoval bool) error {
	m2mTable, ok := resourceType.TenantAccessTable()
	if !ok {
		return errors.Errorf("entity %s does not have access table", resourceType)
	}

	persist, err := persistence.FromCtx(ctx)
	if err != nil {
		return err
	}

	var stmtBuilder strings.Builder

	stmtBuilder.WriteString(fmt.Sprintf("SELECT id FROM %s WHERE", g.tableName))

	tenantIsolation, err := NewTenantIsolationCondition(resourceType, tenant, true)
	if err != nil {
		return err
	}

	conditions = append(conditions, tenantIsolation)

	err = writeEnumeratedConditions(&stmtBuilder, conditions)
	if err != nil {
		return errors.Wrap(err, "while writing enumerated conditions")
	}
	allArgs := getAllArgs(conditions)

	query := getQueryFromBuilder(stmtBuilder)
	log.C(ctx).Debugf("Executing DB query: %s", query)

	var ids IDs
	err = persist.SelectContext(ctx, &ids, query, allArgs...)
	if err = persistence.MapSQLError(ctx, err, resourceType, resource.Delete, "while selecting objects from '%s' table by conditions", g.tableName); err != nil {
		return err
	}

	if len(ids) == 0 {
		return apperrors.NewUnauthorizedError(apperrors.ShouldBeOwnerMsg)
	}

	if requireSingleRemoval && len(ids) != 1 {
		return apperrors.NewInternalError("delete should remove single row, but removed %d rows", len(ids))
	}

	stmtBuilder.Reset()

	stmtBuilder.WriteString(fmt.Sprintf("DELETE FROM %s WHERE", m2mTable))

	deleteConditions := Conditions{NewInConditionForStringValues(M2MResourceIDColumn, ids)}
	err = writeEnumeratedConditions(&stmtBuilder, deleteConditions)
	if err != nil {
		return errors.Wrap(err, "while writing enumerated conditions")
	}

	allArgs = getAllArgs(deleteConditions)

	query = getQueryFromBuilder(stmtBuilder)
	log.C(ctx).Debugf("Executing DB query: %s", query)

	_, err = persist.ExecContext(ctx, query, allArgs...)
	return persistence.MapSQLError(ctx, err, resourceType, resource.Delete, "while deleting objects from '%s' table", m2mTable)
}

func (g *universalDeleter) unsafeDeleteChildEntity(ctx context.Context, resourceType resource.Type, tenant string, conditions Conditions, requireSingleRemoval bool) error {
	persist, err := persistence.FromCtx(ctx)
	if err != nil {
		return err
	}

	var stmtBuilder strings.Builder
	stmtBuilder.WriteString(fmt.Sprintf("DELETE FROM %s WHERE", g.tableName))

	tenantIsolation, err := NewTenantIsolationCondition(resourceType, tenant, true)
	if err != nil {
		return err
	}

	conditions = append(conditions, tenantIsolation)

	err = writeEnumeratedConditions(&stmtBuilder, conditions)
	if err != nil {
		return errors.Wrap(err, "while writing enumerated conditions")
	}
	allArgs := getAllArgs(conditions)

	query := getQueryFromBuilder(stmtBuilder)
	log.C(ctx).Debugf("Executing DB query: %s", query)

	res, err := persist.ExecContext(ctx, query, allArgs...)
	if err = persistence.MapSQLError(ctx, err, resourceType, resource.Delete, "while deleting object from '%s' table", g.tableName); err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "while checking affected rows")
	}

	if affected == 0 {
		return apperrors.NewUnauthorizedError(apperrors.ShouldBeOwnerMsg)
	}

	if requireSingleRemoval && affected != 1 {
		return apperrors.NewInternalError("delete should remove single row, but removed %d rows", affected)
	}

	return nil
}

// IDs keeps IDs retrieved from the Compass storage.
type IDs []string

// Len returns the length of the IDs
func (i IDs) Len() int {
	return len(i)
}
