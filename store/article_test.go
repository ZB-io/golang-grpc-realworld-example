// Remove duplicate declarations of these types since they are already defined in the package
type ArticleStore struct {
    db *gorm.DB
}

type T struct {
    common 
    isEnvSet bool
    context *testContext
}

type ExpectedBegin struct {
    commonExpectation
    delay time.Duration
}

type ExpectedCommit struct {
    commonExpectation
}

type ExpectedExec struct {
    queryBasedExpectation
    result driver.Result 
    delay time.Duration
}

type ExpectedRollback struct {
    commonExpectation
}

type DB struct {
    // ...
}

type Rows struct {
    // ...
}

type Article struct {
    // ...
}

type User struct {
    // ...
}
