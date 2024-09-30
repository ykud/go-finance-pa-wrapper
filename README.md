# go-finance-pa-wrapper
Calculation functions from go-finance for Planning Analytics

# Overview

This is a simple wrapper for [go-finance](https://pkg.go.dev/github.com/alpeb/go-finance) calculation functions to make them available for Planning Analytics `ExecuteCommand`.
[irr](irr.go) allows executing IRR or XIRR calculations on the file in `key,date,casfhlow` format, like so:
```
Project A, 2024-01-01,-1000
Project A, 2024-02-07,400
Project A, 2024-03-01,650
```
and generates a `key,result` csv file.
```
Project A, 0.416184717
```

Execution parameters
```
irr.exe -in file_with_cashflows -out file_with_results -function irr_or_xirr -guess guess_value -log log_file
```
Execution example:
```
irr.exe -in examples\irr_test_cashflows.csv -out examples\irr_test_result.csv -function xirr -guess 0.1 -log calc.log
```
