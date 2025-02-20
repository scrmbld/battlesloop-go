module github.com/scrmbld/battlesloop-go/server

go 1.23.3

replace github.com/scrmbld/battlesloop-go/sloopGame => ../sloopGame

replace github.com/scrmbld/battlesloop-go/sloopNet => ../sloopNet

require (
	github.com/scrmbld/battlesloop-go/sloopGame v0.0.0-00010101000000-000000000000
	github.com/scrmbld/battlesloop-go/sloopNet v0.0.0-00010101000000-000000000000
)
