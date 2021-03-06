package postgres_test

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/services/sqlstore/postgres"

	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/web"
)

var trx *dbx.Trx

var demoTenant *models.Tenant
var avengersTenant *models.Tenant
var gotTenant *models.Tenant
var jonSnow *models.User
var aryaStark *models.User
var sansaStark *models.User
var tonyStark *models.User

var demoTenantCtx context.Context
var avengersTenantCtx context.Context
var gotTenantCtx context.Context

var jonSnowCtx context.Context
var aryaStarkCtx context.Context
var sansaStarkCtx context.Context
var tonyStarkCtx context.Context

func SetupDatabaseTest(t *testing.T) context.Context {
	RegisterT(t)
	bus.Init(postgres.Service{})

	u, _ := url.Parse("http://cdn.test.fider.io")
	req := web.Request{URL: u}
	ctx := context.WithValue(context.Background(), app.RequestCtxKey, req)

	trx, _ = dbx.BeginTx(ctx)
	trxCtx := context.WithValue(ctx, app.TransactionCtxKey, trx)

	getDemo := &query.GetTenantByDomain{Domain: "demo"}
	getAvengers := &query.GetTenantByDomain{Domain: "avengers"}
	getGameOfThrones := &query.GetTenantByDomain{Domain: "got"}
	bus.Dispatch(trxCtx, getDemo, getAvengers, getGameOfThrones)
	demoTenant = getDemo.Result
	avengersTenant = getAvengers.Result
	gotTenant = getGameOfThrones.Result

	demoTenantCtx = withTenant(trxCtx, demoTenant)
	avengersTenantCtx = withTenant(trxCtx, avengersTenant)
	gotTenantCtx = withTenant(trxCtx, gotTenant)

	getJonSnow := &query.GetUserByEmail{Email: "jon.snow@got.com"}
	getAryaStark := &query.GetUserByEmail{Email: "arya.stark@got.com"}
	getSansaStark := &query.GetUserByEmail{Email: "sansa.stark@got.com"}
	bus.Dispatch(demoTenantCtx, getJonSnow, getSansaStark, getAryaStark)
	jonSnow = getJonSnow.Result
	aryaStark = getAryaStark.Result
	sansaStark = getSansaStark.Result

	getTonyStark := &query.GetUserByEmail{Email: "tony.stark@avengers.com"}
	bus.Dispatch(avengersTenantCtx, getTonyStark)
	tonyStark = getTonyStark.Result

	jonSnowCtx = withUser(trxCtx, jonSnow)
	aryaStarkCtx = withUser(trxCtx, aryaStark)
	sansaStarkCtx = withUser(trxCtx, sansaStark)
	tonyStarkCtx = withUser(trxCtx, tonyStark)
	return trxCtx
}

func TeardownDatabaseTest() {
	trx.Rollback()
}

func TestMain(m *testing.M) {
	dbx.Seed()

	code := m.Run()
	os.Exit(code)
}
