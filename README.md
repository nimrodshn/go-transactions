# go-transactions
A set of API's for creating and managing SQL transactions over [Contexts](https://golang.org/pkg/context/).

## Creating a new transaction.
Creating a new transaction over a context is a simple as:
```
ctx, err := NewContext(context.Background())
if err != nil {
    // Handle error here.
}
```

## Fetching a transaction from a context:
To fetch a transaction from a context:
```
tx, err := FromContext(ctx)
if err != nil {
    // Handle error here.
}
// resolve the transaction.
defer transaction.Resolve(ctx) 
```

## Marking a transaction for rollback
To mark a transaction for rollback: `MarkForRollback(ctx)`. 

## License
MIT.